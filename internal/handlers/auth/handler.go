package auth

import (
	"encoding/json"
	"errors"
	"net/http"

	"cloud-notes/internal/logger"
	"cloud-notes/internal/render"
	"cloud-notes/internal/security"
	"cloud-notes/internal/services/auth"

	"github.com/go-playground/validator/v10"
)

type Handler struct {
	log logger.Logger
	srv auth.Service
	val *validator.Validate
}

func New(log logger.Logger, srv auth.Service) Handler {
	return Handler{
		log: log,
		srv: srv,
		val: validator.New(),
	}
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.auth.Register"
	_ = h.log.With(logger.String("op", op))
	ctx := r.Context()

	request := new(RegisterRequest)
	if err := json.NewDecoder(r.Body).Decode(request); err != nil {
		render.InvalidJSONError(w)
		return
	}

	if err := h.val.StructCtx(ctx, request); err != nil {
		render.ValidationError(w, request, err)
		return
	}

	err := h.srv.Register(ctx, &auth.RegisterInput{
		Login:     request.Login,
		Password:  request.Password,
		FirstName: request.FirstName,
		Timezone:  request.Timezone,
	})

	switch {
	case err == nil:
		render.Empty(w)
	case errors.Is(err, auth.ErrLoginAlreadyExists):
		render.Error(w, http.StatusConflict, err)
	default:
		render.ServerError(w, http.StatusInternalServerError)
	}
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.auth.Login"
	_ = h.log.With(logger.String("op", op))
	ctx := r.Context()

	request := new(LoginRequest)
	if err := json.NewDecoder(r.Body).Decode(request); err != nil {
		render.InvalidJSONError(w)
		return
	}

	if err := h.val.StructCtx(ctx, request); err != nil {
		render.ValidationError(w, request, err)
		return
	}

	input := &auth.LoginInput{
		Login:    request.Login,
		Password: request.Password,
	}
	if r.UserAgent() != "" {
		userAgent := r.UserAgent()
		input.UserAgent = &userAgent
	}

	output, err := h.srv.Login(ctx, input)

	switch {
	case err == nil:
		render.JSON(w, http.StatusOK, LoginResponse{
			AccessToken: output.AccessToken,
		})
	case errors.Is(err, auth.ErrUserNotFound):
		render.Error(w, http.StatusNotFound, err)
	case errors.Is(err, auth.ErrInvalidPassword):
		render.Error(w, http.StatusForbidden, err)
	default:
		render.ServerError(w, http.StatusInternalServerError)
	}
}

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.auth.Logout"
	_ = h.log.With(logger.String("op", op))
	ctx := r.Context()

	claims := security.GetClaims(ctx)
	err := h.srv.Logout(ctx, claims.SessionID)

	switch { // nolint
	case err == nil:
		render.Empty(w)
	default:
		render.ServerError(w, http.StatusInternalServerError)
	}
}
func (h *Handler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.auth.ChangePassword"
	_ = h.log.With(logger.String("op", op))
	ctx := r.Context()

	request := new(ChangePasswordRequest)
	if err := json.NewDecoder(r.Body).Decode(request); err != nil {
		render.InvalidJSONError(w)
		return
	}

	if err := h.val.StructCtx(ctx, request); err != nil {
		render.ValidationError(w, request, err)
		return
	}

	claims := security.GetClaims(ctx)
	err := h.srv.ChangePassword(ctx, &auth.ChangePasswordInput{
		UserID:      claims.UserID,
		OldPassword: request.OldPassword,
		NewPassword: request.NewPassword,
	})

	switch {
	case err == nil:
		render.Empty(w)
	case errors.Is(err, auth.ErrInvalidPassword):
		render.Error(w, http.StatusForbidden, err)
	default:
		render.ServerError(w, http.StatusInternalServerError)
	}
}
