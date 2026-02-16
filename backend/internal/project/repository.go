package project

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
	Create(ctx context.Context, project *Project) error
	GetByID(ctx context.Context, id uuid.UUID) (*Project, error)
	GetByCommunity(ctx context.Context, communityID uuid.UUID, limit, offset int) ([]Project, error)
	Update(ctx context.Context, project *Project) error
	Delete(ctx context.Context, id uuid.UUID) error
	Archive(ctx context.Context, id uuid.UUID) error
}

type PgxRepository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) Repository {
	return &PgxRepository{db: db}
}

func (r *PgxRepository) Create(ctx context.Context, project *Project) error {
	query := `
		INSERT INTO projects (id, community_id, name, description, key, is_archived, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`
	_, err := r.db.Exec(ctx, query,
		project.ID, project.CommunityID, project.Name, project.Description,
		project.Key, project.IsArchived, project.CreatedAt, project.UpdatedAt,
	)
	return err
}

func (r *PgxRepository) GetByID(ctx context.Context, id uuid.UUID) (*Project, error) {
	query := `SELECT id, community_id, name, description, key, is_archived, created_at, updated_at FROM projects WHERE id = $1`
	var p Project
	err := r.db.QueryRow(ctx, query, id).Scan(
		&p.ID, &p.CommunityID, &p.Name, &p.Description, &p.Key,
		&p.IsArchived, &p.CreatedAt, &p.UpdatedAt,
	)
	return &p, err
}

func (r *PgxRepository) GetByCommunity(ctx context.Context, communityID uuid.UUID, limit, offset int) ([]Project, error) {
	query := `
		SELECT id, community_id, name, description, key, is_archived, created_at, updated_at
		FROM projects WHERE community_id = $1 AND is_archived = false
		ORDER BY created_at DESC LIMIT $2 OFFSET $3
	`
	rows, err := r.db.Query(ctx, query, communityID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var projects []Project
	for rows.Next() {
		var p Project
		if err := rows.Scan(
			&p.ID, &p.CommunityID, &p.Name, &p.Description, &p.Key,
			&p.IsArchived, &p.CreatedAt, &p.UpdatedAt,
		); err != nil {
			return nil, err
		}
		projects = append(projects, p)
	}
	return projects, nil
}

func (r *PgxRepository) Update(ctx context.Context, project *Project) error {
	query := `
		UPDATE projects SET name=$1, description=$2, key=$3, is_archived=$4, updated_at=$5
		WHERE id=$6
	`
	_, err := r.db.Exec(ctx, query, project.Name, project.Description, project.Key, project.IsArchived, project.UpdatedAt, project.ID)
	return err
}

func (r *PgxRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM projects WHERE id = $1`
	_, err := r.db.Exec(ctx, query, id)
	return err
}

func (r *PgxRepository) Archive(ctx context.Context, id uuid.UUID) error {
	query := `UPDATE projects SET is_archived = true, updated_at = NOW() WHERE id = $1`
	_, err := r.db.Exec(ctx, query, id)
	return err
}
