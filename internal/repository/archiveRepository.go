package repository

import (
	"archive/zip"
	"doodocsProg/internal/models"
	"fmt"
	"os"
	"path/filepath"

	"github.com/gabriel-vasile/mimetype"
)

type ArchiveRepository struct {
	StoragePath string
}

func NewArchiveRepository(storagePath string) *ArchiveRepository {
	return &ArchiveRepository{
		StoragePath: storagePath,
	}
}

func (r *ArchiveRepository) GetArchiveInfo(filePath string) (*models.ArchiveInfo, error) {
	archivePath := filepath.Join(r.StoragePath, filePath)

	// Открываем архивный файл.
	reader, err := zip.OpenReader(archivePath)
	if err != nil {
		return nil, handleError(err, "Failed to open archive file")
	}
	defer reader.Close()

	archiveInfo := &models.ArchiveInfo{
		Filename:    filePath,
		ArchiveSize: float64(reader.File[0].UncompressedSize64),
		TotalSize:   0,
		TotalFiles:  float64(len(reader.File)),
		Files:       []models.FileInfo{},
	}

	for _, file := range reader.File {
		fileInfo := models.FileInfo{
			FilePath: file.Name,
			Size:     float64(file.UncompressedSize64),
			MimeType: getMimeType(file.Name),
		}
		archiveInfo.TotalSize += fileInfo.Size
		archiveInfo.Files = append(archiveInfo.Files, fileInfo)
	}

	return archiveInfo, nil
}

func (r *ArchiveRepository) SaveArchive(archiveName string, data []byte) error {
	archivePath := filepath.Join(r.StoragePath, archiveName)
	if err := os.MkdirAll(r.StoragePath, 0755); err != nil {
		return handleError(err, "Failed to create storage directory")
	}

	file, err := os.Create(archivePath)
	if err != nil {
		return handleError(err, "Failed to create archive file")
	}
	defer file.Close()

	_, err = file.Write(data)
	if err != nil {
		return handleError(err, "Failed to write archive data")
	}

	return nil
}

func handleError(err error, message string) error {
	if err != nil {
		return fmt.Errorf("%s: %w", message, err)
	}
	return nil
}

func getMimeType(filename string) string {
	file, err := os.Open(filename)
	if err != nil {
		return "application/octet-stream"
	}
	defer file.Close()

	buffer := make([]byte, 512) // Читаем первые 512 байт файла для определения MIME-типа.
	_, err = file.Read(buffer)
	if err != nil {
		return "application/octet-stream" // MIME-тип по умолчанию, если не удалось прочитать файл
	}

	// Используем библиотеку mimetype для определения MIME-типа.
	mimeType := mimetype.Detect(buffer)
	return mimeType.String()
}
