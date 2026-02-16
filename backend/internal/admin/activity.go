package admin

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ActivityLog struct {
	ID          uuid.UUID              `json:"id"`
	CommunityID uuid.UUID              `json:"community_id"`
	ProjectID   *uuid.UUID             `json:"project_id,omitempty"`
	MemberID    *uuid.UUID             `json:"member_id,omitempty"`
	Action      string                 `json:"action"`
	EntityType  string                 `json:"entity_type"`
	EntityID    *uuid.UUID             `json:"entity_id,omitempty"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
	CreatedAt   time.Time              `json:"created_at"`
}

type ActivityLogRepository interface {
	Create(ctx context.Context, log *ActivityLog) error
	GetByCommunity(ctx context.Context, communityID uuid.UUID, limit, offset int) ([]ActivityLog, error)
	GetByProject(ctx context.Context, projectID uuid.UUID, limit, offset int) ([]ActivityLog, error)
}

type PgxActivityLogRepository struct {
	db *pgxpool.Pool
}

func NewActivityLogRepository(db *pgxpool.Pool) ActivityLogRepository {
	return &PgxActivityLogRepository{db: db}
}

func (r *PgxActivityLogRepository) Create(ctx context.Context, log *ActivityLog) error {
	metadata, _ := json.Marshal(log.Metadata)
	query := `
		INSERT INTO activity_logs (id, community_id, project_id, member_id, action, entity_type, entity_id, metadata, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`
	_, err := r.db.Exec(ctx, query,
		log.ID, log.CommunityID, log.ProjectID, log.MemberID,
		log.Action, log.EntityType, log.EntityID, metadata, log.CreatedAt,
	)
	return err
}

func (r *PgxActivityLogRepository) GetByCommunity(ctx context.Context, communityID uuid.UUID, limit, offset int) ([]ActivityLog, error) {
	query := `
		SELECT id, community_id, project_id, member_id, action, entity_type, entity_id, metadata, created_at
		FROM activity_logs WHERE community_id = $1
		ORDER BY created_at DESC LIMIT $2 OFFSET $3
	`
	rows, err := r.db.Query(ctx, query, communityID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []ActivityLog
	for rows.Next() {
		var l ActivityLog
		var metadata []byte
		if err := rows.Scan(
			&l.ID, &l.CommunityID, &l.ProjectID, &l.MemberID,
			&l.Action, &l.EntityType, &l.EntityID, &metadata, &l.CreatedAt,
		); err != nil {
			return nil, err
		}
		json.Unmarshal(metadata, &l.Metadata)
		logs = append(logs, l)
	}
	return logs, nil
}

func (r *PgxActivityLogRepository) GetByProject(ctx context.Context, projectID uuid.UUID, limit, offset int) ([]ActivityLog, error) {
	query := `
		SELECT id, community_id, project_id, member_id, action, entity_type, entity_id, metadata, created_at
		FROM activity_logs WHERE project_id = $1
		ORDER BY created_at DESC LIMIT $2 OFFSET $3
	`
	rows, err := r.db.Query(ctx, query, projectID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []ActivityLog
	for rows.Next() {
		var l ActivityLog
		var metadata []byte
		if err := rows.Scan(
			&l.ID, &l.CommunityID, &l.ProjectID, &l.MemberID,
			&l.Action, &l.EntityType, &l.EntityID, &metadata, &l.CreatedAt,
		); err != nil {
			return nil, err
		}
		json.Unmarshal(metadata, &l.Metadata)
		logs = append(logs, l)
	}
	return logs, nil
}

type ActivityService struct {
	repo ActivityLogRepository
}

func NewActivityService(repo ActivityLogRepository) *ActivityService {
	return &ActivityService{repo: repo}
}

func (s *ActivityService) Log(ctx context.Context, communityID, projectID, memberID uuid.UUID, action, entityType string, entityID *uuid.UUID, metadata map[string]interface{}) error {
	log := &ActivityLog{
		ID:          uuid.New(),
		CommunityID: communityID,
		ProjectID:   &projectID,
		MemberID:    &memberID,
		Action:      action,
		EntityType:  entityType,
		EntityID:    entityID,
		Metadata:    metadata,
		CreatedAt:   time.Now(),
	}
	return s.repo.Create(ctx, log)
}

func (s *ActivityService) GetRecentActivity(ctx context.Context, communityID uuid.UUID, limit, offset int) ([]ActivityLog, error) {
	if limit <= 0 {
		limit = 20
	}
	return s.repo.GetByCommunity(ctx, communityID, limit, offset)
}

func (s *ActivityService) GetProjectActivity(ctx context.Context, projectID uuid.UUID, limit, offset int) ([]ActivityLog, error) {
	if limit <= 0 {
		limit = 20
	}
	return s.repo.GetByProject(ctx, projectID, limit, offset)
}
