package task

import (
	"context"

	"github.com/google/uuid"
	"github.com/monachy/projek/internal/common"
)

type Service struct {
	repo      Repository
	labelRepo LabelRepository
}

func NewService(repo Repository, labelRepo LabelRepository) *Service {
	return &Service{
		repo:      repo,
		labelRepo: labelRepo,
	}
}

func (s *Service) Create(ctx context.Context, boardID uuid.UUID, input CreateTaskInput) (*Task, error) {
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
		CreatedAt:   common.Now(),
		UpdatedAt:   common.Now(),
	}

	if err := s.repo.Create(ctx, task); err != nil {
		return nil, err
	}

	return task, nil
}

func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (*Task, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *Service) GetByBoard(ctx context.Context, boardID uuid.UUID, filter TaskFilter) ([]Task, error) {
	return s.repo.GetByBoard(ctx, boardID, filter)
}

func (s *Service) GetByEpic(ctx context.Context, epicID uuid.UUID) ([]Task, error) {
	return s.repo.GetByEpic(ctx, epicID)
}

func (s *Service) Update(ctx context.Context, id uuid.UUID, input UpdateTaskInput) (*Task, error) {
	task, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if input.Title != "" {
		task.Title = input.Title
	}
	if input.Description != "" {
		task.Description = input.Description
	}
	if input.Status != "" {
		task.Status = input.Status
	}
	if input.Position != nil {
		task.Position = *input.Position
	}
	if input.StoryPoints != nil {
		task.StoryPoints = input.StoryPoints
	}
	if input.DueDate != nil {
		task.DueDate = input.DueDate
	}
	if input.AssigneeID != nil {
		task.AssigneeID = input.AssigneeID
	}
	if input.EpicID != nil {
		task.EpicID = input.EpicID
	}

	if err := s.repo.Update(ctx, task); err != nil {
		return nil, err
	}

	return task, nil
}

func (s *Service) Move(ctx context.Context, id uuid.UUID, status TaskStatus, position int) (*Task, error) {
	task, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	task.Status = status
	task.Position = position

	if err := s.repo.Update(ctx, task); err != nil {
		return nil, err
	}

	return task, nil
}

func (s *Service) Delete(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}

func (s *Service) Search(ctx context.Context, query string, limit int) ([]Task, error) {
	if limit <= 0 {
		limit = 10
	}
	return s.repo.Search(ctx, query, limit)
}
