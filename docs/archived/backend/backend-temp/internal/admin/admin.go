package admin

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/monachy/projek/internal/common"
)

type SettingsRepository interface {
	GetByCommunity(ctx context.Context, communityID uuid.UUID) (*CommunitySettings, error)
	Update(ctx context.Context, settings *CommunitySettings) error
}

type CommunitySettings struct {
	ID                       uuid.UUID `json:"id"`
	CommunityID              uuid.UUID `json:"community_id"`
	AllowMemberRegistration  bool      `json:"allow_member_registration"`
	RequireEmailVerification bool      `json:"require_email_verification"`
	CreatedAt                time.Time `json:"created_at"`
	UpdatedAt                time.Time `json:"updated_at"`
}

type PgxSettingsRepository struct {
	db *pgxpool.Pool
}

func NewSettingsRepository(db *pgxpool.Pool) SettingsRepository {
	return &PgxSettingsRepository{db: db}
}

func (r *PgxSettingsRepository) GetByCommunity(ctx context.Context, communityID uuid.UUID) (*CommunitySettings, error) {
	query := `SELECT id, community_id, allow_member_registration, require_email_verification, created_at, updated_at FROM community_settings WHERE community_id = $1`
	var s CommunitySettings
	err := r.db.QueryRow(ctx, query, communityID).Scan(
		&s.ID, &s.CommunityID, &s.AllowMemberRegistration, &s.RequireEmailVerification, &s.CreatedAt, &s.UpdatedAt,
	)
	return &s, err
}

func (r *PgxSettingsRepository) Update(ctx context.Context, settings *CommunitySettings) error {
	settings.UpdatedAt = time.Now()
	query := `UPDATE community_settings SET allow_member_registration=$1, require_email_verification=$2, updated_at=$3 WHERE community_id=$4`
	_, err := r.db.Exec(ctx, query, settings.AllowMemberRegistration, settings.RequireEmailVerification, settings.UpdatedAt, settings.CommunityID)
	return err
}

type Service struct {
	settingsRepo SettingsRepository
}

func NewService(settingsRepo SettingsRepository) *Service {
	return &Service{settingsRepo: settingsRepo}
}

func (s *Service) GetSettings(ctx context.Context, communityID uuid.UUID) (*CommunitySettings, error) {
	settings, err := s.settingsRepo.GetByCommunity(ctx, communityID)
	if err != nil {
		return nil, common.ErrNotFound
	}
	return settings, nil
}

func (s *Service) UpdateSettings(ctx context.Context, communityID uuid.UUID, allowRegistration, requireVerification bool) (*CommunitySettings, error) {
	settings, err := s.settingsRepo.GetByCommunity(ctx, communityID)
	if err != nil {
		settings = &CommunitySettings{
			ID:                       uuid.New(),
			CommunityID:              communityID,
			AllowMemberRegistration:  allowRegistration,
			RequireEmailVerification: requireVerification,
			CreatedAt:                common.Now(),
			UpdatedAt:                common.Now(),
		}
	}

	settings.AllowMemberRegistration = allowRegistration
	settings.RequireEmailVerification = requireVerification

	if err := s.settingsRepo.Update(ctx, settings); err != nil {
		return nil, err
	}

	return settings, nil
}

type DashboardStats struct {
	TotalMembers   int `json:"total_members"`
	TotalProjects  int `json:"total_projects"`
	TotalTasks     int `json:"total_tasks"`
	ActiveTasks    int `json:"active_tasks"`
	CompletedTasks int `json:"completed_tasks"`
}

func (s *Service) GetDashboardStats(ctx context.Context, communityID uuid.UUID) (*DashboardStats, error) {
	var stats DashboardStats

	err := s.settingsRepo.(*PgxSettingsRepository).db.QueryRow(ctx,
		"SELECT COUNT(*) FROM members WHERE community_id = $1", communityID).Scan(&stats.TotalMembers)
	if err != nil {
		return nil, err
	}

	err = s.settingsRepo.(*PgxSettingsRepository).db.QueryRow(ctx,
		"SELECT COUNT(*) FROM projects WHERE community_id = $1", communityID).Scan(&stats.TotalProjects)
	if err != nil {
		return nil, err
	}

	err = s.settingsRepo.(*PgxSettingsRepository).db.QueryRow(ctx,
		`SELECT COUNT(*), 
			SUM(CASE WHEN status = 'done' THEN 1 ELSE 0 END)
			FROM tasks t
			JOIN boards b ON t.board_id = b.id
			JOIN projects p ON b.project_id = p.id
			WHERE p.community_id = $1`, communityID).Scan(&stats.TotalTasks, &stats.CompletedTasks)
	if err != nil {
		return nil, err
	}

	stats.ActiveTasks = stats.TotalTasks - stats.CompletedTasks

	return &stats, nil
}
