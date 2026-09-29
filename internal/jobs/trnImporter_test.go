package jobs

import (
	"converterapi/internal/config"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"testing"
)

// Соединение, к которому не удалось подключиться, не должно оставлять
// половинчатое состояние: следующий заход джобы обязан попробовать заново.
func TestSftpConnFailedConnectLeavesNothing(t *testing.T) {
	// Слушатель поднимаем и сразу закрываем - порт гарантированно свободен
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("занять порт: %v", err)
	}
	addr := l.Addr().(*net.TCPAddr)
	l.Close()

	old := config.Config.Processing.SFTP
	defer func() { config.Config.Processing.SFTP = old }()
	config.Config.Processing.SFTP.Host = "127.0.0.1"
	config.Config.Processing.SFTP.Port = strconv.Itoa(addr.Port)

	var c sftpConn
	if _, err := c.client(); err == nil {
		t.Fatal("подключение к закрытому порту должно возвращать ошибку")
	}
	if c.sftp != nil || c.ssh != nil {
		t.Errorf("после неудачи соединение должно остаться пустым: sftp=%v ssh=%v", c.sftp, c.ssh)
	}
}

// drop на незанятом соединении не должен паниковать и обязан быть идемпотентным:
// его зовут и из обработки ошибок, и при остановке сервиса.
func TestSftpConnDropIdempotent(t *testing.T) {
	var c sftpConn
	c.drop()
	c.drop()
	if c.sftp != nil || c.ssh != nil {
		t.Error("drop обязан обнулять обе сущности")
	}
	CloseSFTP()
}

// Файл, уже лежащий в каталоге загрузки, повторно скачивать не нужно - раньше
// учитывался только каталог success, куда его перекладывает внешний обработчик.
func TestTrackerCountsDownloadDir(t *testing.T) {
	downloads := t.TempDir()
	success := t.TempDir()

	if err := os.WriteFile(filepath.Join(downloads, "trn_1.json"), []byte("{}"), 0644); err != nil {
		t.Fatalf("подготовка файла: %v", err)
	}
	if err := os.WriteFile(filepath.Join(success, "trn_2.json"), []byte("{}"), 0644); err != nil {
		t.Fatalf("подготовка файла: %v", err)
	}

	tracker, err := NewImportedFilesTracker(downloads, success)
	if err != nil {
		t.Fatalf("трекер: %v", err)
	}

	for _, name := range []string{"trn_1.json", "trn_2.json"} {
		if !tracker.isImported(name) {
			t.Errorf("%s уже скачан, повторная загрузка не нужна", name)
		}
	}
	if tracker.isImported("trn_3.json") {
		t.Error("незнакомый файл должен скачиваться")
	}
}

// Отсутствующий каталог трекер создаёт, а не падает
func TestTrackerCreatesMissingDir(t *testing.T) {
	missing := filepath.Join(t.TempDir(), "success")
	if _, err := NewImportedFilesTracker(missing); err != nil {
		t.Fatalf("трекер: %v", err)
	}
	if _, err := os.Stat(missing); err != nil {
		t.Errorf("каталог должен быть создан: %v", err)
	}
}

// Партнёр называет пакеты по-своему - имя приводим к виду, который понимает
// сканер, иначе файл теряется молча
func TestOfflineLocalName(t *testing.T) {
	tests := map[string]string{
		"CreateCards_Out_2026.09.29_16.36.44":     "CreateCardsOut_2026.09.29_16.36.44.xml",
		"CreateCardsOut_2026.09.29_16.36.44.xml":  "CreateCardsOut_2026.09.29_16.36.44.xml",
		"createcardsout20260929":                  "CreateCardsOut_20260929.xml",
		"Reissue_Cards_Out_001":                   "ReissueCardsOut_001.xml",
		"RelinkPreIssuedCards_Out_5":              "RelinkPreIssuedCardsOut_5.xml",
		"CreatePreIssuedCards_2026.09.29":         "CreatePreIssuedCards_2026.09.29.xml",
		"CreateCustomerAndAccount_2026.09.29.xml": "CreateCustomerAndAccount_2026.09.29.xml",

		// не наши файлы
		"trn_20260929.json": "",
		"readme.txt":        "",
		"":                  "",
	}

	for remote, want := range tests {
		if got := offlineLocalName(remote); got != want {
			t.Errorf("offlineLocalName(%q) = %q, want %q", remote, got, want)
		}
	}
}

// Транзакционные файлы офлайн-импортом не забираются: у них свой маршрут
func TestOfflineLocalNameSkipsTransactions(t *testing.T) {
	for _, name := range []string{"G2BTRN-20260929.json", "out.json", "20260929.json"} {
		if got := offlineLocalName(name); got != "" {
			t.Errorf("%q не пакет выпуска, получено %q", name, got)
		}
	}
}
