package storage

import (
	"fmt"
	"io"
	"os"
)

func LoadFile(path string) ([]byte, error) {
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
func MoveFile(sourcePath, destPath string, content []byte) (err error) {
	// Проверяем, существует ли исходный файл
	if _, err := os.Stat(sourcePath); os.IsNotExist(err) {
		return fmt.Errorf("исходный файл не существует: %s", sourcePath)
	}

	// Create создает файл или усекает (очищает) существующий
	if content != nil {
		file, err := os.Create(sourcePath)
		if err != nil {
			return fmt.Errorf("Ошибка os.Create: %v", err)
		}

		_, err = file.Write(content)
		if err != nil {
			return fmt.Errorf("Ошибка записи: %v", err)
		}
		file.Close()
	}

	// Перемещаем файл
	err = os.Rename(sourcePath, destPath)
	if err != nil {
		return fmt.Errorf("Oшибка перемещения файла: %w", err)
	}

	return nil
}
