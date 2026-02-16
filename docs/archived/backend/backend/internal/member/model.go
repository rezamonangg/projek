package member

import (
	"time"

	"github.com/google/uuid"
)

type Member struct {
	ID           uuid.UUID  `json:"id"`
	CommunityID  uuid.UUID  `json:"community_id"`
	Email        string     `json:"email"`
	PasswordHash string     `json:"-"`
	FirstName    string     `json:"first_name"`
	LastName     string     `json:"last_name"`
	DateOfBirth  *time.Time `json:"date_of_birth,omitempty"`
	Role         Role       `json:"role"`
	IsActive     bool       `json:"is_active"`
	LastLoginAt  *time.Time `json:"last_login_at,omitempty"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

type Role string

const (
	RoleAdmin  Role = "admin"
	RoleMember Role = "member"
)

func (r Role) String() string {
	return string(r)
}

func (r Role) IsValid() bool {
	switch r {
	case RoleAdmin, RoleMember:
		return true
	}
	return false
}

type Invitation struct {
	ID          uuid.UUID  `json:"id"`
	CommunityID uuid.UUID  `json:"community_id"`
	Email       string     `json:"email"`
	Token       string     `json:"token,omitempty"`
	Role        Role       `json:"role"`
	InvitedBy   uuid.UUID  `json:"invited_by"`
	AcceptedAt  *time.Time `json:"accepted_at,omitempty"`
	ExpiresAt   time.Time  `json:"expires_at"`
	CreatedAt   time.Time  `json:"created_at"`
}

type CreateMemberInput struct {
	Email       string     `json:"email" validate:"required,email"`
	Password    string     `json:"password" validate:"required,min=8"`
	FirstName   string     `json:"first_name" validate:"required"`
	LastName    string     `json:"last_name"`
	DateOfBirth *time.Time `json:"date_of_birth"`
	Role        Role       `json:"role"`
}

type UpdateMemberInput struct {
	FirstName   string     `json:"first_name"`
	LastName    string     `json:"last_name"`
	DateOfBirth *time.Time `json:"date_of_birth"`
	Role        Role       `json:"role"`
	IsActive    *bool      `json:"is_active"`
}
