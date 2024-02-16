package handler

import (
	"database/sql"
	"errors"
	"github.com/dfalgout/dream/config"
	"log/slog"
	"net/http"
	"time"

	"github.com/dfalgout/dream/helpers"
	"github.com/dfalgout/dream/service"
	"github.com/dfalgout/dream/types"
	"github.com/dfalgout/dream/view/user"
	"github.com/labstack/echo/v4"
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
	e.POST(config.ActionSendCode, h.sendCode)
	e.POST(config.ActionVerifyCode, h.verifyCode)
	e.PATCH(config.ActionUsers, h.updateUser)

	e.GET(config.Onboarding, h.Onboarding)

	e.GET(config.LoginPage, h.loginPage)
	e.GET(config.Logout, h.logout)
}

type updateUserInput struct {
	FullName *string `form:"fullName"`
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
	Email string `form:"email" validate:"required,email"`
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
	return Render(c, user.Verify(input.Email))
}

type verifyCodeInput struct {
	Email string `form:"email" validate:"required,email"`
	Code  string `form:"code" validate:"required,len=6"`
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
		if errors.Is(err, sql.ErrNoRows) {
			h.logger.Error("failed to send code", "err", err)
			return c.NoContent(http.StatusPreconditionFailed)
		}
		h.logger.Error("failed to send code", "err", err)
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
	return Navigate(c, config.Onboarding)
}

func (h *userHandlers) loginPage(c echo.Context) error {
	return Render(c, user.Login())
}

func (h *userHandlers) Onboarding(c echo.Context) error {
	return Render(c, user.OnboardingPage())
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
	return Navigate(c, config.LoginPage)
}
