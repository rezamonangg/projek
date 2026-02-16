package fixtures

import (
	"time"

	"github.com/google/uuid"
	"github.com/monachy/projek/internal/task"
)

var (
	TaskID  = uuid.MustParse("99999999-9999-9999-9999-999999999999")
	LabelID = uuid.MustParse("aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa")
)

func Task() *task.Task {
	now := time.Now()
	return &task.Task{
		ID:          TaskID,
		BoardID:     BoardID,
		EpicID:      nil,
		Title:       "Test Task",
		Description: "A test task",
		Status:      task.StatusTodo,
		Position:    0,
		StoryPoints: nil,
		DueDate:     nil,
		AssigneeID:  nil,
		ReporterID:  MemberID,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

func TaskWithEpic() *task.Task {
	t := Task()
	t.EpicID = &EpicID
	return t
}

func Label() *task.Label {
	now := time.Now()
	return &task.Label{
		ID:        LabelID,
		ProjectID: ProjectID,
		Name:      "Bug",
		Color:     "red",
		CreatedAt: now,
	}
}

func CreateTaskInput() task.CreateTaskInput {
	return task.CreateTaskInput{
		Title:       "Test Task",
		Description: "A test task",
		EpicID:      nil,
		Status:      task.StatusTodo,
		StoryPoints: nil,
		DueDate:     nil,
		AssigneeID:  nil,
		ReporterID:  MemberID,
	}
}

func CreateLabelInput() task.CreateLabelInput {
	return task.CreateLabelInput{
		Name:  "Bug",
		Color: "red",
	}
}
