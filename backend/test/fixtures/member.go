package fixtures

import (
	"time"

	"github.com/google/uuid"
	"github.com/monachy/projek/internal/member"
)

var (
	CommunityID = uuid.MustParse("11111111-1111-1111-1111-111111111111")
	MemberID    = uuid.MustParse("22222222-2222-2222-2222-222222222222")
	MemberID2   = uuid.MustParse("33333333-3333-3333-3333-333333333333")
)

func Member() *member.Member {
	now := time.Now()
	return &member.Member{
		ID:           MemberID,
		CommunityID:  CommunityID,
		Email:        "john@example.com",
		PasswordHash: "$2a$10$hashedpassword",
		FirstName:    "John",
		LastName:     "Doe",
		DateOfBirth:  nil,
		Role:         member.RoleMember,
		IsActive:     true,
		LastLoginAt:  nil,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
}

func MemberAdmin() *member.Member {
	m := Member()
	m.ID = uuid.MustParse("44444444-4444-4444-4444-444444444444")
	m.Email = "admin@example.com"
	m.Role = member.RoleAdmin
	return m
}

func CreateMemberInput() member.CreateMemberInput {
	return member.CreateMemberInput{
		Email:       "john@example.com",
		Password:    "password123",
		FirstName:   "John",
		LastName:    "Doe",
		DateOfBirth: nil,
		Role:        member.RoleMember,
	}
}

func UpdateMemberInput() member.UpdateMemberInput {
	trueVal := true
	return member.UpdateMemberInput{
		FirstName:   "JohnUpdated",
		LastName:    "DoeUpdated",
		DateOfBirth: nil,
		Role:        member.RoleAdmin,
		IsActive:    &trueVal,
	}
}
