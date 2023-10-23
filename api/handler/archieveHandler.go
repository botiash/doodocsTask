package handler

import (
	"archive/zip"
	"doodocsProg/internal/repository"
	"doodocsProg/internal/usecase"
	"io"
	"io/ioutil"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

func GetArchiveInfo(c *gin.Context) {
	file, err := c.FormFile("file")
	if file == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file provided"})
		return
	}

	if !isValidZipFile(file.Filename) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Not zip file"})
		return
	}
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to retrieve the file"})
		return
	}

	tempFileName := file.Filename

	tempFile := filepath.Join("tmp", tempFileName)
	if err := c.SaveUploadedFile(file, tempFile); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save the file"})
		return
	}

	data, err := ioutil.ReadFile(tempFile)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read archive data"})
		return
	}

	// Создаем экземпляр репозитория ArchiveRepository и вызываем его метод SaveArchive.
	archiveRepo := repository.NewArchiveRepository("temp/")
	err = archiveRepo.SaveArchive(tempFileName, data)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save the archive"})
		return
	}

	// Создаем экземпляр UseCase и вызываем его метод для получения информации о файле.
	archiveInfo, err := usecase.GetArchiveInfo(tempFile)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get archive information"})
		return
	}

	c.JSON(http.StatusOK, archiveInfo)
}

func CreateArchieve(c *gin.Context) {
	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse form data"})
		return
	}

	files := form.File["file[]"]
	if len(files) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No files provided"})
		return
	}
	for _, fileHeader := range files {
		mimeType := mime.TypeByExtension(filepath.Ext(fileHeader.Filename))
		if !isValidMimeType(mimeType) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file format"})
			return
		}
	}

	// Создаем временное имя для архива.
	tempFileName := "my_archive.zip"

	// Открываем временный файл для записи архива.
	tempFilePath := filepath.Join("tmp", tempFileName)
	archiveFile, err := os.Create(tempFilePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create archive"})
		return
	}
	defer archiveFile.Close()

	// Создаем ZIP-архив и добавляем файлы в него.
	archiveWriter := zip.NewWriter(archiveFile)
	for _, fileHeader := range files {
		file, err := fileHeader.Open()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open file"})
			return
		}
		defer file.Close()

		// Создаем файл в архиве с именем, равным оригинальному имени файла.
		archiveFileInZip, err := archiveWriter.Create(fileHeader.Filename)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add file to archive"})
			return
		}

		// Копируем содержимое файла в архив.
		_, err = io.Copy(archiveFileInZip, file)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to write file to archive"})
			return
		}
	}

	// Закрываем ZIP-архив.
	archiveWriter.Close()

	// Отправляем ZIP-архив клиенту.
	c.File(tempFilePath)
}

func isValidMimeType(mimeType string) bool {
	validMimeTypes := []string{
		"application/vnd.openxmlformats-officedocument.wordprocessingml.document",
		"application/xml",
		"image/jpeg",
		"image/png",
	}

	for _, validType := range validMimeTypes {
		if mimeType == validType {
			return true
		}
	}

	return false
}

func isValidZipFile(filename string) bool {
	return strings.HasSuffix(filename, ".zip")
}
