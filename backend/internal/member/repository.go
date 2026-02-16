package member

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
	Create(ctx context.Context, member *Member) error
	GetByID(ctx context.Context, id uuid.UUID) (*Member, error)
	GetByEmail(ctx context.Context, email string) (*Member, error)
	GetByCommunity(ctx context.Context, communityID uuid.UUID, limit, offset int) ([]Member, error)
	Update(ctx context.Context, member *Member) error
	Delete(ctx context.Context, id uuid.UUID) error
	UpdateLastLogin(ctx context.Context, id uuid.UUID) error
	CountByCommunity(ctx context.Context, communityID uuid.UUID) (int, error)
}

type PgxRepository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) Repository {
	return &PgxRepository{db: db}
}

func (r *PgxRepository) Create(ctx context.Context, member *Member) error {
	query := `
		INSERT INTO members (id, community_id, email, password_hash, first_name, last_name, date_of_birth, role, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`
	_, err := r.db.Exec(ctx, query,
		member.ID, member.CommunityID, member.Email, member.PasswordHash,
		member.FirstName, member.LastName, member.DateOfBirth, member.Role,
		member.IsActive, member.CreatedAt, member.UpdatedAt,
	)
	return err
}

func (r *PgxRepository) GetByID(ctx context.Context, id uuid.UUID) (*Member, error) {
	query := `
		SELECT id, community_id, email, password_hash, first_name, last_name, date_of_birth, role, is_active, last_login_at, created_at, updated_at
		FROM members WHERE id = $1
	`
	var m Member
	err := r.db.QueryRow(ctx, query, id).Scan(
		&m.ID, &m.CommunityID, &m.Email, &m.PasswordHash, &m.FirstName, &m.LastName,
		&m.DateOfBirth, &m.Role, &m.IsActive, &m.LastLoginAt, &m.CreatedAt, &m.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &m, nil
}

func (r *PgxRepository) GetByEmail(ctx context.Context, email string) (*Member, error) {
	query := `
		SELECT id, community_id, email, password_hash, first_name, last_name, date_of_birth, role, is_active, last_login_at, created_at, updated_at
		FROM members WHERE email = $1
	`
	var m Member
	err := r.db.QueryRow(ctx, query, email).Scan(
		&m.ID, &m.CommunityID, &m.Email, &m.PasswordHash, &m.FirstName, &m.LastName,
		&m.DateOfBirth, &m.Role, &m.IsActive, &m.LastLoginAt, &m.CreatedAt, &m.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &m, nil
}

func (r *PgxRepository) GetByCommunity(ctx context.Context, communityID uuid.UUID, limit, offset int) ([]Member, error) {
	query := `
		SELECT id, community_id, email, password_hash, first_name, last_name, date_of_birth, role, is_active, last_login_at, created_at, updated_at
		FROM members WHERE community_id = $1
		ORDER BY created_at DESC LIMIT $2 OFFSET $3
	`
	rows, err := r.db.Query(ctx, query, communityID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var members []Member
	for rows.Next() {
		var m Member
		if err := rows.Scan(
			&m.ID, &m.CommunityID, &m.Email, &m.PasswordHash, &m.FirstName, &m.LastName,
			&m.DateOfBirth, &m.Role, &m.IsActive, &m.LastLoginAt, &m.CreatedAt, &m.UpdatedAt,
		); err != nil {
			return nil, err
		}
		members = append(members, m)
	}
	return members, nil
}

func (r *PgxRepository) Update(ctx context.Context, member *Member) error {
	member.UpdatedAt = time.Now()
	query := `
		UPDATE members SET first_name=$1, last_name=$2, date_of_birth=$3, role=$4, is_active=$5, updated_at=$6
		WHERE id=$7
	`
	_, err := r.db.Exec(ctx, query, member.FirstName, member.LastName, member.DateOfBirth, member.Role, member.IsActive, member.UpdatedAt, member.ID)
	return err
}

func (r *PgxRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM members WHERE id = $1`
	_, err := r.db.Exec(ctx, query, id)
	return err
}

func (r *PgxRepository) UpdateLastLogin(ctx context.Context, id uuid.UUID) error {
	query := `UPDATE members SET last_login_at = $1 WHERE id = $2`
	_, err := r.db.Exec(ctx, query, time.Now(), id)
	return err
}

func (r *PgxRepository) CountByCommunity(ctx context.Context, communityID uuid.UUID) (int, error) {
	query := `SELECT COUNT(*) FROM members WHERE community_id = $1`
	var count int
	err := r.db.QueryRow(ctx, query, communityID).Scan(&count)
	return count, err
}
