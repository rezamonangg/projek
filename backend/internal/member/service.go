package member

import (
	"context"

	"github.com/google/uuid"
	"github.com/monachy/projek/internal/common"
)

var ErrMemberNotFound = common.ErrNotFound

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Create(ctx context.Context, communityID uuid.UUID, input CreateMemberInput) (*Member, error) {
	hash, err := common.HashPassword(input.Password)
	if err != nil {
		return nil, err
	}

	role := input.Role
	if role == "" {
		role = RoleMember
	}

	member := &Member{
		ID:           uuid.New(),
		CommunityID:  communityID,
		Email:        input.Email,
		PasswordHash: hash,
		FirstName:    input.FirstName,
		LastName:     input.LastName,
		DateOfBirth:  input.DateOfBirth,
		Role:         role,
		IsActive:     true,
		CreatedAt:    common.Now(),
		UpdatedAt:    common.Now(),
	}

	if err := s.repo.Create(ctx, member); err != nil {
		return nil, err
	}

	return member, nil
}

func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (*Member, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *Service) GetByEmail(ctx context.Context, email string) (*Member, error) {
	return s.repo.GetByEmail(ctx, email)
}

func (s *Service) GetByCommunity(ctx context.Context, communityID uuid.UUID, limit, offset int) ([]Member, error) {
	return s.repo.GetByCommunity(ctx, communityID, limit, offset)
}

func (s *Service) Update(ctx context.Context, id uuid.UUID, input UpdateMemberInput) (*Member, error) {
	m, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if input.FirstName != "" {
		m.FirstName = input.FirstName
	}
	if input.LastName != "" {
		m.LastName = input.LastName
	}
	if input.DateOfBirth != nil {
		m.DateOfBirth = input.DateOfBirth
	}
	if input.Role != "" {
		m.Role = input.Role
	}
	if input.IsActive != nil {
		m.IsActive = *input.IsActive
	}

	if err := s.repo.Update(ctx, m); err != nil {
		return nil, err
	}

	return m, nil
}

func (s *Service) Delete(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}

func (s *Service) ChangePassword(ctx context.Context, id uuid.UUID, newPassword string) error {
	m, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	hash, err := common.HashPassword(newPassword)
	if err != nil {
		return err
	}

	m.PasswordHash = hash
	return s.repo.Update(ctx, m)
}
