package notes

import (
	"cloud-notes/internal/logger"
	"cloud-notes/internal/render"
	"cloud-notes/internal/security"
	"cloud-notes/internal/services/notes"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type Handler struct {
	log logger.Logger
	srv notes.Service
	val *validator.Validate
}

func New(log logger.Logger, srv notes.Service) Handler {
	return Handler{
		log: log,
		srv: srv,
		val: validator.New(),
	}
}

func (h *Handler) CreateNote(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.notes.CreateNote"
	_ = h.log.With(logger.String("op", op))
	ctx := r.Context()

	request := new(CreateNoteRequest)
	if err := json.NewDecoder(r.Body).Decode(request); err != nil {
		render.InvalidJSONError(w)
		return
	}

	if err := h.val.StructCtx(ctx, request); err != nil {
		render.ValidationError(w, request, err)
		return
	}

	claims := security.GetClaims(ctx)
	output, err := h.srv.CreateNote(ctx, &notes.CreateNoteInput{
		UserID: claims.UserID,
		Title:  request.Title,
		Text:   request.Text,
		Pinned: request.Pinned,
	})

	switch { // nolint
	case err == nil:
		render.JSON(w, http.StatusOK, &NoteResponse{
			ID:        output.ID,
			Title:     output.Title,
			Text:      output.Text,
			Pinned:    output.Pinned,
			UpdatedAt: output.UpdatedAt,
			CreatedAt: output.CreatedAt,
		})
	default:
		render.ServerError(w, http.StatusInternalServerError)
	}
}

func (h *Handler) GetNotes(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.notes.GetNotes"
	_ = h.log.With(logger.String("op", op))
	ctx := r.Context()

	claims := security.GetClaims(ctx)
	output, err := h.srv.GetNotes(ctx, claims.UserID)

	switch { // nolint
	case err == nil:
		response := new(GetNotesResponse)
		for _, note := range output.Notes {
			response.Notes = append(response.Notes, &NoteResponse{
				ID:        note.ID,
				Title:     note.Title,
				Text:      note.Text,
				Pinned:    note.Pinned,
				UpdatedAt: note.UpdatedAt,
				CreatedAt: note.CreatedAt,
			})
		}
		render.JSON(w, http.StatusOK, response)
	default:
		render.ServerError(w, http.StatusInternalServerError)
	}
}

func (h *Handler) UpdateNote(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.notes.UpdateNote"
	_ = h.log.With(logger.String("op", op))
	ctx := r.Context()

	noteID, err := uuid.Parse(chi.URLParam(r, "note-id"))
	if err != nil {
		render.Error(w, http.StatusBadRequest,
			errors.New("invalid note id"))
		return
	}

	request := new(UpdateNoteRequest)
	if err := json.NewDecoder(r.Body).Decode(request); err != nil {
		render.InvalidJSONError(w)
		return
	}

	if err := h.val.StructCtx(ctx, request); err != nil {
		render.ValidationError(w, request, err)
		return
	}

	claims := security.GetClaims(ctx)
	output, err := h.srv.UpdateNote(ctx, &notes.UpdateNoteInput{
		UserID: claims.UserID,
		NoteID: noteID,
		Title:  request.Title,
		Text:   request.Text,
		Pinned: request.Pinned,
	})

	switch { // nolint
	case err == nil:
		render.JSON(w, http.StatusOK, &NoteResponse{
			ID:        output.ID,
			Title:     output.Title,
			Text:      output.Text,
			Pinned:    output.Pinned,
			UpdatedAt: output.UpdatedAt,
			CreatedAt: output.CreatedAt,
		})
	case errors.Is(err, notes.ErrNoteNotFound):
		render.Error(w, http.StatusNotFound, err)
	default:
		render.ServerError(w, http.StatusInternalServerError)
	}
}

func (h *Handler) DeleteNote(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.notes.DeleteNote"
	_ = h.log.With(logger.String("op", op))
	ctx := r.Context()

	noteID, err := uuid.Parse(chi.URLParam(r, "note-id"))
	if err != nil {
		render.Error(w, http.StatusBadRequest,
			errors.New("invalid note id"))
		return
	}

	claims := security.GetClaims(ctx)
	err = h.srv.DeleteNote(ctx, &notes.DeleteNoteInput{
		UserID: claims.UserID,
		NoteID: noteID,
	})

	switch { // nolint
	case err == nil:
		render.Empty(w)
	case errors.Is(err, notes.ErrNoteNotFound):
		render.Error(w, http.StatusNotFound, err)
	default:
		render.ServerError(w, http.StatusInternalServerError)
	}
}
