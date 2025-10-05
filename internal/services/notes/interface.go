package user

import (
	"context"

	"github.com/google/uuid"
)

type Service interface {
	CreateNote(ctx context.Context, input *CreateNoteInput) (*NoteOutput, error)
	GetNotes(ctx context.Context, userID uuid.UUID) (*GetNotesOutput, error)
	UpdateNote(ctx context.Context, input *UpdateNoteInput) (*NoteOutput, error)
	DeleteNote(ctx context.Context, input *DeleteNoteInput) error
}
