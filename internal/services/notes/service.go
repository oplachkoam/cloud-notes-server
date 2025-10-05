package user

import (
	"context"
	"fmt"
	"time"

	"cloud-notes/internal/logger"
	"cloud-notes/internal/storage"

	"github.com/google/uuid"
)

type service struct {
	log logger.Logger
	st  storage.Storage
}

func New(log logger.Logger, st storage.Storage) Service {
	return &service{
		log: log,
		st:  st,
	}
}

func (s *service) CreateNote(
	ctx context.Context, input *CreateNoteInput) (*NoteOutput, error) {
	const op = "services.notes.CreateNote"
	_ = s.log.With(logger.String("op", op))

	note := &storage.Note{
		ID:        uuid.New(),
		UserID:    input.UserID,
		Title:     input.Title,
		Text:      input.Text,
		Pinned:    input.Pinned,
		UpdatedAt: nil,
		CreatedAt: time.Now(),
	}

	err := s.st.Notes().Create(ctx, note)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &NoteOutput{
		ID:        note.ID,
		Title:     note.Title,
		Text:      note.Text,
		Pinned:    note.Pinned,
		UpdatedAt: note.UpdatedAt,
		CreatedAt: note.CreatedAt,
	}, nil
}

func (s *service) GetNotes(
	ctx context.Context, userID uuid.UUID) (*GetNotesOutput, error) {
	const op = "services.notes.GetNotes"
	_ = s.log.With(logger.String("op", op))

	notes, err := s.st.Notes().GetByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	output := new(GetNotesOutput)
	for _, note := range notes {
		output.Notes = append(output.Notes, &NoteOutput{
			ID:        note.ID,
			Title:     note.Title,
			Text:      note.Text,
			Pinned:    note.Pinned,
			UpdatedAt: note.UpdatedAt,
			CreatedAt: note.CreatedAt,
		})
	}

	return output, nil
}

func (s *service) UpdateNote(
	ctx context.Context, input *UpdateNoteInput) (*NoteOutput, error) {
	const op = "services.notes.UpdateNote"
	_ = s.log.With(logger.String("op", op))

	note, err := s.st.Notes().GetByID(ctx, input.NoteID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if note == nil || note.UserID != input.UserID {
		return nil, ErrNoteNotFound
	}

	updatedAt := time.Now()
	note.Title = input.Title
	note.Text = input.Text
	note.Pinned = input.Pinned
	note.UpdatedAt = &updatedAt
	err = s.st.Notes().Update(ctx, note)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &NoteOutput{
		ID:        note.ID,
		Title:     note.Title,
		Text:      note.Text,
		Pinned:    note.Pinned,
		UpdatedAt: note.UpdatedAt,
		CreatedAt: note.CreatedAt,
	}, nil
}

func (s *service) DeleteNote(
	ctx context.Context, input *DeleteNoteInput) error {
	const op = "services.notes.DeleteNote"
	_ = s.log.With(logger.String("op", op))

	note, err := s.st.Notes().GetByID(ctx, input.NoteID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if note == nil || note.UserID != input.UserID {
		return ErrNoteNotFound
	}

	err = s.st.Notes().Delete(ctx, note.ID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
