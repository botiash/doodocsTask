package usecase

import (
	"archive/zip"
	"doodocsProg/internal/models"
	"errors"
	"io/ioutil"

	"github.com/gabriel-vasile/mimetype"
)

// GetArchiveInfo возвращает информацию о файле архива.
func GetArchiveInfo(filePath string) (*models.ArchiveInfo, error) {
	archiveInfo := &models.ArchiveInfo{}

	// Открываем архивный файл.
	r, err := zip.OpenReader(filePath)
	if err != nil {
		return nil, err
	}
	defer r.Close()

	if len(r.File) == 0 {
		return nil, errors.New("archive is empty")
	}

	archiveInfo.Filename = r.File[0].Name
	archiveInfo.ArchiveSize = float64(r.File[0].UncompressedSize64)
	archiveInfo.TotalSize = 0
	archiveInfo.TotalFiles = float64(len(r.File))
	archiveInfo.Files = []models.FileInfo{}
	
	// Перебираем файлы в архиве и собираем информацию о них.
	for _, f := range r.File {
		fileInfo := models.FileInfo{
			FilePath: f.Name,
			Size:     float64(f.UncompressedSize64),
			MimeType: GetMimeType(f.Name, f),
		}
		archiveInfo.TotalSize += fileInfo.Size
		archiveInfo.Files = append(archiveInfo.Files, fileInfo)
	}

	return archiveInfo, nil
}

func GetMimeType(filename string, f *zip.File) string {
	// Открываем файл в архиве.
	file, err := f.Open()
	if err != nil {
		return "application/octet-stream"
	}
	defer file.Close()

	// Читаем содержимое файла и определяем MIME-тип на основе его содержимого.
	data, err := ioutil.ReadAll(file)
	if err != nil {
		return "application/octet-stream"
	}

	// Определяем MIME-тип на основе содержимого файла.
	mime := mimetype.Detect(data)

	return mime.String()
}
