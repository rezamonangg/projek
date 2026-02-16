package auth

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/monachy/projek/internal/common"
	"github.com/monachy/projek/internal/member"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestAuthService_Login(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockMemberRepo := member.NewMockRepository(ctrl)
	mockSessionStore := NewMockSessionStore(ctrl)
	svc := NewAuthService(mockMemberRepo, mockSessionStore)

	ctx := context.Background()
	passwordHash, err := common.HashPassword("password123")
	require.NoError(t, err)

	testMember := &member.Member{
		ID:           uuid.New(),
		CommunityID:  uuid.New(),
		Email:        "john@example.com",
		PasswordHash: passwordHash,
		FirstName:    "John",
		LastName:     "Doe",
		IsActive:     true,
	}

	mockMemberRepo.EXPECT().GetByEmail(ctx, "john@example.com").Return(testMember, nil)
	mockSessionStore.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)
	mockMemberRepo.EXPECT().UpdateLastLogin(gomock.Any(), testMember.ID).Return(nil)

	session, m, err := svc.Login(ctx, "john@example.com", "password123")

	require.NoError(t, err)
	assert.NotNil(t, session)
	assert.NotNil(t, m)
	assert.Equal(t, testMember.Email, m.Email)
}

func TestAuthService_Login_MemberNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockMemberRepo := member.NewMockRepository(ctrl)
	mockSessionStore := NewMockSessionStore(ctrl)
	svc := NewAuthService(mockMemberRepo, mockSessionStore)

	ctx := context.Background()

	mockMemberRepo.EXPECT().GetByEmail(ctx, "john@example.com").Return(nil, common.ErrNotFound)

	session, m, err := svc.Login(ctx, "john@example.com", "password123")

	assert.Nil(t, session)
	assert.Nil(t, m)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid credentials")
}

func TestAuthService_Login_InactiveMember(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockMemberRepo := member.NewMockRepository(ctrl)
	mockSessionStore := NewMockSessionStore(ctrl)
	svc := NewAuthService(mockMemberRepo, mockSessionStore)

	ctx := context.Background()
	testMember := &member.Member{
		ID:          uuid.New(),
		CommunityID: uuid.New(),
		Email:       "john@example.com",
		IsActive:    false,
	}

	mockMemberRepo.EXPECT().GetByEmail(ctx, "john@example.com").Return(testMember, nil)

	session, m, err := svc.Login(ctx, "john@example.com", "password123")

	assert.Nil(t, session)
	assert.Nil(t, m)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "inactive")
}

func TestAuthService_Login_InvalidPassword(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockMemberRepo := member.NewMockRepository(ctrl)
	mockSessionStore := NewMockSessionStore(ctrl)
	svc := NewAuthService(mockMemberRepo, mockSessionStore)

	ctx := context.Background()
	passwordHash, err := common.HashPassword("correctpassword")
	require.NoError(t, err)

	testMember := &member.Member{
		ID:           uuid.New(),
		CommunityID:  uuid.New(),
		Email:        "john@example.com",
		PasswordHash: passwordHash,
		FirstName:    "John",
		LastName:     "Doe",
		IsActive:     true,
	}

	mockMemberRepo.EXPECT().GetByEmail(ctx, "john@example.com").Return(testMember, nil)

	session, m, err := svc.Login(ctx, "john@example.com", "wrongpassword")

	assert.Nil(t, session)
	assert.Nil(t, m)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid credentials")
}

func TestAuthService_Login_SessionStoreError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockMemberRepo := member.NewMockRepository(ctrl)
	mockSessionStore := NewMockSessionStore(ctrl)
	svc := NewAuthService(mockMemberRepo, mockSessionStore)

	ctx := context.Background()
	passwordHash, err := common.HashPassword("password123")
	require.NoError(t, err)

	testMember := &member.Member{
		ID:           uuid.New(),
		CommunityID:  uuid.New(),
		Email:        "john@example.com",
		PasswordHash: passwordHash,
		FirstName:    "John",
		IsActive:     true,
	}

	mockMemberRepo.EXPECT().GetByEmail(ctx, "john@example.com").Return(testMember, nil)
	mockSessionStore.EXPECT().Create(gomock.Any(), gomock.Any()).Return(errors.New("redis error"))

	session, m, err := svc.Login(ctx, "john@example.com", "password123")

	assert.Nil(t, session)
	assert.Nil(t, m)
	assert.Error(t, err)
}

func TestAuthService_Logout(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockMemberRepo := member.NewMockRepository(ctrl)
	mockSessionStore := NewMockSessionStore(ctrl)
	svc := NewAuthService(mockMemberRepo, mockSessionStore)

	ctx := context.Background()
	sessionID := "session-123"

	mockSessionStore.EXPECT().Delete(ctx, sessionID).Return(nil)

	err := svc.Logout(ctx, sessionID)

	require.NoError(t, err)
}

func TestAuthService_ValidateSession(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockMemberRepo := member.NewMockRepository(ctrl)
	mockSessionStore := NewMockSessionStore(ctrl)
	svc := NewAuthService(mockMemberRepo, mockSessionStore)

	ctx := context.Background()
	sessionID := "session-123"
	expectedSession := &Session{
		ID:          sessionID,
		MemberID:    uuid.New(),
		CommunityID: uuid.New(),
		ExpiresAt:   time.Now().Add(time.Hour),
	}

	mockSessionStore.EXPECT().Get(ctx, sessionID).Return(expectedSession, nil)

	session, err := svc.ValidateSession(ctx, sessionID)

	require.NoError(t, err)
	assert.Equal(t, expectedSession, session)
}

func TestAuthService_ValidateSession_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockMemberRepo := member.NewMockRepository(ctrl)
	mockSessionStore := NewMockSessionStore(ctrl)
	svc := NewAuthService(mockMemberRepo, mockSessionStore)

	ctx := context.Background()
	sessionID := "session-123"

	mockSessionStore.EXPECT().Get(ctx, sessionID).Return(nil, ErrInvalidSession)

	session, err := svc.ValidateSession(ctx, sessionID)

	assert.Nil(t, session)
	assert.Error(t, err)
	assert.True(t, errors.Is(err, ErrInvalidSession))
}

func TestAuthService_GetMemberBySession(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockMemberRepo := member.NewMockRepository(ctrl)
	mockSessionStore := NewMockSessionStore(ctrl)
	svc := NewAuthService(mockMemberRepo, mockSessionStore)

	ctx := context.Background()
	sessionID := "session-123"
	memberID := uuid.New()
	testMember := &member.Member{
		ID:          memberID,
		CommunityID: uuid.New(),
		Email:       "john@example.com",
	}
	session := &Session{
		ID:       sessionID,
		MemberID: memberID,
	}

	mockSessionStore.EXPECT().Get(ctx, sessionID).Return(session, nil)
	mockMemberRepo.EXPECT().GetByID(ctx, memberID).Return(testMember, nil)

	m, err := svc.GetMemberBySession(ctx, sessionID)

	require.NoError(t, err)
	assert.Equal(t, testMember, m)
}
