package helpers

import (
	"time"

	"github.com/dfalgout/dream/db"
	"github.com/dfalgout/dream/types"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

func CreateToken(user *db.User) (*string, error) {
	claims := &types.JwtClaims{
		Email:   user.Email,
		IsAdmin: user.IsAdmin,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "tasai",
			Subject:   user.ID,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 72)),
		},
	}

	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString([]byte("secret"))
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func GetSessionUser(c echo.Context) *types.SessionUser {
	user, ok := c.Get("user").(*jwt.Token)
	if !ok {
		return nil
	}
	claims, ok := user.Claims.(*types.JwtClaims)
	if !ok {
		return nil
	}
	return &types.SessionUser{
		ID:      claims.Subject,
		Email:   claims.Email,
		IsAdmin: claims.IsAdmin,
	}
}
