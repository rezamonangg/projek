package auth

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/monachy/projek/internal/common"
	"github.com/monachy/projek/internal/member"
	"github.com/redis/go-redis/v9"
)

var (
	ErrInvalidSession = errors.New("invalid session")
	ErrSessionExpired = errors.New("session expired")
)

type Session struct {
	ID          string
	MemberID    uuid.UUID
	CommunityID uuid.UUID
	ExpiresAt   time.Time
}

type SessionStore interface {
	Create(ctx context.Context, session *Session) error
	Get(ctx context.Context, sessionID string) (*Session, error)
	Delete(ctx context.Context, sessionID string) error
	DeleteByMemberID(ctx context.Context, memberID uuid.UUID) error
}

type RedisSessionStore struct {
	client *redis.Client
	ttl    time.Duration
}

func NewSessionStore(client *redis.Client, ttl time.Duration) SessionStore {
	return &RedisSessionStore{client: client, ttl: ttl}
}

func (s *RedisSessionStore) Create(ctx context.Context, session *Session) error {
	key := "session:" + session.ID
	return s.client.Set(ctx, key, session.MemberID.String(), s.ttl).Err()
}

func (s *RedisSessionStore) Get(ctx context.Context, sessionID string) (*Session, error) {
	key := "session:" + sessionID
	memberIDStr, err := s.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, ErrInvalidSession
		}
		return nil, err
	}

	memberID, err := uuid.Parse(memberIDStr)
	if err != nil {
		return nil, ErrInvalidSession
	}

	ttl, err := s.client.TTL(ctx, "session:"+sessionID).Result()
	if err != nil || ttl <= 0 {
		return nil, ErrSessionExpired
	}

	return &Session{
		ID:        sessionID,
		MemberID:  memberID,
		ExpiresAt: time.Now().Add(ttl),
	}, nil
}

func (s *RedisSessionStore) Delete(ctx context.Context, sessionID string) error {
	return s.client.Del(ctx, "session:"+sessionID).Err()
}

func (s *RedisSessionStore) DeleteByMemberID(ctx context.Context, memberID uuid.UUID) error {
	pattern := "session:*"
	iter := s.client.Scan(ctx, 0, pattern, 1000).Iterator()

	for iter.Next(ctx) {
		key := iter.Val()
		memberIDStr, err := s.client.Get(ctx, key).Result()
		if err != nil {
			continue
		}
		if memberIDStr == memberID.String() {
			s.client.Del(ctx, key)
		}
	}
	return iter.Err()
}

type AuthService struct {
	memberRepo   member.Repository
	sessionStore SessionStore
}

func NewAuthService(memberRepo member.Repository, sessionStore SessionStore) *AuthService {
	return &AuthService{
		memberRepo:   memberRepo,
		sessionStore: sessionStore,
	}
}

func (s *AuthService) Login(ctx context.Context, email, password string) (*Session, *member.Member, error) {
	m, err := s.memberRepo.GetByEmail(ctx, email)
	if err != nil {
		return nil, nil, errors.New("invalid credentials")
	}

	if !m.IsActive {
		return nil, nil, errors.New("account is inactive")
	}

	if !common.CheckPassword(password, m.PasswordHash) {
		return nil, nil, errors.New("invalid credentials")
	}

	session := &Session{
		ID:          uuid.New().String(),
		MemberID:    m.ID,
		CommunityID: m.CommunityID,
		ExpiresAt:   time.Now().Add(7 * 24 * time.Hour),
	}

	if err := s.sessionStore.Create(ctx, session); err != nil {
		return nil, nil, err
	}

	if err := s.memberRepo.UpdateLastLogin(ctx, m.ID); err != nil {
		return nil, nil, err
	}

	return session, m, nil
}

func (s *AuthService) Logout(ctx context.Context, sessionID string) error {
	return s.sessionStore.Delete(ctx, sessionID)
}

func (s *AuthService) ValidateSession(ctx context.Context, sessionID string) (*Session, error) {
	return s.sessionStore.Get(ctx, sessionID)
}

func (s *AuthService) GetMemberBySession(ctx context.Context, sessionID string) (*member.Member, error) {
	session, err := s.sessionStore.Get(ctx, sessionID)
	if err != nil {
		return nil, err
	}
	return s.memberRepo.GetByID(ctx, session.MemberID)
}
