package handler

import (
	"github.com/a-h/templ"
	"github.com/dfalgout/dream/config"
	"github.com/dfalgout/dream/helpers"
	"github.com/dfalgout/dream/types"
	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"net/http"
	"slices"
)

var (
	PublicRoutes = []string{config.ActionSendCode, config.ActionVerifyCode, config.LoginPage}
)

func Render(c echo.Context, component templ.Component) error {
	c.Response().Header().Set("Content-Type", "text/html")
	return component.Render(c.Request().Context(), c.Response().Writer)
}

func Redirect(c echo.Context, url string) error {
	RedirectHeader(c, url)
	return c.Redirect(http.StatusFound, url)
}

func RedirectHeader(c echo.Context, url string) {
	c.Response().Header().Set("HX-Redirect", url)
}

func Navigate(c echo.Context, url string) error {
	c.Response().Header().Set("HX-Replace-Url", url)
	c.Response().Header().Set("HX-Refresh", "true")
	return c.Redirect(http.StatusFound, url)
}

func Fallback(c echo.Context) error {
	requester := helpers.GetSessionUser(c)
	if requester == nil {
		return Redirect(c, config.LoginPage)
	}
	return Redirect(c, config.Onboarding)
}

func RouteGuard(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		requester := helpers.GetSessionUser(c)
		path := c.Request().URL.Path
		if requester == nil {
			if slices.Contains(PublicRoutes, path) {
				return next(c)
			} else {
				return Navigate(c, config.LoginPage)
			}
		} else {
			if slices.Contains(PublicRoutes, path) {
				return Navigate(c, config.Onboarding)
			}
			return next(c)
		}
	}
}

func JwtGuard() echo.MiddlewareFunc {
	jwtConfig := echojwt.Config{
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(types.JwtClaims)
		},
		SigningKey:  []byte("secret"),
		TokenLookup: "cookie:token",
		Skipper: func(c echo.Context) bool {
			token, err := c.Cookie("token")
			if err != nil || token == nil {
				return slices.Contains(PublicRoutes, c.Request().URL.Path)
			}
			return false
		},
		ErrorHandler: func(c echo.Context, err error) error {
			if slices.Contains(PublicRoutes, c.Request().URL.Path) {
				return nil
			}
			return Navigate(c, config.LoginPage)
		},
	}
	return echojwt.WithConfig(jwtConfig)
}
