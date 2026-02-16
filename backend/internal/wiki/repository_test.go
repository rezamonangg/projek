package wiki

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

var testProjectID = uuid.MustParse("11111111-1111-1111-1111-111111111111")
var testPageID = uuid.MustParse("22222222-2222-2222-2222-222222222222")

func testWikiPage() *WikiPage {
	now := time.Now()
	return &WikiPage{
		ID:        testPageID,
		ProjectID: testProjectID,
		Title:     "Test Page",
		Slug:      "test-page",
		Content:   "# Test Content",
		ParentID:  nil,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

func TestRepository_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := NewMockRepository(ctrl)
	ctx := context.Background()
	page := testWikiPage()

	repo.EXPECT().Create(ctx, page).Return(nil)

	err := repo.Create(ctx, page)

	require.NoError(t, err)
}

func TestRepository_GetByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := NewMockRepository(ctrl)
	ctx := context.Background()
	expected := testWikiPage()

	repo.EXPECT().GetByID(ctx, testPageID).Return(expected, nil)

	page, err := repo.GetByID(ctx, testPageID)

	require.NoError(t, err)
	assert.Equal(t, expected, page)
}

func TestRepository_GetByID_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := NewMockRepository(ctrl)
	ctx := context.Background()

	repo.EXPECT().GetByID(ctx, testPageID).Return(nil, common.ErrNotFound)

	page, err := repo.GetByID(ctx, testPageID)

	assert.Nil(t, page)
	assert.Error(t, err)
	assert.True(t, errors.Is(err, common.ErrNotFound))
}

func TestRepository_GetBySlug(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := NewMockRepository(ctrl)
	ctx := context.Background()
	expected := testWikiPage()

	repo.EXPECT().GetBySlug(ctx, testProjectID, "test-page").Return(expected, nil)

	page, err := repo.GetBySlug(ctx, testProjectID, "test-page")

	require.NoError(t, err)
	assert.Equal(t, expected, page)
}

func TestRepository_GetByProject(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := NewMockRepository(ctrl)
	ctx := context.Background()
	expected := []WikiPage{*testWikiPage()}

	repo.EXPECT().GetByProject(ctx, testProjectID).Return(expected, nil)

	pages, err := repo.GetByProject(ctx, testProjectID)

	require.NoError(t, err)
	assert.Equal(t, expected, pages)
}

func TestRepository_Update(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := NewMockRepository(ctrl)
	ctx := context.Background()
	page := testWikiPage()
	page.Title = "Updated Page"

	repo.EXPECT().Update(ctx, page).Return(nil)

	err := repo.Update(ctx, page)

	require.NoError(t, err)
}

func TestRepository_Delete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := NewMockRepository(ctrl)
	ctx := context.Background()

	repo.EXPECT().Delete(ctx, testPageID).Return(nil)

	err := repo.Delete(ctx, testPageID)

	require.NoError(t, err)
}
