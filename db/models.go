// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0

package db

import (
	"time"
)

type Channel struct {
	ID          string
	Name        string
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type ChannelUser struct {
	UserID    string
	ChannelID string
	CreatedAt time.Time
}

type Message struct {
	ID        string
	Message   string
	UserID    string
	ChannelID string
	CreatedAt time.Time
}

type User struct {
	ID              string
	Email           string
	FullName        *string
	IsAdmin         bool
	VerifyCode      *string
	VerifyExpiresAt *time.Time
	CreatedAt       time.Time
	UpdatedAt       time.Time
}
