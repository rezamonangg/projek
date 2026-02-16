package project

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
var testProjectID = uuid.MustParse("22222222-2222-2222-2222-222222222222")

func testProject() *Project {
	now := time.Now()
	return &Project{
		ID:          testProjectID,
		CommunityID: testCommunityID,
		Name:        "Test Project",
		Description: "A test project",
		Key:         "TEST",
		IsArchived:  false,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

func testCreateProjectInput() CreateProjectInput {
	return CreateProjectInput{
		Name:        "Test Project",
		Description: "A test project",
		Key:         "TEST",
	}
}

func TestRepository_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := NewMockRepository(ctrl)
	ctx := context.Background()
	p := testProject()

	repo.EXPECT().Create(ctx, p).Return(nil)

	err := repo.Create(ctx, p)

	require.NoError(t, err)
}

func TestRepository_GetByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := NewMockRepository(ctrl)
	ctx := context.Background()
	expected := testProject()

	repo.EXPECT().GetByID(ctx, testProjectID).Return(expected, nil)

	p, err := repo.GetByID(ctx, testProjectID)

	require.NoError(t, err)
	assert.Equal(t, expected, p)
}

func TestRepository_GetByID_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := NewMockRepository(ctrl)
	ctx := context.Background()

	repo.EXPECT().GetByID(ctx, testProjectID).Return(nil, common.ErrNotFound)

	p, err := repo.GetByID(ctx, testProjectID)

	assert.Nil(t, p)
	assert.Error(t, err)
	assert.True(t, errors.Is(err, common.ErrNotFound))
}

func TestRepository_GetByCommunity(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := NewMockRepository(ctrl)
	ctx := context.Background()
	expected := []Project{*testProject()}

	repo.EXPECT().GetByCommunity(ctx, testCommunityID, 10, 0).Return(expected, nil)

	projects, err := repo.GetByCommunity(ctx, testCommunityID, 10, 0)

	require.NoError(t, err)
	assert.Equal(t, expected, projects)
}

func TestRepository_Update(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := NewMockRepository(ctrl)
	ctx := context.Background()
	p := testProject()
	p.Name = "Updated Project"

	repo.EXPECT().Update(ctx, p).Return(nil)

	err := repo.Update(ctx, p)

	require.NoError(t, err)
}

func TestRepository_Delete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := NewMockRepository(ctrl)
	ctx := context.Background()

	repo.EXPECT().Delete(ctx, testProjectID).Return(nil)

	err := repo.Delete(ctx, testProjectID)

	require.NoError(t, err)
}

func TestRepository_Archive(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := NewMockRepository(ctrl)
	ctx := context.Background()

	repo.EXPECT().Archive(ctx, testProjectID).Return(nil)

	err := repo.Archive(ctx, testProjectID)

	require.NoError(t, err)
}
