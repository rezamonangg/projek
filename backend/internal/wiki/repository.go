package wiki

import (
	"context"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/monachy/projek/internal/common"
)

type Repository interface {
	Create(ctx context.Context, page *WikiPage) error
	GetByID(ctx context.Context, id uuid.UUID) (*WikiPage, error)
	GetBySlug(ctx context.Context, projectID uuid.UUID, slug string) (*WikiPage, error)
	GetByProject(ctx context.Context, projectID uuid.UUID) ([]WikiPage, error)
	GetChildren(ctx context.Context, parentID uuid.UUID) ([]WikiPage, error)
	Update(ctx context.Context, page *WikiPage) error
	Delete(ctx context.Context, id uuid.UUID) error
	CreateVersion(ctx context.Context, version *WikiPageVersion) error
	GetVersions(ctx context.Context, pageID uuid.UUID) ([]WikiPageVersion, error)
	GetVersion(ctx context.Context, pageID uuid.UUID, version int) (*WikiPageVersion, error)
}

type PgxRepository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) Repository {
	return &PgxRepository{db: db}
}

func (r *PgxRepository) Create(ctx context.Context, page *WikiPage) error {
	query := `
		INSERT INTO wiki_pages (id, project_id, parent_id, title, slug, content, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`
	_, err := r.db.Exec(ctx, query,
		page.ID, page.ProjectID, page.ParentID, page.Title, page.Slug,
		page.Content, page.CreatedAt, page.UpdatedAt,
	)
	return err
}

func (r *PgxRepository) GetByID(ctx context.Context, id uuid.UUID) (*WikiPage, error) {
	query := `SELECT id, project_id, parent_id, title, slug, content, created_at, updated_at FROM wiki_pages WHERE id = $1`
	var p WikiPage
	err := r.db.QueryRow(ctx, query, id).Scan(
		&p.ID, &p.ProjectID, &p.ParentID, &p.Title, &p.Slug, &p.Content, &p.CreatedAt, &p.UpdatedAt,
	)
	return &p, err
}

func (r *PgxRepository) GetBySlug(ctx context.Context, projectID uuid.UUID, slug string) (*WikiPage, error) {
	query := `SELECT id, project_id, parent_id, title, slug, content, created_at, updated_at FROM wiki_pages WHERE project_id = $1 AND slug = $2`
	var p WikiPage
	err := r.db.QueryRow(ctx, query, projectID, slug).Scan(
		&p.ID, &p.ProjectID, &p.ParentID, &p.Title, &p.Slug, &p.Content, &p.CreatedAt, &p.UpdatedAt,
	)
	return &p, err
}

func (r *PgxRepository) GetByProject(ctx context.Context, projectID uuid.UUID) ([]WikiPage, error) {
	query := `
		SELECT id, project_id, parent_id, title, slug, content, created_at, updated_at
		FROM wiki_pages WHERE project_id = $1 ORDER BY title
	`
	rows, err := r.db.Query(ctx, query, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var pages []WikiPage
	for rows.Next() {
		var p WikiPage
		if err := rows.Scan(&p.ID, &p.ProjectID, &p.ParentID, &p.Title, &p.Slug, &p.Content, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, err
		}
		pages = append(pages, p)
	}
	return pages, nil
}

func (r *PgxRepository) GetChildren(ctx context.Context, parentID uuid.UUID) ([]WikiPage, error) {
	query := `
		SELECT id, project_id, parent_id, title, slug, content, created_at, updated_at
		FROM wiki_pages WHERE parent_id = $1 ORDER BY title
	`
	rows, err := r.db.Query(ctx, query, parentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var pages []WikiPage
	for rows.Next() {
		var p WikiPage
		if err := rows.Scan(&p.ID, &p.ProjectID, &p.ParentID, &p.Title, &p.Slug, &p.Content, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, err
		}
		pages = append(pages, p)
	}
	return pages, nil
}

func (r *PgxRepository) Update(ctx context.Context, page *WikiPage) error {
	page.UpdatedAt = time.Now()
	query := `UPDATE wiki_pages SET title=$1, slug=$2, content=$3, parent_id=$4, updated_at=$5 WHERE id=$6`
	_, err := r.db.Exec(ctx, query, page.Title, page.Slug, page.Content, page.ParentID, page.UpdatedAt, page.ID)
	return err
}

func (r *PgxRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM wiki_pages WHERE id = $1`
	_, err := r.db.Exec(ctx, query, id)
	return err
}

func (r *PgxRepository) CreateVersion(ctx context.Context, version *WikiPageVersion) error {
	query := `
		INSERT INTO wiki_page_versions (id, wiki_page_id, content, version, created_at, created_by)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	_, err := r.db.Exec(ctx, query,
		version.ID, version.WikiPageID, version.Content, version.Version, version.CreatedAt, version.CreatedBy,
	)
	return err
}

func (r *PgxRepository) GetVersions(ctx context.Context, pageID uuid.UUID) ([]WikiPageVersion, error) {
	query := `
		SELECT id, wiki_page_id, content, version, created_at, created_by
		FROM wiki_page_versions WHERE wiki_page_id = $1 ORDER BY version DESC
	`
	rows, err := r.db.Query(ctx, query, pageID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var versions []WikiPageVersion
	for rows.Next() {
		var v WikiPageVersion
		if err := rows.Scan(&v.ID, &v.WikiPageID, &v.Content, &v.Version, &v.CreatedAt, &v.CreatedBy); err != nil {
			return nil, err
		}
		versions = append(versions, v)
	}
	return versions, nil
}

func (r *PgxRepository) GetVersion(ctx context.Context, pageID uuid.UUID, version int) (*WikiPageVersion, error) {
	query := `SELECT id, wiki_page_id, content, version, created_at, created_by FROM wiki_page_versions WHERE wiki_page_id = $1 AND version = $2`
	var v WikiPageVersion
	err := r.db.QueryRow(ctx, query, pageID, version).Scan(&v.ID, &v.WikiPageID, &v.Content, &v.Version, &v.CreatedAt, &v.CreatedBy)
	return &v, err
}

func GenerateSlug(title string) string {
	slug := strings.ToLower(title)
	slug = strings.ReplaceAll(slug, " ", "-")
	slug = strings.Map(func(r rune) rune {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '-' {
			return r
		}
		return -1
	}, slug)
	return slug
}

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Create(ctx context.Context, projectID uuid.UUID, input CreateWikiPageInput) (*WikiPage, error) {
	slug := GenerateSlug(input.Title)

	page := &WikiPage{
		ID:        uuid.New(),
		ProjectID: projectID,
		ParentID:  input.ParentID,
		Title:     input.Title,
		Slug:      slug,
		Content:   input.Content,
		CreatedAt: common.Now(),
		UpdatedAt: common.Now(),
	}

	if err := s.repo.Create(ctx, page); err != nil {
		return nil, err
	}

	return page, nil
}

func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (*WikiPage, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *Service) GetBySlug(ctx context.Context, projectID uuid.UUID, slug string) (*WikiPage, error) {
	return s.repo.GetBySlug(ctx, projectID, slug)
}

func (s *Service) GetByProject(ctx context.Context, projectID uuid.UUID) ([]WikiPage, error) {
	return s.repo.GetByProject(ctx, projectID)
}

func (s *Service) GetChildren(ctx context.Context, parentID uuid.UUID) ([]WikiPage, error) {
	return s.repo.GetChildren(ctx, parentID)
}

func (s *Service) Update(ctx context.Context, id uuid.UUID, input UpdateWikiPageInput, userID uuid.UUID) (*WikiPage, error) {
	page, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	oldContent := page.Content
	newContent := input.Content

	if input.Title != "" {
		page.Title = input.Title
		page.Slug = GenerateSlug(input.Title)
	}
	if input.Content != "" {
		page.Content = input.Content
	}
	if input.ParentID != nil {
		page.ParentID = input.ParentID
	}

	if err := s.repo.Update(ctx, page); err != nil {
		return nil, err
	}

	if oldContent != newContent {
		versions, _ := s.repo.GetVersions(ctx, id)
		newVersion := len(versions) + 1

		version := &WikiPageVersion{
			ID:         uuid.New(),
			WikiPageID: id,
			Content:    oldContent,
			Version:    newVersion,
			CreatedAt:  common.Now(),
			CreatedBy:  userID,
		}
		s.repo.CreateVersion(ctx, version)
	}

	return page, nil
}

func (s *Service) Delete(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}

func (s *Service) GetVersions(ctx context.Context, pageID uuid.UUID) ([]WikiPageVersion, error) {
	return s.repo.GetVersions(ctx, pageID)
}
