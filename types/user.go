package types

import (
	"time"

	"github.com/dfalgout/dream/db"
	"github.com/golang-jwt/jwt/v5"
)

type JwtClaims struct {
	Email   string `json:"email"`
	IsAdmin bool   `json:"isAdmin"`
	jwt.RegisteredClaims
}

type SessionUser struct {
	ID      string `json:"id"`
	Email   string `json:"email"`
	IsAdmin bool   `json:"isAdmin"`
}

type UpsertUserIdInput struct {
	ID              string
	Email           *string
	FullName        *string
	IsAdmin         *bool
	VerifyCode      *string
	VerifyExpiresAt *time.Time
}

type UpsertUserEmailInput struct {
	Email           string
	FullName        *string
	IsAdmin         *bool
	VerifyCode      *string
	VerifyExpiresAt *time.Time
}

type VerifyCodeInput struct {
	Email string
	Code  string
}

type User struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	FullName  *string   `json:"fullName,omitempty"`
	IsAdmin   bool      `json:"isAdmin"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func NewUser(dbuser *db.User) *User {
	u := &User{
		ID:        dbuser.ID,
		Email:     dbuser.Email,
		IsAdmin:   dbuser.IsAdmin,
		CreatedAt: dbuser.CreatedAt,
		UpdatedAt: dbuser.UpdatedAt,
	}
	if dbuser.FullName != nil {
		u.FullName = dbuser.FullName
	}
	return u
}
