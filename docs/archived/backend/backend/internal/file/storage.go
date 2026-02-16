package file

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/monachy/projek/internal/common"
)

type Storage interface {
	Save(path string, data []byte) error
	Load(path string) ([]byte, error)
	Delete(path string) error
	GetURL(path string) string
}

type LocalStorage struct {
	basePath string
	baseURL  string
}

func NewLocalStorage(basePath, baseURL string) *LocalStorage {
	return &LocalStorage{basePath: basePath, baseURL: baseURL}
}

func (s *LocalStorage) Save(path string, data []byte) error {
	fullPath := filepath.Join(s.basePath, path)
	dir := filepath.Dir(fullPath)

	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	return os.WriteFile(fullPath, data, 0644)
}

func (s *LocalStorage) Load(path string) ([]byte, error) {
	fullPath := filepath.Join(s.basePath, path)
	return os.ReadFile(fullPath)
}

func (s *LocalStorage) Delete(path string) error {
	fullPath := filepath.Join(s.basePath, path)
	return os.Remove(fullPath)
}

func (s *LocalStorage) GetURL(path string) string {
	return fmt.Sprintf("%s/%s", s.baseURL, path)
}

type S3Storage struct {
	bucket  string
	baseURL string
	region  string
}

func NewS3Storage(bucket, region, baseURL string) *S3Storage {
	return &S3Storage{bucket: bucket, region: region, baseURL: baseURL}
}

func (s *S3Storage) Save(path string, data []byte) error {
	return fmt.Errorf("S3 storage not implemented: would save to %s/%s", s.bucket, path)
}

func (s *S3Storage) Load(path string) ([]byte, error) {
	return nil, fmt.Errorf("S3 storage not implemented")
}

func (s *S3Storage) Delete(path string) error {
	return fmt.Errorf("S3 storage not implemented")
}

func (s *S3Storage) GetURL(path string) string {
	return fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", s.bucket, s.region, path)
}

func NewStorageFromConfig(cfg *common.Config) Storage {
	storageType := os.Getenv("STORAGE_TYPE")
	if storageType == "s3" {
		return NewS3Storage(
			os.Getenv("S3_BUCKET"),
			os.Getenv("AWS_REGION"),
			os.Getenv("S3_BASE_URL"),
		)
	}
	return NewLocalStorage(
		"/var/lib/projek/uploads",
		"/uploads",
	)
}

type FileAttachment struct {
	ID               uuid.UUID  `json:"id"`
	ProjectID        uuid.UUID  `json:"project_id"`
	TaskID           *uuid.UUID `json:"task_id,omitempty"`
	UploaderID       uuid.UUID  `json:"uploader_id"`
	Filename         string     `json:"filename"`
	OriginalFilename string     `json:"original_filename"`
	ContentType      string     `json:"content_type"`
	Size             int64      `json:"size"`
	StoragePath      string     `json:"storage_path"`
	StorageType      string     `json:"storage_type"`
	CreatedAt        int64      `json:"created_at"`
}

type CreateFileInput struct {
	ProjectID   uuid.UUID  `json:"project_id"`
	TaskID      *uuid.UUID `json:"task_id"`
	UploaderID  uuid.UUID  `json:"uploader_id" validate:"required"`
	Filename    string     `json:"filename" validate:"required"`
	ContentType string     `json:"content_type" validate:"required"`
	Data        []byte     `json:"-"`
}
