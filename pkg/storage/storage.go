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
	err = copyAndRemove(sourcePath, destPath)
	if err != nil {
		return fmt.Errorf("Oшибка перемещения файла: %w", err)
	}

	return nil
}

// copyAndRemove копирует файл и удаляет исходник
func copyAndRemove(sourcePath, destPath string) error {
	// 1. Копируем файл с сохранением прав доступа
	err := copyFile(sourcePath, destPath)
	if err != nil {
		return fmt.Errorf("копирование файла: %w", err)
	}

	// 2. Удаляем исходный файл
	err = os.Remove(sourcePath)
	if err != nil {
		// Пытаемся удалить скопированный файл, чтобы не оставить мусор
		os.Remove(destPath)
		return fmt.Errorf("удаление исходного файла: %w", err)
	}

	return nil
}

// copyFile копирует содержимое и права доступа
func copyFile(sourcePath, destPath string) error {
	// Открываем исходный файл
	src, err := os.Open(sourcePath)
	if err != nil {
		return err
	}
	defer src.Close()

	// Получаем информацию о файле для прав доступа
	srcInfo, err := src.Stat()
	if err != nil {
		return err
	}

	// Создаём целевой файл с теми же правами
	dest, err := os.OpenFile(destPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, srcInfo.Mode())
	if err != nil {
		return err
	}
	defer dest.Close()

	// Копируем содержимое
	_, err = io.Copy(dest, src)
	if err != nil {
		return err
	}

	// Гарантируем запись на диск
	err = dest.Sync()
	if err != nil {
		return err
	}

	return nil
}
