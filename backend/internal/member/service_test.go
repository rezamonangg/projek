package member

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/monachy/projek/internal/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

var testCommunityID = uuid.MustParse("11111111-1111-1111-1111-111111111111")
var testMemberID = uuid.MustParse("22222222-2222-2222-2222-222222222222")

func testMember() *Member {
	now := time.Now()
	return &Member{
		ID:           testMemberID,
		CommunityID:  testCommunityID,
		Email:        "john@example.com",
		PasswordHash: "$2a$10$hashedpassword",
		FirstName:    "John",
		LastName:     "Doe",
		DateOfBirth:  nil,
		Role:         RoleMember,
		IsActive:     true,
		LastLoginAt:  nil,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
}

func testCreateMemberInput() CreateMemberInput {
	return CreateMemberInput{
		Email:     "john@example.com",
		Password:  "password123",
		FirstName: "John",
		LastName:  "Doe",
	}
}

func testUpdateMemberInput() UpdateMemberInput {
	trueVal := true
	return UpdateMemberInput{
		FirstName: "JohnUpdated",
		LastName:  "DoeUpdated",
		Role:      RoleAdmin,
		IsActive:  &trueVal,
	}
}

func TestService_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := NewMockRepository(ctrl)
	svc := NewService(repo)

	ctx := context.Background()
	input := testCreateMemberInput()

	repo.EXPECT().Create(ctx, gomock.Any()).Return(nil)

	member, err := svc.Create(ctx, testCommunityID, input)

	require.NoError(t, err)
	assert.NotNil(t, member)
	assert.Equal(t, input.Email, member.Email)
	assert.Equal(t, input.FirstName, member.FirstName)
	assert.Equal(t, input.LastName, member.LastName)
	assert.True(t, member.IsActive)
}

func TestService_Create_PasswordHashing(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := NewMockRepository(ctrl)
	svc := NewService(repo)

	ctx := context.Background()
	input := testCreateMemberInput()

	repo.EXPECT().Create(ctx, gomock.Any()).DoAndReturn(func(c context.Context, m *Member) error {
		assert.NotEqual(t, input.Password, m.PasswordHash, "password should be hashed")
		assert.NotEmpty(t, m.PasswordHash)
		return nil
	})

	member, err := svc.Create(ctx, testCommunityID, input)

	require.NoError(t, err)
	assert.NotNil(t, member)
}

func TestService_Create_DefaultRole(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := NewMockRepository(ctrl)
	svc := NewService(repo)

	ctx := context.Background()
	input := testCreateMemberInput()
	input.Role = ""

	repo.EXPECT().Create(ctx, gomock.Any()).DoAndReturn(func(c context.Context, m *Member) error {
		assert.Equal(t, RoleMember, m.Role, "default role should be member")
		return nil
	})

	member, err := svc.Create(ctx, testCommunityID, input)

	require.NoError(t, err)
	assert.NotNil(t, member)
	assert.Equal(t, RoleMember, member.Role)
}

func TestService_GetByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := NewMockRepository(ctrl)
	svc := NewService(repo)

	ctx := context.Background()
	expectedMember := testMember()

	repo.EXPECT().GetByID(ctx, expectedMember.ID).Return(expectedMember, nil)

	member, err := svc.GetByID(ctx, expectedMember.ID)

	require.NoError(t, err)
	assert.Equal(t, expectedMember, member)
}

func TestService_GetByID_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := NewMockRepository(ctrl)
	svc := NewService(repo)

	ctx := context.Background()

	repo.EXPECT().GetByID(ctx, testMemberID).Return(nil, common.ErrNotFound)

	member, err := svc.GetByID(ctx, testMemberID)

	assert.Nil(t, member)
	assert.Error(t, err)
	assert.True(t, errors.Is(err, common.ErrNotFound))
}

func TestService_GetByEmail(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := NewMockRepository(ctrl)
	svc := NewService(repo)

	ctx := context.Background()
	expectedMember := testMember()

	repo.EXPECT().GetByEmail(ctx, expectedMember.Email).Return(expectedMember, nil)

	member, err := svc.GetByEmail(ctx, expectedMember.Email)

	require.NoError(t, err)
	assert.Equal(t, expectedMember, member)
}

func TestService_GetByCommunity(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := NewMockRepository(ctrl)
	svc := NewService(repo)

	ctx := context.Background()
	expectedMembers := []Member{*testMember()}

	repo.EXPECT().GetByCommunity(ctx, testCommunityID, 10, 0).Return(expectedMembers, nil)

	members, err := svc.GetByCommunity(ctx, testCommunityID, 10, 0)

	require.NoError(t, err)
	assert.Equal(t, expectedMembers, members)
}

func TestService_Update(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := NewMockRepository(ctrl)
	svc := NewService(repo)

	ctx := context.Background()
	existingMember := testMember()
	updateInput := testUpdateMemberInput()

	repo.EXPECT().GetByID(ctx, existingMember.ID).Return(existingMember, nil)
	repo.EXPECT().Update(ctx, gomock.Any()).Return(nil)

	member, err := svc.Update(ctx, existingMember.ID, updateInput)

	require.NoError(t, err)
	assert.NotNil(t, member)
	assert.Equal(t, updateInput.FirstName, member.FirstName)
	assert.Equal(t, updateInput.LastName, member.LastName)
	assert.Equal(t, updateInput.Role, member.Role)
}

func TestService_Update_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := NewMockRepository(ctrl)
	svc := NewService(repo)

	ctx := context.Background()

	repo.EXPECT().GetByID(ctx, testMemberID).Return(nil, common.ErrNotFound)

	member, err := svc.Update(ctx, testMemberID, testUpdateMemberInput())

	assert.Nil(t, member)
	assert.Error(t, err)
}

func TestService_Delete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := NewMockRepository(ctrl)
	svc := NewService(repo)

	ctx := context.Background()

	repo.EXPECT().Delete(ctx, testMemberID).Return(nil)

	err := svc.Delete(ctx, testMemberID)

	require.NoError(t, err)
}

func TestService_ChangePassword(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := NewMockRepository(ctrl)
	svc := NewService(repo)

	ctx := context.Background()
	existingMember := testMember()
	newPassword := "newpassword123"

	repo.EXPECT().GetByID(ctx, existingMember.ID).Return(existingMember, nil)
	repo.EXPECT().Update(ctx, gomock.Any()).DoAndReturn(func(c context.Context, m *Member) error {
		assert.NotEmpty(t, m.PasswordHash, "password should be hashed")
		return nil
	})

	err := svc.ChangePassword(ctx, existingMember.ID, newPassword)

	require.NoError(t, err)
}

func TestService_ChangePassword_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := NewMockRepository(ctrl)
	svc := NewService(repo)

	ctx := context.Background()

	repo.EXPECT().GetByID(ctx, testMemberID).Return(nil, common.ErrNotFound)

	err := svc.ChangePassword(ctx, testMemberID, "newpassword123")

	assert.Error(t, err)
}
