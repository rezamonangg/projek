package wiki

import (
	"time"

	"github.com/google/uuid"
)

type WikiPage struct {
	ID        uuid.UUID  `json:"id"`
	ProjectID uuid.UUID  `json:"project_id"`
	ParentID  *uuid.UUID `json:"parent_id,omitempty"`
	Title     string     `json:"title"`
	Slug      string     `json:"slug"`
	Content   string     `json:"content"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

type WikiPageVersion struct {
	ID         uuid.UUID `json:"id"`
	WikiPageID uuid.UUID `json:"wiki_page_id"`
	Content    string    `json:"content"`
	Version    int       `json:"version"`
	CreatedAt  time.Time `json:"created_at"`
	CreatedBy  uuid.UUID `json:"created_by"`
}

type CreateWikiPageInput struct {
	ParentID *uuid.UUID `json:"parent_id"`
	Title    string     `json:"title" validate:"required"`
	Content  string     `json:"content"`
}

type UpdateWikiPageInput struct {
	Title    string     `json:"title"`
	Content  string     `json:"content"`
	ParentID *uuid.UUID `json:"parent_id"`
}
