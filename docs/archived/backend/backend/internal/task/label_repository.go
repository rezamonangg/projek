package task

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type LabelRepository interface {
	Create(ctx context.Context, label *Label) error
	GetByID(ctx context.Context, id uuid.UUID) (*Label, error)
	GetByProject(ctx context.Context, projectID uuid.UUID) ([]Label, error)
	Update(ctx context.Context, label *Label) error
	Delete(ctx context.Context, id uuid.UUID) error
	AddToTask(ctx context.Context, taskID, labelID uuid.UUID) error
	RemoveFromTask(ctx context.Context, taskID, labelID uuid.UUID) error
	GetByTask(ctx context.Context, taskID uuid.UUID) ([]Label, error)
}

type PgxLabelRepository struct {
	db *pgxpool.Pool
}

func NewLabelRepository(db *pgxpool.Pool) LabelRepository {
	return &PgxLabelRepository{db: db}
}

func (r *PgxLabelRepository) Create(ctx context.Context, label *Label) error {
	query := `
		INSERT INTO labels (id, project_id, name, color, created_at)
		VALUES ($1, $2, $3, $4, $5)
	`
	_, err := r.db.Exec(ctx, query, label.ID, label.ProjectID, label.Name, label.Color, label.CreatedAt)
	return err
}

func (r *PgxLabelRepository) GetByID(ctx context.Context, id uuid.UUID) (*Label, error) {
	query := `SELECT id, project_id, name, color, created_at FROM labels WHERE id = $1`
	var l Label
	err := r.db.QueryRow(ctx, query, id).Scan(&l.ID, &l.ProjectID, &l.Name, &l.Color, &l.CreatedAt)
	return &l, err
}

func (r *PgxLabelRepository) GetByProject(ctx context.Context, projectID uuid.UUID) ([]Label, error) {
	query := `SELECT id, project_id, name, color, created_at FROM labels WHERE project_id = $1 ORDER BY name`
	rows, err := r.db.Query(ctx, query, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var labels []Label
	for rows.Next() {
		var l Label
		if err := rows.Scan(&l.ID, &l.ProjectID, &l.Name, &l.Color, &l.CreatedAt); err != nil {
			return nil, err
		}
		labels = append(labels, l)
	}
	return labels, nil
}

func (r *PgxLabelRepository) Update(ctx context.Context, label *Label) error {
	query := `UPDATE labels SET name=$1, color=$2 WHERE id=$3`
	_, err := r.db.Exec(ctx, query, label.Name, label.Color, label.ID)
	return err
}

func (r *PgxLabelRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM labels WHERE id = $1`
	_, err := r.db.Exec(ctx, query, id)
	return err
}

func (r *PgxLabelRepository) AddToTask(ctx context.Context, taskID, labelID uuid.UUID) error {
	query := `INSERT INTO task_labels (task_id, label_id) VALUES ($1, $2) ON CONFLICT DO NOTHING`
	_, err := r.db.Exec(ctx, query, taskID, labelID)
	return err
}

func (r *PgxLabelRepository) RemoveFromTask(ctx context.Context, taskID, labelID uuid.UUID) error {
	query := `DELETE FROM task_labels WHERE task_id = $1 AND label_id = $2`
	_, err := r.db.Exec(ctx, query, taskID, labelID)
	return err
}

func (r *PgxLabelRepository) GetByTask(ctx context.Context, taskID uuid.UUID) ([]Label, error) {
	query := `
		SELECT l.id, l.project_id, l.name, l.color, l.created_at
		FROM labels l
		JOIN task_labels tl ON l.id = tl.label_id
		WHERE tl.task_id = $1
	`
	rows, err := r.db.Query(ctx, query, taskID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var labels []Label
	for rows.Next() {
		var l Label
		if err := rows.Scan(&l.ID, &l.ProjectID, &l.Name, &l.Color, &l.CreatedAt); err != nil {
			return nil, err
		}
		labels = append(labels, l)
	}
	return labels, nil
}

type TaskService struct {
	repo      Repository
	labelRepo LabelRepository
}

func NewTaskService(repo Repository, labelRepo LabelRepository) *TaskService {
	return &TaskService{repo: repo, labelRepo: labelRepo}
}

func (s *TaskService) Create(ctx context.Context, boardID uuid.UUID, input CreateTaskInput) (*Task, error) {
	status := input.Status
	if status == "" {
		status = StatusBacklog
	}

	task := &Task{
		ID:          uuid.New(),
		BoardID:     boardID,
		EpicID:      input.EpicID,
		Title:       input.Title,
		Description: input.Description,
		Status:      status,
		Position:    0,
		StoryPoints: input.StoryPoints,
		DueDate:     input.DueDate,
		AssigneeID:  input.AssigneeID,
		ReporterID:  input.ReporterID,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := s.repo.Create(ctx, task); err != nil {
		return nil, err
	}

	return task, nil
}

func (s *TaskService) GetByID(ctx context.Context, id uuid.UUID) (*Task, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *TaskService) GetByBoard(ctx context.Context, boardID uuid.UUID, filter TaskFilter) ([]Task, error) {
	if filter.Limit == 0 {
		filter.Limit = 50
	}
	return s.repo.GetByBoard(ctx, boardID, filter)
}

func (s *TaskService) Update(ctx context.Context, id uuid.UUID, input UpdateTaskInput) (*Task, error) {
	t, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if input.Title != "" {
		t.Title = input.Title
	}
	if input.Description != "" {
		t.Description = input.Description
	}
	if input.EpicID != nil {
		t.EpicID = input.EpicID
	}
	if input.Status != "" {
		t.Status = input.Status
	}
	if input.Position != nil {
		t.Position = *input.Position
	}
	if input.StoryPoints != nil {
		t.StoryPoints = input.StoryPoints
	}
	if input.DueDate != nil {
		t.DueDate = input.DueDate
	}
	if input.AssigneeID != nil {
		t.AssigneeID = input.AssigneeID
	}

	if err := s.repo.Update(ctx, t); err != nil {
		return nil, err
	}

	return t, nil
}

func (s *TaskService) UpdateStatus(ctx context.Context, id uuid.UUID, status TaskStatus) error {
	return s.repo.UpdateStatus(ctx, id, status)
}

func (s *TaskService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}

func (s *TaskService) CreateLabel(ctx context.Context, projectID uuid.UUID, input CreateLabelInput) (*Label, error) {
	color := input.Color
	if color == "" {
		color = "#808080"
	}

	label := &Label{
		ID:        uuid.New(),
		ProjectID: projectID,
		Name:      input.Name,
		Color:     color,
		CreatedAt: time.Now(),
	}

	if err := s.labelRepo.Create(ctx, label); err != nil {
		return nil, err
	}

	return label, nil
}

func (s *TaskService) GetLabels(ctx context.Context, projectID uuid.UUID) ([]Label, error) {
	return s.labelRepo.GetByProject(ctx, projectID)
}

func (s *TaskService) DeleteLabel(ctx context.Context, id uuid.UUID) error {
	return s.labelRepo.Delete(ctx, id)
}

func (s *TaskService) AddLabelToTask(ctx context.Context, taskID, labelID uuid.UUID) error {
	return s.labelRepo.AddToTask(ctx, taskID, labelID)
}

func (s *TaskService) RemoveLabelFromTask(ctx context.Context, taskID, labelID uuid.UUID) error {
	return s.labelRepo.RemoveFromTask(ctx, taskID, labelID)
}

func (s *TaskService) GetTaskLabels(ctx context.Context, taskID uuid.UUID) ([]Label, error) {
	return s.labelRepo.GetByTask(ctx, taskID)
}
