package storage

import (
	"fmt"
	"io"
	"os"
)

func DownloadFile(path string) ([]byte, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, err
	}
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("не удалось открыть файл: %w", err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("ошибка чтения файла: %w", err)
	}
	return data, nil
}

// MoveFile использует os.Rename для перемещения файла
func MoveFile(sourcePath, destPath string) error {
	// Проверяем, существует ли исходный файл
	if _, err := os.Stat(sourcePath); os.IsNotExist(err) {
		return fmt.Errorf("исходный файл не существует: %s", sourcePath)
	}

	// Проверяем, не существует ли уже файл в месте назначения
	if _, err := os.Stat(destPath); err == nil {
		return fmt.Errorf("файл назначения уже существует: %s", destPath)
	}

	// Перемещаем файл
	err := os.Rename(sourcePath, destPath)
	if err != nil {
		return fmt.Errorf("ошибка перемещения файла: %w", err)
	}

	return nil
}
