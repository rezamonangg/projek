package community

import (
	"time"

	"github.com/google/uuid"
)

type Community struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Slug      string    `json:"slug"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateCommunityInput struct {
	Name string `json:"name" validate:"required,min=2,max=100"`
	Slug string `json:"slug" validate:"required,alphanum,min=2,max=50"`
}

type CommunitySettings struct {
	ID                       uuid.UUID `json:"id"`
	CommunityID              uuid.UUID `json:"community_id"`
	AllowMemberRegistration  bool      `json:"allow_member_registration"`
	RequireEmailVerification bool      `json:"require_email_verification"`
	CreatedAt                time.Time `json:"created_at"`
	UpdatedAt                time.Time `json:"updated_at"`
}
