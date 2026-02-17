package project

import (
	"context"

	"github.com/google/uuid"
	"github.com/monachy/projek/internal/common"
	"github.com/monachy/projek/internal/task"
)

type Service struct {
	repo      Repository
	epicRepo  EpicRepository
	boardRepo BoardRepository
	taskRepo  task.Repository
	labelRepo task.LabelRepository
}

func NewService(repo Repository, epicRepo EpicRepository, boardRepo BoardRepository, taskRepo task.Repository, labelRepo task.LabelRepository) *Service {
	return &Service{
		repo:      repo,
		epicRepo:  epicRepo,
		boardRepo: boardRepo,
		taskRepo:  taskRepo,
		labelRepo: labelRepo,
	}
}

func (s *Service) Create(ctx context.Context, communityID uuid.UUID, input CreateProjectInput) (*Project, error) {
	project := &Project{
		ID:          uuid.New(),
		CommunityID: communityID,
		Name:        input.Name,
		Description: input.Description,
		Key:         input.Key,
		IsArchived:  false,
		CreatedAt:   common.Now(),
		UpdatedAt:   common.Now(),
	}

	if err := s.repo.Create(ctx, project); err != nil {
		return nil, err
	}

	defaultBoard := &Board{
		ID:        uuid.New(),
		ProjectID: project.ID,
		Name:      "Main Board",
		CreatedAt: common.Now(),
		UpdatedAt: common.Now(),
	}
	s.boardRepo.Create(ctx, defaultBoard)

	return project, nil
}

func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (*Project, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *Service) GetByCommunity(ctx context.Context, communityID uuid.UUID, limit, offset int) ([]Project, error) {
	if limit <= 0 {
		limit = 20
	}
	return s.repo.GetByCommunity(ctx, communityID, limit, offset)
}

func (s *Service) Update(ctx context.Context, id uuid.UUID, input UpdateProjectInput) (*Project, error) {
	p, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if input.Name != "" {
		p.Name = input.Name
	}
	if input.Description != "" {
		p.Description = input.Description
	}
	if input.Key != "" {
		p.Key = input.Key
	}
	if input.IsArchived != nil {
		p.IsArchived = *input.IsArchived
	}

	if err := s.repo.Update(ctx, p); err != nil {
		return nil, err
	}

	return p, nil
}

func (s *Service) Archive(ctx context.Context, id uuid.UUID) error {
	return s.repo.Archive(ctx, id)
}

func (s *Service) Delete(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}

func (s *Service) CreateEpic(ctx context.Context, projectID uuid.UUID, input CreateEpicInput) (*Epic, error) {
	epic := &Epic{
		ID:          uuid.New(),
		ProjectID:   projectID,
		Name:        input.Name,
		Description: input.Description,
		StartDate:   input.StartDate,
		EndDate:     input.EndDate,
		CreatedAt:   common.Now(),
		UpdatedAt:   common.Now(),
	}

	if err := s.epicRepo.Create(ctx, epic); err != nil {
		return nil, err
	}

	return epic, nil
}

func (s *Service) GetEpics(ctx context.Context, projectID uuid.UUID) ([]Epic, error) {
	return s.epicRepo.GetByProject(ctx, projectID)
}

func (s *Service) UpdateEpic(ctx context.Context, id uuid.UUID, input CreateEpicInput) (*Epic, error) {
	e, err := s.epicRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	e.Name = input.Name
	e.Description = input.Description
	e.StartDate = input.StartDate
	e.EndDate = input.EndDate

	if err := s.epicRepo.Update(ctx, e); err != nil {
		return nil, err
	}

	return e, nil
}

func (s *Service) DeleteEpic(ctx context.Context, id uuid.UUID) error {
	return s.epicRepo.Delete(ctx, id)
}

func (s *Service) CreateBoard(ctx context.Context, projectID uuid.UUID, input CreateBoardInput) (*Board, error) {
	board := &Board{
		ID:        uuid.New(),
		ProjectID: projectID,
		Name:      input.Name,
		CreatedAt: common.Now(),
		UpdatedAt: common.Now(),
	}

	if err := s.boardRepo.Create(ctx, board); err != nil {
		return nil, err
	}

	return board, nil
}

func (s *Service) GetBoards(ctx context.Context, projectID uuid.UUID) ([]Board, error) {
	return s.boardRepo.GetByProject(ctx, projectID)
}

func (s *Service) GetBoardByID(ctx context.Context, id uuid.UUID) (*Board, error) {
	return s.boardRepo.GetByID(ctx, id)
}

func (s *Service) DeleteBoard(ctx context.Context, id uuid.UUID) error {
	return s.boardRepo.Delete(ctx, id)
}
