package fixtures

import (
	"time"

	"github.com/google/uuid"
	"github.com/monachy/projek/internal/project"
)

var (
	ProjectID = uuid.MustParse("66666666-6666-6666-6666-666666666666")
	EpicID    = uuid.MustParse("77777777-7777-7777-7777-777777777777")
	BoardID   = uuid.MustParse("88888888-8888-8888-8888-888888888888")
)

func Project() *project.Project {
	now := time.Now()
	return &project.Project{
		ID:          ProjectID,
		CommunityID: CommunityID,
		Name:        "Test Project",
		Description: "A test project",
		Key:         "TEST",
		IsArchived:  false,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

func Epic() *project.Epic {
	now := time.Now()
	return &project.Epic{
		ID:          EpicID,
		ProjectID:   ProjectID,
		Name:        "Test Epic",
		Description: "A test epic",
		StartDate:   nil,
		EndDate:     nil,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

func Board() *project.Board {
	now := time.Now()
	return &project.Board{
		ID:        BoardID,
		ProjectID: ProjectID,
		Name:      "Test Board",
		CreatedAt: now,
		UpdatedAt: now,
	}
}

func CreateProjectInput() project.CreateProjectInput {
	return project.CreateProjectInput{
		Name:        "Test Project",
		Description: "A test project",
		Key:         "TEST",
	}
}

func CreateEpicInput() project.CreateEpicInput {
	return project.CreateEpicInput{
		Name:        "Test Epic",
		Description: "A test epic",
		StartDate:   nil,
		EndDate:     nil,
	}
}

func CreateBoardInput() project.CreateBoardInput {
	return project.CreateBoardInput{
		Name: "Test Board",
	}
}
