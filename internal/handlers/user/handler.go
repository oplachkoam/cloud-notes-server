package user

import (
	"encoding/json"
	"net/http"

	"cloud-notes/internal/logger"
	"cloud-notes/internal/render"
	"cloud-notes/internal/security"
	"cloud-notes/internal/services/user"

	"github.com/go-playground/validator/v10"
)

type Handler struct {
	log logger.Logger
	srv user.Service
	val *validator.Validate
}

func New(log logger.Logger, srv user.Service) Handler {
	return Handler{
		log: log,
		srv: srv,
		val: validator.New(),
	}
}

func (h *Handler) GetProfile(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.user.GetProfile"
	_ = h.log.With(logger.String("op", op))
	ctx := r.Context()

	claims := security.GetClaims(ctx)
	output, err := h.srv.GetProfile(ctx, claims.UserID)

	switch { // nolint
	case err == nil:
		render.JSON(w, http.StatusOK, &GetProfileResponse{
			Login:     output.Login,
			FirstName: output.FirstName,
			Timezone:  output.Timezone,
			CreatedAt: output.CreatedAt,
		})
	default:
		render.ServerError(w, http.StatusInternalServerError)
	}
}

func (h *Handler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.user.UpdateProfile"
	_ = h.log.With(logger.String("op", op))
	ctx := r.Context()

	request := new(UpdateProfileRequest)
	if err := json.NewDecoder(r.Body).Decode(request); err != nil {
		render.InvalidJSONError(w)
		return
	}

	if err := h.val.StructCtx(ctx, request); err != nil {
		render.ValidationError(w, request, err)
		return
	}

	claims := security.GetClaims(ctx)
	err := h.srv.UpdateProfile(ctx, &user.UpdateProfileInput{
		UserID:    claims.UserID,
		FirstName: request.FirstName,
		Timezone:  request.Timezone,
	})

	switch { // nolint
	case err == nil:
		render.Empty(w)
	default:
		render.ServerError(w, http.StatusInternalServerError)
	}
}

func (h *Handler) DeleteProfile(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.user.DeleteProfile"
	_ = h.log.With(logger.String("op", op))
	ctx := r.Context()

	claims := security.GetClaims(ctx)
	err := h.srv.DeleteProfile(ctx, claims.UserID)

	switch { // nolint
	case err == nil:
		render.Empty(w)
	default:
		render.ServerError(w, http.StatusInternalServerError)
	}
}
