package fixtures

import (
	"time"

	"github.com/google/uuid"
	"github.com/monachy/projek/internal/community"
)

var (
	CommunityID2 = uuid.MustParse("55555555-5555-5555-5555-555555555555")
)

func Community() *community.Community {
	now := time.Now()
	return &community.Community{
		ID:        CommunityID,
		Name:      "Test Community",
		Slug:      "test-community",
		CreatedAt: now,
		UpdatedAt: now,
	}
}

func Community2() *community.Community {
	now := time.Now()
	return &community.Community{
		ID:        CommunityID2,
		Name:      "Second Community",
		Slug:      "second-community",
		CreatedAt: now,
		UpdatedAt: now,
	}
}

func CreateCommunityInput() community.CreateCommunityInput {
	return community.CreateCommunityInput{
		Name: "Test Community",
		Slug: "test-community",
	}
}
