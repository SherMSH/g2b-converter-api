package jobs

import (
	"converterapi/internal/config"
	"net"
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
