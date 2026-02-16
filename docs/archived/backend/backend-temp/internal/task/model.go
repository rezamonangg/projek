package task

import (
	"time"

	"github.com/google/uuid"
)

type Task struct {
	ID          uuid.UUID  `json:"id"`
	BoardID     uuid.UUID  `json:"board_id"`
	EpicID      *uuid.UUID `json:"epic_id,omitempty"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Status      TaskStatus `json:"status"`
	Position    int        `json:"position"`
	StoryPoints *int       `json:"story_points,omitempty"`
	DueDate     *time.Time `json:"due_date,omitempty"`
	AssigneeID  *uuid.UUID `json:"assignee_id,omitempty"`
	ReporterID  uuid.UUID  `json:"reporter_id"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

type TaskStatus string

const (
	StatusBacklog    TaskStatus = "backlog"
	StatusTodo       TaskStatus = "todo"
	StatusInProgress TaskStatus = "inprogress"
	StatusCodeReview TaskStatus = "codereview"
	StatusInTest     TaskStatus = "intest"
	StatusNeedDeploy TaskStatus = "needdeploy"
	StatusDone       TaskStatus = "done"
)

func (s TaskStatus) IsValid() bool {
	switch s {
	case StatusBacklog, StatusTodo, StatusInProgress, StatusCodeReview, StatusInTest, StatusNeedDeploy, StatusDone:
		return true
	}
	return false
}

var TaskStatusOrder = map[TaskStatus]int{
	StatusBacklog:    0,
	StatusTodo:       1,
	StatusInProgress: 2,
	StatusCodeReview: 3,
	StatusInTest:     4,
	StatusNeedDeploy: 5,
	StatusDone:       6,
}

type Label struct {
	ID        uuid.UUID `json:"id"`
	ProjectID uuid.UUID `json:"project_id"`
	Name      string    `json:"name"`
	Color     string    `json:"color"`
	CreatedAt time.Time `json:"created_at"`
}

type TaskLabel struct {
	TaskID  uuid.UUID `json:"task_id"`
	LabelID uuid.UUID `json:"label_id"`
}

type CreateTaskInput struct {
	Title       string     `json:"title" validate:"required"`
	Description string     `json:"description"`
	EpicID      *uuid.UUID `json:"epic_id"`
	Status      TaskStatus `json:"status"`
	StoryPoints *int       `json:"story_points"`
	DueDate     *time.Time `json:"due_date"`
	AssigneeID  *uuid.UUID `json:"assignee_id"`
	ReporterID  uuid.UUID  `json:"reporter_id" validate:"required"`
}

type UpdateTaskInput struct {
	Title       string     `json:"title"`
	Description string     `json:"description"`
	EpicID      *uuid.UUID `json:"epic_id"`
	Status      TaskStatus `json:"status"`
	Position    *int       `json:"position"`
	StoryPoints *int       `json:"story_points"`
	DueDate     *time.Time `json:"due_date"`
	AssigneeID  *uuid.UUID `json:"assignee_id"`
}

type CreateLabelInput struct {
	Name  string `json:"name" validate:"required"`
	Color string `json:"color"`
}

type TaskFilter struct {
	Status     TaskStatus
	AssigneeID uuid.UUID
	EpicID     uuid.UUID
	Search     string
	Limit      int
	Offset     int
}
