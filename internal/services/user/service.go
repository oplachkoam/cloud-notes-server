package user

import (
	"context"
	"fmt"

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

func (s *service) GetProfile(
	ctx context.Context, userID uuid.UUID) (*GetProfileOutput, error) {
	const op = "services.user.GetProfile"
	_ = s.log.With(logger.String("op", op))

	user, err := s.st.Users().GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &GetProfileOutput{
		Login:     user.Login,
		FirstName: user.FirstName,
		Timezone:  user.Timezone,
		CreatedAt: user.CreatedAt,
	}, nil
}

func (s *service) UpdateProfile(
	ctx context.Context, input *UpdateProfileInput) error {
	const op = "services.user.UpdateProfile"
	_ = s.log.With(logger.String("op", op))

	user, err := s.st.Users().GetByID(ctx, input.UserID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	user.FirstName = input.FirstName
	user.Timezone = input.Timezone
	err = s.st.Users().Update(ctx, user)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil

}

func (s *service) DeleteProfile(ctx context.Context, userID uuid.UUID) error {
	const op = "services.user.DeleteProfile"
	_ = s.log.With(logger.String("op", op))

	user, err := s.st.Users().GetByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	err = s.st.Users().Delete(ctx, user.ID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
