package project

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type EpicRepository interface {
	Create(ctx context.Context, epic *Epic) error
	GetByID(ctx context.Context, id uuid.UUID) (*Epic, error)
	GetByProject(ctx context.Context, projectID uuid.UUID) ([]Epic, error)
	Update(ctx context.Context, epic *Epic) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type BoardRepository interface {
	Create(ctx context.Context, board *Board) error
	GetByID(ctx context.Context, id uuid.UUID) (*Board, error)
	GetByProject(ctx context.Context, projectID uuid.UUID) ([]Board, error)
	Update(ctx context.Context, board *Board) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type PgxEpicRepository struct {
	db *pgxpool.Pool
}

func NewEpicRepository(db *pgxpool.Pool) EpicRepository {
	return &PgxEpicRepository{db: db}
}

func (r *PgxEpicRepository) Create(ctx context.Context, epic *Epic) error {
	query := `
		INSERT INTO epics (id, project_id, name, description, start_date, end_date, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`
	_, err := r.db.Exec(ctx, query,
		epic.ID, epic.ProjectID, epic.Name, epic.Description,
		epic.StartDate, epic.EndDate, epic.CreatedAt, epic.UpdatedAt,
	)
	return err
}

func (r *PgxEpicRepository) GetByID(ctx context.Context, id uuid.UUID) (*Epic, error) {
	query := `SELECT id, project_id, name, description, start_date, end_date, created_at, updated_at FROM epics WHERE id = $1`
	var e Epic
	err := r.db.QueryRow(ctx, query, id).Scan(
		&e.ID, &e.ProjectID, &e.Name, &e.Description,
		&e.StartDate, &e.EndDate, &e.CreatedAt, &e.UpdatedAt,
	)
	return &e, err
}

func (r *PgxEpicRepository) GetByProject(ctx context.Context, projectID uuid.UUID) ([]Epic, error) {
	query := `
		SELECT id, project_id, name, description, start_date, end_date, created_at, updated_at
		FROM epics WHERE project_id = $1 ORDER BY created_at DESC
	`
	rows, err := r.db.Query(ctx, query, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var epics []Epic
	for rows.Next() {
		var e Epic
		if err := rows.Scan(
			&e.ID, &e.ProjectID, &e.Name, &e.Description,
			&e.StartDate, &e.EndDate, &e.CreatedAt, &e.UpdatedAt,
		); err != nil {
			return nil, err
		}
		epics = append(epics, e)
	}
	return epics, nil
}

func (r *PgxEpicRepository) Update(ctx context.Context, epic *Epic) error {
	epic.UpdatedAt = time.Now()
	query := `UPDATE epics SET name=$1, description=$2, start_date=$3, end_date=$4, updated_at=$5 WHERE id=$6`
	_, err := r.db.Exec(ctx, query, epic.Name, epic.Description, epic.StartDate, epic.EndDate, epic.UpdatedAt, epic.ID)
	return err
}

func (r *PgxEpicRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM epics WHERE id = $1`
	_, err := r.db.Exec(ctx, query, id)
	return err
}

type PgxBoardRepository struct {
	db *pgxpool.Pool
}

func NewBoardRepository(db *pgxpool.Pool) BoardRepository {
	return &PgxBoardRepository{db: db}
}

func (r *PgxBoardRepository) Create(ctx context.Context, board *Board) error {
	query := `
		INSERT INTO boards (id, project_id, name, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
	`
	_, err := r.db.Exec(ctx, query, board.ID, board.ProjectID, board.Name, board.CreatedAt, board.UpdatedAt)
	return err
}

func (r *PgxBoardRepository) GetByID(ctx context.Context, id uuid.UUID) (*Board, error) {
	query := `SELECT id, project_id, name, created_at, updated_at FROM boards WHERE id = $1`
	var b Board
	err := r.db.QueryRow(ctx, query, id).Scan(&b.ID, &b.ProjectID, &b.Name, &b.CreatedAt, &b.UpdatedAt)
	return &b, err
}

func (r *PgxBoardRepository) GetByProject(ctx context.Context, projectID uuid.UUID) ([]Board, error) {
	query := `
		SELECT id, project_id, name, created_at, updated_at
		FROM boards WHERE project_id = $1 ORDER BY created_at DESC
	`
	rows, err := r.db.Query(ctx, query, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var boards []Board
	for rows.Next() {
		var b Board
		if err := rows.Scan(&b.ID, &b.ProjectID, &b.Name, &b.CreatedAt, &b.UpdatedAt); err != nil {
			return nil, err
		}
		boards = append(boards, b)
	}
	return boards, nil
}

func (r *PgxBoardRepository) Update(ctx context.Context, board *Board) error {
	board.UpdatedAt = time.Now()
	query := `UPDATE boards SET name=$1, updated_at=$2 WHERE id=$3`
	_, err := r.db.Exec(ctx, query, board.Name, board.UpdatedAt, board.ID)
	return err
}

func (r *PgxBoardRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM boards WHERE id = $1`
	_, err := r.db.Exec(ctx, query, id)
	return err
}
