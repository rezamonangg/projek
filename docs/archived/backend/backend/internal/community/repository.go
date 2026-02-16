package community

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
	Create(ctx context.Context, community *Community) error
	GetByID(ctx context.Context, id uuid.UUID) (*Community, error)
	GetBySlug(ctx context.Context, slug string) (*Community, error)
	Update(ctx context.Context, community *Community) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type PgxRepository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) Repository {
	return &PgxRepository{db: db}
}

func (r *PgxRepository) Create(ctx context.Context, community *Community) error {
	query := `
		INSERT INTO communities (id, name, slug, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
	`
	_, err := r.db.Exec(ctx, query, community.ID, community.Name, community.Slug, community.CreatedAt, community.UpdatedAt)
	return err
}

func (r *PgxRepository) GetByID(ctx context.Context, id uuid.UUID) (*Community, error) {
	query := `SELECT id, name, slug, created_at, updated_at FROM communities WHERE id = $1`
	var c Community
	err := r.db.QueryRow(ctx, query, id).Scan(&c.ID, &c.Name, &c.Slug, &c.CreatedAt, &c.UpdatedAt)
	return &c, err
}

func (r *PgxRepository) GetBySlug(ctx context.Context, slug string) (*Community, error) {
	query := `SELECT id, name, slug, created_at, updated_at FROM communities WHERE slug = $1`
	var c Community
	err := r.db.QueryRow(ctx, query, slug).Scan(&c.ID, &c.Name, &c.Slug, &c.CreatedAt, &c.UpdatedAt)
	return &c, err
}

func (r *PgxRepository) Update(ctx context.Context, community *Community) error {
	query := `UPDATE communities SET name=$1, slug=$2, updated_at=$3 WHERE id=$4`
	_, err := r.db.Exec(ctx, query, community.Name, community.Slug, community.UpdatedAt, community.ID)
	return err
}

func (r *PgxRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM communities WHERE id = $1`
	_, err := r.db.Exec(ctx, query, id)
	return err
}
