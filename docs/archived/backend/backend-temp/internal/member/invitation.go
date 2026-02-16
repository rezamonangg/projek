package member

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/monachy/projek/internal/common"
)

type InvitationRepository interface {
	Create(ctx context.Context, invitation *Invitation) error
	GetByToken(ctx context.Context, token string) (*Invitation, error)
	GetByEmail(ctx context.Context, email string) ([]Invitation, error)
	Update(ctx context.Context, invitation *Invitation) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type PgxInvitationRepository struct {
	db *pgxpool.Pool
}

func NewInvitationRepository(db *pgxpool.Pool) InvitationRepository {
	return &PgxInvitationRepository{db: db}
}

func (r *PgxInvitationRepository) Create(ctx context.Context, invitation *Invitation) error {
	query := `
		INSERT INTO invitations (id, community_id, email, token, role, invited_by, expires_at, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`
	_, err := r.db.Exec(ctx, query,
		invitation.ID, invitation.CommunityID, invitation.Email, invitation.Token,
		invitation.Role, invitation.InvitedBy, invitation.ExpiresAt, invitation.CreatedAt,
	)
	return err
}

func (r *PgxInvitationRepository) GetByToken(ctx context.Context, token string) (*Invitation, error) {
	query := `
		SELECT id, community_id, email, token, role, invited_by, accepted_at, expires_at, created_at
		FROM invitations WHERE token = $1
	`
	var i Invitation
	err := r.db.QueryRow(ctx, query, token).Scan(
		&i.ID, &i.CommunityID, &i.Email, &i.Token, &i.Role,
		&i.InvitedBy, &i.AcceptedAt, &i.ExpiresAt, &i.CreatedAt,
	)
	return &i, err
}

func (r *PgxInvitationRepository) GetByEmail(ctx context.Context, email string) ([]Invitation, error) {
	query := `
		SELECT id, community_id, email, token, role, invited_by, accepted_at, expires_at, created_at
		FROM invitations WHERE email = $1 AND accepted_at IS NULL
	`
	rows, err := r.db.Query(ctx, query, email)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var invitations []Invitation
	for rows.Next() {
		var i Invitation
		if err := rows.Scan(
			&i.ID, &i.CommunityID, &i.Email, &i.Token, &i.Role,
			&i.InvitedBy, &i.AcceptedAt, &i.ExpiresAt, &i.CreatedAt,
		); err != nil {
			return nil, err
		}
		invitations = append(invitations, i)
	}
	return invitations, nil
}

func (r *PgxInvitationRepository) Update(ctx context.Context, invitation *Invitation) error {
	query := `
		UPDATE invitations SET email=$1, role=$2, accepted_at=$3, expires_at=$4
		WHERE id=$5
	`
	_, err := r.db.Exec(ctx, query, invitation.Email, invitation.Role, invitation.AcceptedAt, invitation.ExpiresAt, invitation.ID)
	return err
}

func (r *PgxInvitationRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM invitations WHERE id = $1`
	_, err := r.db.Exec(ctx, query, id)
	return err
}

func GenerateToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

type InvitationService struct {
	repo       InvitationRepository
	memberRepo Repository
}

func NewInvitationService(repo InvitationRepository, memberRepo Repository) *InvitationService {
	return &InvitationService{repo: repo, memberRepo: memberRepo}
}

func (s *InvitationService) Create(ctx context.Context, communityID, inviterID uuid.UUID, email string, role Role) (*Invitation, error) {
	token, err := GenerateToken()
	if err != nil {
		return nil, err
	}

	invitation := &Invitation{
		ID:          uuid.New(),
		CommunityID: communityID,
		Email:       email,
		Token:       token,
		Role:        role,
		InvitedBy:   inviterID,
		ExpiresAt:   time.Now().Add(7 * 24 * time.Hour),
		CreatedAt:   time.Now(),
	}

	if err := s.repo.Create(ctx, invitation); err != nil {
		return nil, err
	}

	return invitation, nil
}

func (s *InvitationService) Accept(ctx context.Context, token, password string) (*Member, error) {
	invitation, err := s.repo.GetByToken(ctx, token)
	if err != nil {
		return nil, err
	}

	if time.Now().After(invitation.ExpiresAt) {
		return nil, common.ErrInvalidToken
	}

	if invitation.AcceptedAt != nil {
		return nil, common.ErrConflict
	}

	member, err := s.memberRepo.GetByEmail(ctx, invitation.Email)
	if err == nil {
		return nil, common.ErrUserExists
	}

	member = &Member{
		ID:           uuid.New(),
		CommunityID:  invitation.CommunityID,
		Email:        invitation.Email,
		PasswordHash: "",
		FirstName:    "",
		LastName:     "",
		Role:         invitation.Role,
		IsActive:     true,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := s.memberRepo.Create(ctx, member); err != nil {
		return nil, err
	}

	now := time.Now()
	invitation.AcceptedAt = &now
	s.repo.Update(ctx, invitation)

	return member, nil
}
