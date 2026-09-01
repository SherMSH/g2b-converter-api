package jobs

import (
	"converterapi/internal/config"
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
)

// ImportedFilesTracker отслеживает импортированные файлы
type ImportedFilesTracker struct {
	importedDir string
	importedMap map[string]bool
}

// NewImportedFilesTracker создает новый трекер импортированных файлов
func NewImportedFilesTracker(importedDir string) (*ImportedFilesTracker, error) {
	tracker := &ImportedFilesTracker{
		importedDir: importedDir,
		importedMap: make(map[string]bool),
	}

	// Сканируем уже импортированные файлы
	if err := tracker.scanImportedFiles(); err != nil {
		return nil, fmt.Errorf("ошибка сканирования импортированных файлов: %w", err)
	}

	return tracker, nil
}

// scanImportedFiles сканирует локальную директорию импортированных файлов
func (t *ImportedFilesTracker) scanImportedFiles() error {
	files, err := os.ReadDir(t.importedDir)
	if err != nil {
		if os.IsNotExist(err) {
			// Создаем директорию, если ее нет
			return os.MkdirAll(t.importedDir, 0755)
		}
		return err
	}

	for _, file := range files {
		if !file.IsDir() {
			t.importedMap[file.Name()] = true
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

// connectSFTP устанавливает SFTP соединение
func connectSFTP() (*sftp.Client, error) {
	// Настраиваем SSH клиент
	config := config.Config.Processing.SFTP
	sshConfig := &ssh.ClientConfig{
		User: config.Name,
		Auth: []ssh.AuthMethod{
			ssh.Password(config.Token),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), // В продакшене используйте нормальную проверку ключей
		Timeout:         30 * time.Second,
	}

	// Подключаемся
	addr := fmt.Sprintf("%s:%s", config.Host, config.Port)
	conn, err := ssh.Dial("tcp", addr, sshConfig)
	if err != nil {
		return nil, fmt.Errorf("не удалось подключиться к SSH: %w", err)
	}

	// Создаем SFTP клиент
	client, err := sftp.NewClient(conn)
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("не удалось создать SFTP клиент: %w", err)
	}

	return client, nil
}

// importSingleFile импортирует один файл через SFTP
func importSingleFile(client *sftp.Client, remotePath, localDir string) error {
	// Получаем имя файла
	fileName := filepath.Base(remotePath)
	localPath := filepath.Join(localDir, fileName)

	// Проверяем, не существует ли уже файл локально
	if _, err := os.Stat(localPath); err == nil {
		return fmt.Errorf("локальный файл уже существует: %s", localPath)
	}

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

	// Получаем информацию о файле для проверки размера
	remoteInfo, err := client.Stat(remotePath)
	if err != nil {
		os.Remove(localPath)
		return fmt.Errorf("ошибка получения информации о файле: %w", err)
	}

	if bytesCopied != remoteInfo.Size() {
		os.Remove(localPath)
		return fmt.Errorf("скопировано %d байт, ожидалось %d", bytesCopied, remoteInfo.Size())
	}

	return nil
}

// ImportAllFilesWithFilter импортирует только файлы, соответствующие фильтру
func ImportAllFilesWithFilter(remoteDir, localDir string, filter func(string) bool) error {
	client, err := connectSFTP()
	if err != nil {
		return err
	}
	defer client.Close()

	tracker, err := NewImportedFilesTracker(localDir)
	if err != nil {
		return err
	}

	files, err := client.ReadDir(remoteDir)
	if err != nil {
		return err
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		fileName := file.Name()

		// Применяем фильтр
		if filter != nil && !filter(fileName) {
			fmt.Printf("Файл %s не соответствует фильтру, пропускаем\n", fileName)
			continue
		}

		if tracker.isImported(fileName) {
			fmt.Printf("Файл %s уже импортирован\n", fileName)
			continue
		}

		remotePath := filepath.Join(remoteDir, fileName)
		if err := importSingleFile(client, remotePath, localDir); err != nil {
			fmt.Printf("Ошибка импорта %s: %v\n", fileName, err)
			continue
		}

		tracker.markAsImported(fileName)
		fmt.Printf("Файл %s импортирован\n", fileName)
	}

	return nil
}

var mutx sync.Mutex

func TrnImporter() {
	mutx.Lock()
	defer mutx.Unlock()
	logger.Infof("[JOBS] TRN files importer")

	remoteDir := "out"
	localDir := "/home/sherzodm/trn/d8" //"/srv/g2b/files/trn/d8"

	filter := func(filename string) bool {
		return strings.HasSuffix(filename, ".json")
	}

	if err := ImportAllFilesWithFilter(remoteDir, localDir, filter); err != nil {
		fmt.Printf("Ошибка импорта с фильтром: %v\n", err)
	}
}
