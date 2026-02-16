package file

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
	Create(ctx context.Context, attachment *FileAttachment) error
	GetByID(ctx context.Context, id uuid.UUID) (*FileAttachment, error)
	GetByProject(ctx context.Context, projectID uuid.UUID) ([]FileAttachment, error)
	GetByTask(ctx context.Context, taskID uuid.UUID) ([]FileAttachment, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type PgxRepository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) Repository {
	return &PgxRepository{db: db}
}

func (r *PgxRepository) Create(ctx context.Context, attachment *FileAttachment) error {
	query := `
		INSERT INTO file_attachments (id, project_id, task_id, uploader_id, filename, original_filename, content_type, size, storage_path, storage_type, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`
	_, err := r.db.Exec(ctx, query,
		attachment.ID, attachment.ProjectID, attachment.TaskID, attachment.UploaderID,
		attachment.Filename, attachment.OriginalFilename, attachment.ContentType,
		attachment.Size, attachment.StoragePath, attachment.StorageType, attachment.CreatedAt,
	)
	return err
}

func (r *PgxRepository) GetByID(ctx context.Context, id uuid.UUID) (*FileAttachment, error) {
	query := `
		SELECT id, project_id, task_id, uploader_id, filename, original_filename, content_type, size, storage_path, storage_type, created_at
		FROM file_attachments WHERE id = $1
	`
	var a FileAttachment
	var taskIDPtr, uploaderIDPtr *uuid.UUID
	err := r.db.QueryRow(ctx, query, id).Scan(
		&a.ID, &a.ProjectID, &taskIDPtr, &uploaderIDPtr,
		&a.Filename, &a.OriginalFilename, &a.ContentType,
		&a.Size, &a.StoragePath, &a.StorageType, &a.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	if taskIDPtr != nil {
		a.TaskID = taskIDPtr
	}
	if uploaderIDPtr != nil {
		a.UploaderID = *uploaderIDPtr
	}
	return &a, nil
}

func (r *PgxRepository) GetByProject(ctx context.Context, projectID uuid.UUID) ([]FileAttachment, error) {
	query := `
		SELECT id, project_id, task_id, uploader_id, filename, original_filename, content_type, size, storage_path, storage_type, created_at
		FROM file_attachments WHERE project_id = $1 ORDER BY created_at DESC
	`
	rows, err := r.db.Query(ctx, query, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var attachments []FileAttachment
	for rows.Next() {
		var a FileAttachment
		var taskIDPtr, uploaderIDPtr *uuid.UUID
		if err := rows.Scan(
			&a.ID, &a.ProjectID, &taskIDPtr, &uploaderIDPtr,
			&a.Filename, &a.OriginalFilename, &a.ContentType,
			&a.Size, &a.StoragePath, &a.StorageType, &a.CreatedAt,
		); err != nil {
			return nil, err
		}
		if taskIDPtr != nil {
			a.TaskID = taskIDPtr
		}
		if uploaderIDPtr != nil {
			a.UploaderID = *uploaderIDPtr
		}
		attachments = append(attachments, a)
	}
	return attachments, nil
}

func (r *PgxRepository) GetByTask(ctx context.Context, taskID uuid.UUID) ([]FileAttachment, error) {
	query := `
		SELECT id, project_id, task_id, uploader_id, filename, original_filename, content_type, size, storage_path, storage_type, created_at
		FROM file_attachments WHERE task_id = $1 ORDER BY created_at DESC
	`
	rows, err := r.db.Query(ctx, query, taskID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var attachments []FileAttachment
	for rows.Next() {
		var a FileAttachment
		var taskIDPtr, uploaderIDPtr *uuid.UUID
		if err := rows.Scan(
			&a.ID, &a.ProjectID, &taskIDPtr, &uploaderIDPtr,
			&a.Filename, &a.OriginalFilename, &a.ContentType,
			&a.Size, &a.StoragePath, &a.StorageType, &a.CreatedAt,
		); err != nil {
			return nil, err
		}
		if taskIDPtr != nil {
			a.TaskID = taskIDPtr
		}
		if uploaderIDPtr != nil {
			a.UploaderID = *uploaderIDPtr
		}
		attachments = append(attachments, a)
	}
	return attachments, nil
}

func (r *PgxRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM file_attachments WHERE id = $1`
	_, err := r.db.Exec(ctx, query, id)
	return err
}

type Service struct {
	repo    Repository
	storage Storage
}

func NewService(repo Repository, storage Storage) *Service {
	return &Service{repo: repo, storage: storage}
}

func (s *Service) Upload(ctx context.Context, input CreateFileInput) (*FileAttachment, error) {
	filename := uuid.New().String() + "-" + input.Filename
	storagePath := input.ProjectID.String() + "/" + filename

	if err := s.storage.Save(storagePath, input.Data); err != nil {
		return nil, err
	}

	attachment := &FileAttachment{
		ID:               uuid.New(),
		ProjectID:        input.ProjectID,
		TaskID:           input.TaskID,
		UploaderID:       input.UploaderID,
		Filename:         filename,
		OriginalFilename: input.Filename,
		ContentType:      input.ContentType,
		Size:             int64(len(input.Data)),
		StoragePath:      storagePath,
		StorageType:      "local",
		CreatedAt:        time.Now().Unix(),
	}

	if err := s.repo.Create(ctx, attachment); err != nil {
		s.storage.Delete(storagePath)
		return nil, err
	}

	return attachment, nil
}

func (s *Service) Download(ctx context.Context, id uuid.UUID) ([]byte, error) {
	attachment, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return s.storage.Load(attachment.StoragePath)
}

func (s *Service) GetByProject(ctx context.Context, projectID uuid.UUID) ([]FileAttachment, error) {
	return s.repo.GetByProject(ctx, projectID)
}

func (s *Service) GetByTask(ctx context.Context, taskID uuid.UUID) ([]FileAttachment, error) {
	return s.repo.GetByTask(ctx, taskID)
}

func (s *Service) Delete(ctx context.Context, id uuid.UUID) error {
	attachment, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	s.storage.Delete(attachment.StoragePath)
	return s.repo.Delete(ctx, id)
}

func (s *Service) GetURL(ctx context.Context, id uuid.UUID) (string, error) {
	attachment, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return "", err
	}

	return s.storage.GetURL(attachment.StoragePath), nil
}
