package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/dfalgout/dream/db"
	"github.com/dfalgout/dream/helpers"
	"github.com/dfalgout/dream/types"
	"github.com/lucsky/cuid"
)

type UserService interface {
	UpsertUserById(ctx context.Context, input *types.UpsertUserIdInput) (*types.User, error)
	UpsertUserByEmail(ctx context.Context, input *types.UpsertUserEmailInput) (*types.User, error)
	GetUserByEmail(ctx context.Context, email string) (*types.User, error)
	SendCode(ctx context.Context, email string) error
	VerifyCode(ctx context.Context, input *types.VerifyCodeInput) (*string, error)
}

type userService struct {
	queries *db.Queries
	logger  *slog.Logger
}

func NewUserService(queries *db.Queries, logger *slog.Logger) UserService {
	return &userService{
		queries: queries,
		logger:  logger,
	}
}

func (s *userService) UpsertUserById(ctx context.Context, input *types.UpsertUserIdInput) (*types.User, error) {
	checkedEmail := helpers.CleanEmail(input.Email)
	if checkedEmail == "" {
		return nil, fmt.Errorf("email is required")
	}
	foundUser, err := s.queries.UpdateUserById(ctx, &db.UpdateUserByIdParams{
		ID:              input.ID,
		Email:           &checkedEmail,
		FullName:        input.FullName,
		IsAdmin:         input.IsAdmin,
		VerifyCode:      input.VerifyCode,
		VerifyExpiresAt: input.VerifyExpiresAt,
	})
	if err != nil {
		if errors.Is(sql.ErrNoRows, err) {
			foundUser, err = s.queries.CreateUser(ctx, &db.CreateUserParams{
				ID:              cuid.New(),
				Email:           checkedEmail,
				FullName:        input.FullName,
				IsAdmin:         input.IsAdmin,
				VerifyCode:      input.VerifyCode,
				VerifyExpiresAt: input.VerifyExpiresAt,
			})
			if err != nil {
				s.logger.Error("failed to update user", "err", err)
				return nil, err
			}
		} else {
			s.logger.Error("failed to get user", "err", err)
			return nil, err
		}
	}
	return types.NewUser(foundUser), nil
}

func (s *userService) UpsertUserByEmail(ctx context.Context, input *types.UpsertUserEmailInput) (*types.User, error) {
	checkedEmail := helpers.CleanEmail(&input.Email)
	if checkedEmail == "" {
		return nil, fmt.Errorf("email is required")
	}
	foundUser, err := s.queries.UpdateUserByEmail(ctx, &db.UpdateUserByEmailParams{
		Email:           checkedEmail,
		FullName:        input.FullName,
		IsAdmin:         input.IsAdmin,
		VerifyCode:      input.VerifyCode,
		VerifyExpiresAt: input.VerifyExpiresAt,
	})
	if err != nil {
		if errors.Is(sql.ErrNoRows, err) {
			foundUser, err = s.queries.CreateUser(ctx, &db.CreateUserParams{
				ID:              cuid.New(),
				Email:           checkedEmail,
				FullName:        input.FullName,
				IsAdmin:         input.IsAdmin,
				VerifyCode:      input.VerifyCode,
				VerifyExpiresAt: input.VerifyExpiresAt,
			})
			if err != nil {
				s.logger.Error("failed to update user", "err", err)
				return nil, err
			}
		} else {
			s.logger.Error("failed to get user", "err", err)
			return nil, err
		}
	}
	return types.NewUser(foundUser), nil
}

func (s *userService) GetUserByEmail(ctx context.Context, email string) (*types.User, error) {
	checkedEmail := helpers.CleanEmail(&email)
	foundUser, err := s.queries.GetUserByEmail(ctx, checkedEmail)
	if err != nil {
		s.logger.Error("failed to get user", "err", err)
		return nil, err
	}
	return types.NewUser(foundUser), nil
}

func (s *userService) SendCode(ctx context.Context, email string) error {
	code := helpers.Code()
	now := time.Now()
	expireCode := now.Add(time.Minute * 5).UTC()
	foundUser, err := s.UpsertUserByEmail(ctx, &types.UpsertUserEmailInput{
		Email:           email,
		VerifyCode:      &code,
		VerifyExpiresAt: &expireCode,
	})
	if err != nil {
		s.logger.Error("failed to upsert user", "error", err)
		return err
	}
	s.logger.Info("code sent", slog.String("email", foundUser.Email), slog.String("code", code))
	return nil
}

func (s *userService) VerifyCode(ctx context.Context, input *types.VerifyCodeInput) (*string, error) {
	checkedEmail := strings.ToLower(input.Email)
	foundUser, err := s.queries.ClearAndGetVerifiedUser(ctx, &db.ClearAndGetVerifiedUserParams{
		Email:      checkedEmail,
		VerifyCode: &input.Code,
	})
	if err != nil {
		s.logger.Error("failed to get user", "err", err)
		return nil, err
	}
	token, err := helpers.CreateToken(foundUser)
	if err != nil {
		s.logger.Error("failed to generate token", "error", err)
		return nil, err
	}
	return token, nil
}
