package task

import (
	"context"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
	Create(ctx context.Context, task *Task) error
	GetByID(ctx context.Context, id uuid.UUID) (*Task, error)
	GetByBoard(ctx context.Context, boardID uuid.UUID, filter TaskFilter) ([]Task, error)
	GetByEpic(ctx context.Context, epicID uuid.UUID) ([]Task, error)
	Update(ctx context.Context, task *Task) error
	Delete(ctx context.Context, id uuid.UUID) error
	UpdateStatus(ctx context.Context, id uuid.UUID, status TaskStatus) error
	UpdatePosition(ctx context.Context, id uuid.UUID, position int) error
}

type PgxRepository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) Repository {
	return &PgxRepository{db: db}
}

func (r *PgxRepository) Create(ctx context.Context, task *Task) error {
	query := `
		INSERT INTO tasks (id, board_id, epic_id, title, description, status, position, story_points, due_date, assignee_id, reporter_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
	`
	_, err := r.db.Exec(ctx, query,
		task.ID, task.BoardID, task.EpicID, task.Title, task.Description,
		task.Status, task.Position, task.StoryPoints, task.DueDate,
		task.AssigneeID, task.ReporterID, task.CreatedAt, task.UpdatedAt,
	)
	return err
}

func (r *PgxRepository) GetByID(ctx context.Context, id uuid.UUID) (*Task, error) {
	query := `
		SELECT id, board_id, epic_id, title, description, status, position, story_points, due_date, assignee_id, reporter_id, created_at, updated_at
		FROM tasks WHERE id = $1
	`
	var t Task
	err := r.db.QueryRow(ctx, query, id).Scan(
		&t.ID, &t.BoardID, &t.EpicID, &t.Title, &t.Description,
		&t.Status, &t.Position, &t.StoryPoints, &t.DueDate,
		&t.AssigneeID, &t.ReporterID, &t.CreatedAt, &t.UpdatedAt,
	)
	return &t, err
}

func (r *PgxRepository) GetByBoard(ctx context.Context, boardID uuid.UUID, filter TaskFilter) ([]Task, error) {
	query := `
		SELECT id, board_id, epic_id, title, description, status, position, story_points, due_date, assignee_id, reporter_id, created_at, updated_at
		FROM tasks WHERE board_id = $1
	`
	args := []interface{}{boardID}
	argIdx := 2

	if filter.Status != "" {
		query += " AND status = $" + string(rune('0'+argIdx))
		args = append(args, filter.Status)
		argIdx++
	}
	if filter.AssigneeID != uuid.Nil {
		query += " AND assignee_id = $" + string(rune('0'+argIdx))
		args = append(args, filter.AssigneeID)
		argIdx++
	}
	if filter.EpicID != uuid.Nil {
		query += " AND epic_id = $" + string(rune('0'+argIdx))
		args = append(args, filter.EpicID)
		argIdx++
	}
	if filter.Search != "" {
		query += " AND (title ILIKE $" + string(rune('0'+argIdx)) + " OR description ILIKE $" + string(rune('0'+argIdx)) + ")"
		args = append(args, "%"+filter.Search+"%")
		argIdx++
	}

	query += " ORDER BY position ASC"

	if filter.Limit > 0 {
		query += " LIMIT $" + string(rune('0'+argIdx))
		args = append(args, filter.Limit)
		argIdx++
	}
	if filter.Offset > 0 {
		query += " OFFSET $" + string(rune('0'+argIdx))
		args = append(args, filter.Offset)
	}

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []Task
	for rows.Next() {
		var t Task
		if err := rows.Scan(
			&t.ID, &t.BoardID, &t.EpicID, &t.Title, &t.Description,
			&t.Status, &t.Position, &t.StoryPoints, &t.DueDate,
			&t.AssigneeID, &t.ReporterID, &t.CreatedAt, &t.UpdatedAt,
		); err != nil {
			return nil, err
		}
		tasks = append(tasks, t)
	}
	return tasks, nil
}

func (r *PgxRepository) GetByEpic(ctx context.Context, epicID uuid.UUID) ([]Task, error) {
	query := `
		SELECT id, board_id, epic_id, title, description, status, position, story_points, due_date, assignee_id, reporter_id, created_at, updated_at
		FROM tasks WHERE epic_id = $1 ORDER BY position ASC
	`
	rows, err := r.db.Query(ctx, query, epicID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []Task
	for rows.Next() {
		var t Task
		if err := rows.Scan(
			&t.ID, &t.BoardID, &t.EpicID, &t.Title, &t.Description,
			&t.Status, &t.Position, &t.StoryPoints, &t.DueDate,
			&t.AssigneeID, &t.ReporterID, &t.CreatedAt, &t.UpdatedAt,
		); err != nil {
			return nil, err
		}
		tasks = append(tasks, t)
	}
	return tasks, nil
}

func (r *PgxRepository) Update(ctx context.Context, task *Task) error {
	task.UpdatedAt = time.Now()
	query := `
		UPDATE tasks SET title=$1, description=$2, status=$3, position=$4, story_points=$5, due_date=$6, assignee_id=$7, epic_id=$8, updated_at=$9
		WHERE id=$10
	`
	_, err := r.db.Exec(ctx, query,
		task.Title, task.Description, task.Status, task.Position, task.StoryPoints,
		task.DueDate, task.AssigneeID, task.EpicID, task.UpdatedAt, task.ID,
	)
	return err
}

func (r *PgxRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM tasks WHERE id = $1`
	_, err := r.db.Exec(ctx, query, id)
	return err
}

func (r *PgxRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status TaskStatus) error {
	query := `UPDATE tasks SET status=$1, updated_at=$2 WHERE id=$3`
	_, err := r.db.Exec(ctx, query, status, time.Now(), id)
	return err
}

func (r *PgxRepository) UpdatePosition(ctx context.Context, id uuid.UUID, position int) error {
	query := `UPDATE tasks SET position=$1, updated_at=$2 WHERE id=$3`
	_, err := r.db.Exec(ctx, query, position, time.Now(), id)
	return err
}

func (r *PgxRepository) Search(ctx context.Context, query string, limit int) ([]Task, error) {
	searchPattern := "%" + strings.ToLower(query) + "%"
	sql := `
		SELECT id, board_id, epic_id, title, description, status, position, story_points, due_date, assignee_id, reporter_id, created_at, updated_at
		FROM tasks 
		WHERE LOWER(title) LIKE $1 OR LOWER(description) LIKE $1
		LIMIT $2
	`
	rows, err := r.db.Query(ctx, sql, searchPattern, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []Task
	for rows.Next() {
		var t Task
		if err := rows.Scan(
			&t.ID, &t.BoardID, &t.EpicID, &t.Title, &t.Description,
			&t.Status, &t.Position, &t.StoryPoints, &t.DueDate,
			&t.AssigneeID, &t.ReporterID, &t.CreatedAt, &t.UpdatedAt,
		); err != nil {
			return nil, err
		}
		tasks = append(tasks, t)
	}
	return tasks, nil
}
