package scripts

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/monachy/projek/internal/common"
	"github.com/monachy/projek/internal/community"
	"github.com/monachy/projek/internal/member"
	"github.com/monachy/projek/test/fixtures"
	"github.com/rs/zerolog"
)

func Seed(ctx context.Context, db *pgxpool.Pool, logger zerolog.Logger) error {
	communityRepo := community.NewRepository(db)
	memberRepo := member.NewRepository(db)

	logger.Info().Msg("seeding community...")
	comm, err := seedCommunity(ctx, communityRepo, logger)
	if err != nil {
		return err
	}

	logger.Info().Msg("seeding admin member...")
	if err := seedAdmin(ctx, memberRepo, comm.ID, logger); err != nil {
		return err
	}

	logger.Info().Msg("seeding complete")
	return nil
}

func seedCommunity(ctx context.Context, repo community.Repository, logger zerolog.Logger) (*community.Community, error) {
	commID := fixtures.CommunityID

	existing, err := repo.GetByID(ctx, commID)
	if err == nil {
		logger.Info().Str("id", commID.String()).Msg("community already exists, skipping")
		return existing, nil
	}

	if !errors.Is(err, pgx.ErrNoRows) {
		return nil, err
	}

	now := time.Now()
	comm := &community.Community{
		ID:        commID,
		Name:      "Test Community",
		Slug:      "test-community",
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := repo.Create(ctx, comm); err != nil {
		return nil, err
	}

	logger.Info().Str("id", comm.ID.String()).Str("slug", comm.Slug).Msg("created community")
	return comm, nil
}

func seedAdmin(ctx context.Context, repo member.Repository, communityID uuid.UUID, logger zerolog.Logger) error {
	adminID := uuid.MustParse("44444444-4444-4444-4444-444444444444")

	existing, err := repo.GetByID(ctx, adminID)
	if err == nil {
		logger.Info().Str("id", existing.ID.String()).Msg("admin already exists, skipping")
		return nil
	}

	if !errors.Is(err, pgx.ErrNoRows) {
		return err
	}

	passwordHash, err := common.HashPassword("admin123")
	if err != nil {
		return err
	}

	now := time.Now()
	admin := &member.Member{
		ID:           adminID,
		CommunityID:  communityID,
		Email:        "admin@example.com",
		PasswordHash: passwordHash,
		FirstName:    "Admin",
		LastName:     "User",
		DateOfBirth:  nil,
		Role:         member.RoleAdmin,
		IsActive:     true,
		LastLoginAt:  nil,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	if err := repo.Create(ctx, admin); err != nil {
		return err
	}

	logger.Info().Str("id", admin.ID.String()).Str("email", admin.Email).Msg("created admin member")
	return nil
}
