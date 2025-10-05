package notes

import (
	"time"

	"github.com/google/uuid"
)

type NoteResponse struct {
	ID        uuid.UUID  `json:"id"`
	Title     *string    `json:"title"`
	Text      *string    `json:"text"`
	Pinned    bool       `json:"pinned"`
	UpdatedAt *time.Time `json:"updated_at"`
	CreatedAt time.Time  `json:"created_at"`
}

type CreateNoteRequest struct {
	Title  *string `json:"title" validate:"required,min=1,max=1000"`
	Text   *string `json:"text" validate:"required,min=1,max=10000"`
	Pinned bool    `json:"pinned"`
}

type GetNotesResponse struct {
	Notes []*NoteResponse `json:"notes"`
}

type UpdateNoteRequest struct {
	Title  *string `json:"title" validate:"required,min=1,max=1000"`
	Text   *string `json:"text" validate:"required,min=1,max=10000"`
	Pinned bool    `json:"pinned"`
}
