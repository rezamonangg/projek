package community

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

func testCommunity() *Community {
	now := time.Now()
	return &Community{
		ID:        testCommunityID,
		Name:      "Test Community",
		Slug:      "test-community",
		CreatedAt: now,
		UpdatedAt: now,
	}
}

func testCreateCommunityInput() CreateCommunityInput {
	return CreateCommunityInput{
		Name: "Test Community",
		Slug: "test-community",
	}
}

func TestRepository_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := NewMockRepository(ctrl)
	ctx := context.Background()
	c := testCommunity()

	repo.EXPECT().Create(ctx, c).Return(nil)

	err := repo.Create(ctx, c)

	require.NoError(t, err)
}

func TestRepository_GetByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := NewMockRepository(ctrl)
	ctx := context.Background()
	expected := testCommunity()

	repo.EXPECT().GetByID(ctx, testCommunityID).Return(expected, nil)

	c, err := repo.GetByID(ctx, testCommunityID)

	require.NoError(t, err)
	assert.Equal(t, expected, c)
}

func TestRepository_GetByID_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := NewMockRepository(ctrl)
	ctx := context.Background()

	repo.EXPECT().GetByID(ctx, testCommunityID).Return(nil, common.ErrNotFound)

	c, err := repo.GetByID(ctx, testCommunityID)

	assert.Nil(t, c)
	assert.Error(t, err)
	assert.True(t, errors.Is(err, common.ErrNotFound))
}

func TestRepository_GetBySlug(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := NewMockRepository(ctrl)
	ctx := context.Background()
	expected := testCommunity()

	repo.EXPECT().GetBySlug(ctx, "test-community").Return(expected, nil)

	c, err := repo.GetBySlug(ctx, "test-community")

	require.NoError(t, err)
	assert.Equal(t, expected, c)
}

func TestRepository_Update(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := NewMockRepository(ctrl)
	ctx := context.Background()
	c := testCommunity()
	c.Name = "Updated Name"

	repo.EXPECT().Update(ctx, c).Return(nil)

	err := repo.Update(ctx, c)

	require.NoError(t, err)
}

func TestRepository_Delete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := NewMockRepository(ctrl)
	ctx := context.Background()

	repo.EXPECT().Delete(ctx, testCommunityID).Return(nil)

	err := repo.Delete(ctx, testCommunityID)

	require.NoError(t, err)
}
