package user

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrNoteNotFound = errors.New("note not found")
)

type NoteOutput struct {
	ID        uuid.UUID
	Title     *string
	Text      *string
	Pinned    bool
	UpdatedAt *time.Time
	CreatedAt time.Time
}

type CreateNoteInput struct {
	UserID uuid.UUID
	Title  *string
	Text   *string
	Pinned bool
}

type GetNotesOutput struct {
	Notes []*NoteOutput
}

type UpdateNoteInput struct {
	UserID uuid.UUID
	NoteID uuid.UUID
	Title  *string
	Text   *string
	Pinned bool
}

type DeleteNoteInput struct {
	UserID uuid.UUID
	NoteID uuid.UUID
}
