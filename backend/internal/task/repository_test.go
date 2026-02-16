package task

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

var testBoardID = uuid.MustParse("11111111-1111-1111-1111-111111111111")
var testTaskID = uuid.MustParse("22222222-2222-2222-2222-222222222222")
var testMemberID = uuid.MustParse("33333333-3333-3333-3333-333333333333")

func testTask() *Task {
	now := time.Now()
	return &Task{
		ID:          testTaskID,
		BoardID:     testBoardID,
		Title:       "Test Task",
		Description: "A test task",
		Status:      StatusTodo,
		Position:    0,
		ReporterID:  testMemberID,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

func testCreateTaskInput() CreateTaskInput {
	return CreateTaskInput{
		Title:      "Test Task",
		ReporterID: testMemberID,
		Status:     StatusTodo,
	}
}

func TestRepository_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := NewMockRepository(ctrl)
	ctx := context.Background()
	task := testTask()

	repo.EXPECT().Create(ctx, task).Return(nil)

	err := repo.Create(ctx, task)

	require.NoError(t, err)
}

func TestRepository_GetByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := NewMockRepository(ctrl)
	ctx := context.Background()
	expected := testTask()

	repo.EXPECT().GetByID(ctx, testTaskID).Return(expected, nil)

	task, err := repo.GetByID(ctx, testTaskID)

	require.NoError(t, err)
	assert.Equal(t, expected, task)
}

func TestRepository_GetByID_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := NewMockRepository(ctrl)
	ctx := context.Background()

	repo.EXPECT().GetByID(ctx, testTaskID).Return(nil, common.ErrNotFound)

	task, err := repo.GetByID(ctx, testTaskID)

	assert.Nil(t, task)
	assert.Error(t, err)
	assert.True(t, errors.Is(err, common.ErrNotFound))
}

func TestRepository_GetByBoard(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := NewMockRepository(ctrl)
	ctx := context.Background()
	expected := []Task{*testTask()}

	repo.EXPECT().GetByBoard(ctx, testBoardID, TaskFilter{}).Return(expected, nil)

	tasks, err := repo.GetByBoard(ctx, testBoardID, TaskFilter{})

	require.NoError(t, err)
	assert.Equal(t, expected, tasks)
}

func TestRepository_Update(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := NewMockRepository(ctrl)
	ctx := context.Background()
	task := testTask()
	task.Title = "Updated Task"

	repo.EXPECT().Update(ctx, task).Return(nil)

	err := repo.Update(ctx, task)

	require.NoError(t, err)
}

func TestRepository_Delete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := NewMockRepository(ctrl)
	ctx := context.Background()

	repo.EXPECT().Delete(ctx, testTaskID).Return(nil)

	err := repo.Delete(ctx, testTaskID)

	require.NoError(t, err)
}

func TestTaskStatus_IsValid(t *testing.T) {
	tests := []struct {
		status TaskStatus
		valid  bool
	}{
		{StatusBacklog, true},
		{StatusTodo, true},
		{StatusInProgress, true},
		{StatusCodeReview, true},
		{StatusInTest, true},
		{StatusNeedDeploy, true},
		{StatusDone, true},
		{"invalid", false},
		{"", false},
	}

	for _, tt := range tests {
		t.Run(string(tt.status), func(t *testing.T) {
			assert.Equal(t, tt.valid, tt.status.IsValid())
		})
	}
}

func TestTaskStatusOrder(t *testing.T) {
	assert.Less(t, TaskStatusOrder[StatusBacklog], TaskStatusOrder[StatusTodo])
	assert.Less(t, TaskStatusOrder[StatusTodo], TaskStatusOrder[StatusInProgress])
	assert.Less(t, TaskStatusOrder[StatusInProgress], TaskStatusOrder[StatusDone])
}
