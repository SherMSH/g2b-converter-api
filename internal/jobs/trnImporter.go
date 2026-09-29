package jobs

import (
	configpkg "converterapi/internal/config"
	"converterapi/pkg/logger"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/knownhosts"
)

// ImportedFilesTracker отслеживает импортированные файлы
type ImportedFilesTracker struct {
	dirs        []string
	importedMap map[string]bool
}

// NewImportedFilesTracker создает новый трекер импортированных файлов.
//
// Смотрим сразу в два каталога: куда файлы скачиваются и куда их переносит
// обработчик после разбора. Раньше учитывался только второй, а кладёт туда
// файлы внешний процесс - пока он этого не сделал, файл считался неимпортированным
// и скачивался заново на каждом заходе джобы.
func NewImportedFilesTracker(dirs ...string) (*ImportedFilesTracker, error) {
	tracker := &ImportedFilesTracker{
		dirs:        dirs,
		importedMap: make(map[string]bool),
	}

	// Сканируем уже импортированные файлы
	if err := tracker.scanImportedFiles(); err != nil {
		return nil, fmt.Errorf("ошибка сканирования импортированных файлов: %w", err)
	}

	return tracker, nil
}

// scanImportedFiles сканирует локальные директории импортированных файлов
func (t *ImportedFilesTracker) scanImportedFiles() error {
	for _, dir := range t.dirs {
		files, err := os.ReadDir(dir)
		if err != nil {
			if os.IsNotExist(err) {
				// Создаем директорию, если ее нет
				if err := os.MkdirAll(dir, 0755); err != nil {
					return err
				}
				continue
			}
			return err
		}

		for _, file := range files {
			if !file.IsDir() {
				t.importedMap[file.Name()] = true
			}
		}
	}

	return nil
}

// isImported проверяет, был ли файл уже импортирован
func (t *ImportedFilesTracker) isImported(fileName string) bool {
	return t.importedMap[fileName]
}

// markAsImported отмечает файл как импортированный
func (t *ImportedFilesTracker) markAsImported(fileName string) {
	t.importedMap[fileName] = true
}

// connectSFTP устанавливает SFTP соединение.
//
// Возвращает обе сущности: sftp.Client.Close() закрывает только подсистему
// sftp внутри SSH-сессии, а само TCP-соединение остаётся висеть на сервере -
// закрывать его нужно отдельно.
func connectSFTP() (*sftp.Client, *ssh.Client, error) {
	// Настраиваем SSH клиент
	config := configpkg.Config.Processing.SFTP
	sshConfig := &ssh.ClientConfig{
		User: config.Name,
		Auth: []ssh.AuthMethod{
			ssh.Password(config.Token),
		},
		HostKeyCallback: hostKeyCallback(config),
		Timeout:         30 * time.Second,
	}

	// Подключаемся
	addr := fmt.Sprintf("%s:%s", config.Host, config.Port)
	conn, err := ssh.Dial("tcp", addr, sshConfig)
	if err != nil {
		return nil, nil, fmt.Errorf("не удалось подключиться к SSH: %w", err)
	}

	// Создаем SFTP клиент
	client, err := sftp.NewClient(conn)
	if err != nil {
		conn.Close()
		return nil, nil, fmt.Errorf("не удалось создать SFTP клиент: %w", err)
	}

	return client, conn, nil
}

// hostKeyCallback задаёт проверку ключа сервера.
//
// Путь к known_hosts берётся из processing.sftp.extra.known_hosts. Если он не
// задан, ключ не проверяется - так работало всегда, и менять это молча нельзя:
// сервис просто перестал бы подключаться. Поэтому о небезопасном режиме
// предупреждаем в логе.
func hostKeyCallback(cfg configpkg.Server) ssh.HostKeyCallback {
	path := cfg.Extra["known_hosts"]
	if path == "" {
		logger.Warnf("[JOBS] SFTP: known_hosts не задан, ключ сервера не проверяется")
		return ssh.InsecureIgnoreHostKey()
	}

	callback, err := knownhosts.New(path)
	if err != nil {
		logger.Errorf("[JOBS] SFTP: не удалось прочитать known_hosts %s: %v", path, err)
		return ssh.InsecureIgnoreHostKey()
	}
	return callback
}

// sftpConn - переиспользуемое соединение с файловым сервером процессинга.
//
// Джоба ходит за файлами каждые несколько секунд, и на каждый заход поднимать
// SSH заново дорого: обмен ключами - самая тяжёлая часть рукопожатия. Поэтому
// соединение живёт между запусками, а переподключаемся только когда оно
// отвалилось.
type sftpConn struct {
	mu   sync.Mutex
	sftp *sftp.Client
	ssh  *ssh.Client
}

// trnSFTP - соединение джобы импорта транзакционных файлов
var trnSFTP sftpConn

// client отдаёт живое соединение, при необходимости переподключаясь.
//
// Живость проверяем дешёвым запросом рабочего каталога: обрыв на той стороне
// молча не обнаруживается, а писать в мёртвое соединение - значит потерять
// заход джобы.
func (c *sftpConn) client() (*sftp.Client, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.sftp != nil {
		if _, err := c.sftp.Getwd(); err == nil {
			return c.sftp, nil
		}
		logger.Warnf("[JOBS] SFTP соединение потеряно, переподключаемся")
		c.closeLocked()
	}

	client, conn, err := connectSFTP()
	if err != nil {
		return nil, err
	}
	c.sftp, c.ssh = client, conn
	logger.Infof("[JOBS] SFTP соединение установлено")
	return c.sftp, nil
}

// drop закрывает соединение, чтобы следующий заход поднял новое.
// Вызывается после ошибок работы с файлами: причина могла быть в соединении.
func (c *sftpConn) drop() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.closeLocked()
}

// closeLocked закрывает обе сущности. Вызывать под mu.
func (c *sftpConn) closeLocked() {
	if c.sftp != nil {
		if err := c.sftp.Close(); err != nil {
			logger.Warnf("[JOBS] закрытие SFTP клиента: %v", err)
		}
		c.sftp = nil
	}
	// Отдельно от sftp: без этого TCP-соединение остаётся открытым на сервере
	if c.ssh != nil {
		if err := c.ssh.Close(); err != nil {
			logger.Warnf("[JOBS] закрытие SSH соединения: %v", err)
		}
		c.ssh = nil
	}
}

// CloseSFTP закрывает соединение импортёра. Вызывается при остановке сервиса.
func CloseSFTP() {
	trnSFTP.drop()
}

// importSingleFile импортирует один файл через SFTP
func importSingleFile(client *sftp.Client, remotePath, localDir string, remoteSize int64, localName ...string) error {
	// Имя на диске может отличаться от удалённого: пакеты выпуска карт
	// приводим к виду, который понимает сканер
	fileName := filepath.Base(remotePath)
	if len(localName) != 0 && localName[0] != "" {
		fileName = localName[0]
	}
	localPath := filepath.Join(localDir, fileName)

	// // Проверяем, не существует ли уже файл локально
	// if _, err := os.Stat(localPath); err == nil {
	// 	return fmt.Errorf("локальный файл уже существует: %s", localPath)
	// }

	// Открываем удаленный файл
	remoteFile, err := client.Open(remotePath)
	if err != nil {
		return fmt.Errorf("не удалось открыть удаленный файл: %w", err)
	}
	defer remoteFile.Close()

	// Создаем локальный файл
	localFile, err := os.Create(localPath)
	if err != nil {
		return fmt.Errorf("не удалось создать локальный файл: %w", err)
	}
	defer localFile.Close()

	// Копируем содержимое
	bytesCopied, err := io.Copy(localFile, remoteFile)
	if err != nil {
		os.Remove(localPath)
		return fmt.Errorf("ошибка копирования: %w", err)
	}

	// Размер берём из листинга каталога: отдельный Stat - это лишний
	// round-trip к серверу на каждый файл
	if bytesCopied != remoteSize {
		os.Remove(localPath)
		return fmt.Errorf("скопировано %d байт, ожидалось %d", bytesCopied, remoteSize)
	}

	return nil
}

// ImportAllFilesWithFilter импортирует только файлы, соответствующие фильтру
func ImportAllFilesWithFilter(remoteDir, localDir string, filter func(string) bool) error {
	client, err := trnSFTP.client()
	if err != nil {
		return err
	}

	tracker, err := NewImportedFilesTracker(localDir, localDir+"/../success")
	if err != nil {
		return err
	}

	files, err := client.ReadDir(remoteDir)
	if err != nil {
		// Каталог мог стать недоступен из-за обрыва - соединение переподнимем
		trnSFTP.drop()
		return err
	}

	var imported int
	for _, file := range files {
		if file.IsDir() {
			continue
		}

		fileName := file.Name()

		// Применяем фильтр
		if filter != nil && !filter(fileName) {
			continue
		}

		if tracker.isImported(fileName) {
			continue
		}

		remotePath := remoteDir + "/" + fileName
		if err := importSingleFile(client, remotePath, localDir, file.Size()); err != nil {
			logger.Errorf("[JOBS] ошибка импорта %s: %v", fileName, err)
			continue
		}

		tracker.markAsImported(fileName)
		imported++
		logger.Infof("[JOBS] файл %s импортирован (%d байт)", fileName, file.Size())
	}

	if imported == 0 {
		logger.Debugf("[JOBS] новых файлов нет")
	}

	return nil
}

var mutx sync.Mutex

func TrnImporter() {
	mutx.Lock()
	defer mutx.Unlock()
	logger.Infof("[JOBS] TRN files importer")

	remoteDir := configpkg.Config.Jobs.TrnImporter.Extra["remote"]
	localDir := configpkg.Config.Jobs.TrnImporter.Extra["local"]

	filter := func(filename string) bool {
		return strings.HasSuffix(filename, ".json")
	}

	if err := ImportAllFilesWithFilter(remoteDir, localDir, filter); err != nil {
		logger.Errorf("[JOBS] TRN импорт: %v", err)
	}

	// В том же каталоге партнёр выкладывает пакеты на выпуск карт - забираем и
	// их, соединение переиспользуется
	if err := ImportOfflinePackets(remoteDir); err != nil {
		logger.Errorf("[JOBS] импорт офлайн-пакетов: %v", err)
	}
}
