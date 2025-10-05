package auth

import (
	"context"
	"fmt"
	"time"

	"cloud-notes/internal/logger"
	"cloud-notes/internal/security"
	"cloud-notes/internal/storage"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type service struct {
	log logger.Logger
	st  storage.Storage
	sec security.Security
}

func New(log logger.Logger, st storage.Storage, sec security.Security) Service {
	return &service{
		log: log,
		st:  st,
		sec: sec,
	}
}

func (s *service) Register(ctx context.Context, input *RegisterInput) error {
	const op = "services.auth.Register"
	log := s.log.With(logger.String("op", op))

	user, err := s.st.Users().GetByLogin(ctx, input.Login)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if user != nil {
		return ErrLoginAlreadyExists
	}

	bytes, err := bcrypt.GenerateFromPassword(
		[]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		log.ErrorContext(ctx, "", logger.Error(err))
		return fmt.Errorf("%s: %w", op, err)
	}
	passwordHash := string(bytes)

	err = s.st.Users().Create(ctx, &storage.User{
		ID:           uuid.New(),
		Login:        input.Login,
		PasswordHash: passwordHash,
		FirstName:    input.FirstName,
		Timezone:     input.Timezone,
		Status:       storage.UserStatusActive,
		CreatedAt:    time.Now(),
	})
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *service) Login(
	ctx context.Context, input *LoginInput) (*LoginOutput, error) {
	const op = "services.auth.Login"
	_ = s.log.With(logger.String("op", op))

	user, err := s.st.Users().GetByLogin(ctx, input.Login)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if user == nil {
		return nil, ErrUserNotFound
	}

	err = bcrypt.CompareHashAndPassword(
		[]byte(user.PasswordHash), []byte(input.Password))
	if err != nil {
		return nil, ErrInvalidPassword
	}

	session := &storage.Session{
		ID:        uuid.New(),
		UserID:    user.ID,
		UserAgent: input.UserAgent,
		CreatedAt: time.Now(),
	}

	token := s.sec.GenerateAccessToken(ctx, &security.Claims{
		UserID:    user.ID,
		SessionID: session.ID,
		CreatedAt: session.CreatedAt,
	})

	err = s.st.Sessions().Create(ctx, session)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &LoginOutput{
		AccessToken: token,
	}, nil
}

func (s *service) Logout(ctx context.Context, sessionID uuid.UUID) error {
	const op = "services.auth.Logout"
	_ = s.log.With(logger.String("op", op))

	err := s.st.Sessions().Delete(ctx, sessionID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *service) ChangePassword(
	ctx context.Context, input *ChangePasswordInput) error {
	const op = "services.auth.ChangePassword"
	log := s.log.With(logger.String("op", op))

	user, err := s.st.Users().GetByID(ctx, input.UserID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	err = bcrypt.CompareHashAndPassword(
		[]byte(user.PasswordHash), []byte(input.OldPassword))
	if err != nil {
		return ErrInvalidPassword
	}

	bytes, err := bcrypt.GenerateFromPassword(
		[]byte(input.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		log.ErrorContext(ctx, "", logger.Error(err))
		return fmt.Errorf("%s: %w", op, err)
	}
	passwordHash := string(bytes)

	user.PasswordHash = passwordHash
	err = s.st.Users().Update(ctx, user)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
