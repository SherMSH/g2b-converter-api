package handlers

import (
	"converterapi/internal/config"
	"converterapi/pkg/logger"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

func PutConvFile(c *gin.Context) {
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Не удалось получить файл: " + err.Error(),
		})
		return
	}
	defer file.Close()

	// Проверка MIME типа
	buffer := make([]byte, 512)
	if _, err := file.Read(buffer); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Не удалось прочитать файл",
		})
		return
	}

	// Возвращаем указатель в начало файла
	if _, err := file.Seek(0, io.SeekStart); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Ошибка при чтении файла",
		})
		return
	}
	mimeType := http.DetectContentType(buffer)
	if mimeType != "text/xml; charset=utf-8" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":         "Тип файла не поддерживается",
			"detected_type": mimeType,
		})
		return
	}

	// Создаем директорию для загрузок, если её нет
	uploadDir := filepath.Join(config.Config.App.Storage.Basepath, config.Config.App.Storage.In)
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Не удалось создать директорию: " + err.Error(),
		})
		return
	}

	filePath := filepath.Join(uploadDir, header.Filename) // Формируем путь для сохранения файла
	dst, err := os.Create(filePath)                       // Создаем файл на диске
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Не удалось создать файл: " + err.Error(),
		})
		return
	}
	defer dst.Close()
	if _, err := io.Copy(dst, file); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Не удалось сохранить файл: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "Файл успешно загружен",
		"filename": header.Filename,
		"size":     header.Size,
	})
}

func GetConvFile(c *gin.Context) {
	// Получаем имя файла из параметра URL
	filename := c.Param("filename")
	logger.Infof("Uploading file: %s", filename)
	filename = filepath.Clean(filename)
	if strings.Contains(filename, "..") || len(filename) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Некорректное имя файла",
		})
		return
	}
	// Путь к файлу
	filePath := filepath.Join(config.Config.App.Storage.Basepath, config.Config.App.Storage.Out, filename)
	// Открываем файл
	file, err := os.Open(filePath)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Файл не найден",
		})
		return
	}
	defer file.Close()

	// Получаем информацию о файле
	stat, err := file.Stat()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Ошибка при получении информации о файле",
		})
		return
	}

	// Устанавливаем заголовки
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Length", fmt.Sprintf("%d", stat.Size()))

	// Потоковая передача файла
	c.Stream(func(w io.Writer) bool {
		// Копируем файл в ответ по частям
		buffer := make([]byte, 8192) // 8KB буфер
		for {
			n, err := file.Read(buffer)
			if n > 0 {
				if _, writeErr := w.Write(buffer[:n]); writeErr != nil {
					return false
				}
			}
			if err == io.EOF {
				break
			}
			if err != nil {
				return false
			}
		}
		return false
	})
}
