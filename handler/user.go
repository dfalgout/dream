package handler

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/dfalgout/dream/helpers"
	"github.com/dfalgout/dream/service"
	"github.com/dfalgout/dream/types"
	"github.com/labstack/echo/v4"
)

const (
	ACTION_SEND_CODE   = "/action/send-code"
	ACTION_VERIFY_CODE = "/action/verify-code"
	ACTION_USERS       = "/action/users"
)

type userHandlers struct {
	logger      *slog.Logger
	userService service.UserService
}

func RegisterUserHandlers(e *echo.Echo, logger *slog.Logger, userService service.UserService) {
	h := &userHandlers{
		logger,
		userService,
	}
	// register routes
	e.POST(ACTION_SEND_CODE, h.sendCode)
	e.POST(ACTION_VERIFY_CODE, h.verifyCode)
	e.PATCH(ACTION_USERS, h.updateUser)
}

type updateUserInput struct {
	FullName *string `json:"fullName"`
}

func (h *userHandlers) updateUser(c echo.Context) error {
	requester := helpers.GetSessionUser(c)
	input := new(updateUserInput)
	if err := helpers.BindAndValidate(c, input); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	createdUser, err := h.userService.UpsertUserById(c.Request().Context(), &types.UpsertUserIdInput{
		ID:       requester.ID,
		Email:    &requester.Email,
		FullName: input.FullName,
	})
	if err != nil {
		h.logger.Error("failed to update user", "err", err)
		return c.String(http.StatusInternalServerError, "failed to create user")
	}
	return c.JSON(http.StatusCreated, createdUser)
}

type sendCodeInput struct {
	Email string `json:"email" validate:"required,email"`
}

func (h *userHandlers) sendCode(c echo.Context) error {
	input := new(sendCodeInput)
	if err := helpers.BindAndValidate(c, input); err != nil {
		return c.NoContent(http.StatusPreconditionFailed)
	}
	err := h.userService.SendCode(c.Request().Context(), input.Email)
	if err != nil {
		h.logger.Error("failed to send code", slog.Any("error", err))
		return c.NoContent(http.StatusInternalServerError)
	}
	return c.NoContent(http.StatusOK)
}

type verifyCodeInput struct {
	Email string `json:"email" validate:"required,email"`
	Code  string `json:"code" validate:"required,len=6"`
}

func (h *userHandlers) verifyCode(c echo.Context) error {
	input := new(verifyCodeInput)
	if err := helpers.BindAndValidate(c, input); err != nil {
		return c.NoContent(http.StatusPreconditionFailed)
	}
	token, err := h.userService.VerifyCode(c.Request().Context(), &types.VerifyCodeInput{
		Email: input.Email,
		Code:  input.Code,
	})
	if err != nil {
		h.logger.Error("failed to send code", slog.Any("error", err))
		return c.NoContent(http.StatusInternalServerError)
	}
	c.SetCookie(&http.Cookie{
		Name:     "token",
		Value:    *token,
		Expires:  time.Now().Add(24 * time.Hour),
		Path:     "/",
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	})
	return c.NoContent(http.StatusOK)
}

func (h *userHandlers) logout(c echo.Context) error {
	c.SetCookie(&http.Cookie{
		Name:     "token",
		Value:    "",
		Expires:  time.Now().Add(-1 * time.Second),
		Path:     "/",
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	})
	return c.NoContent(http.StatusOK)
}
