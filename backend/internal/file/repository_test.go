package file

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
var testTaskID = uuid.MustParse("22222222-2222-2222-2222-222222222222")
var testFileID = uuid.MustParse("33333333-3333-3333-3333-333333333333")

func testFileAttachment() *FileAttachment {
	return &FileAttachment{
		ID:               testFileID,
		ProjectID:        testProjectID,
		TaskID:           &testTaskID,
		UploaderID:       uuid.MustParse("44444444-4444-4444-4444-444444444444"),
		Filename:         "test.pdf",
		OriginalFilename: "test.pdf",
		ContentType:      "application/pdf",
		Size:             1024,
		StoragePath:      "/uploads/test.pdf",
		StorageType:      "local",
		CreatedAt:        time.Now().Unix(),
	}
}

func TestRepository_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := NewMockRepository(ctrl)
	ctx := context.Background()
	attachment := testFileAttachment()

	repo.EXPECT().Create(ctx, attachment).Return(nil)

	err := repo.Create(ctx, attachment)

	require.NoError(t, err)
}

func TestRepository_GetByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := NewMockRepository(ctrl)
	ctx := context.Background()
	expected := testFileAttachment()

	repo.EXPECT().GetByID(ctx, testFileID).Return(expected, nil)

	attachment, err := repo.GetByID(ctx, testFileID)

	require.NoError(t, err)
	assert.Equal(t, expected, attachment)
}

func TestRepository_GetByID_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := NewMockRepository(ctrl)
	ctx := context.Background()

	repo.EXPECT().GetByID(ctx, testFileID).Return(nil, common.ErrNotFound)

	attachment, err := repo.GetByID(ctx, testFileID)

	assert.Nil(t, attachment)
	assert.Error(t, err)
	assert.True(t, errors.Is(err, common.ErrNotFound))
}

func TestRepository_GetByProject(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := NewMockRepository(ctrl)
	ctx := context.Background()
	expected := []FileAttachment{*testFileAttachment()}

	repo.EXPECT().GetByProject(ctx, testProjectID).Return(expected, nil)

	attachments, err := repo.GetByProject(ctx, testProjectID)

	require.NoError(t, err)
	assert.Equal(t, expected, attachments)
}

func TestRepository_GetByTask(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := NewMockRepository(ctrl)
	ctx := context.Background()
	expected := []FileAttachment{*testFileAttachment()}

	repo.EXPECT().GetByTask(ctx, testTaskID).Return(expected, nil)

	attachments, err := repo.GetByTask(ctx, testTaskID)

	require.NoError(t, err)
	assert.Equal(t, expected, attachments)
}

func TestRepository_Delete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := NewMockRepository(ctrl)
	ctx := context.Background()

	repo.EXPECT().Delete(ctx, testFileID).Return(nil)

	err := repo.Delete(ctx, testFileID)

	require.NoError(t, err)
}
