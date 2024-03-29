// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: channel.sql

package db

import (
	"context"
)

const addUserToChannel = `-- name: AddUserToChannel :exec
INSERT INTO channel_users (channel_id, user_id) VALUES (?, ?)
`

type AddUserToChannelParams struct {
	ChannelID string
	UserID    string
}

func (q *Queries) AddUserToChannel(ctx context.Context, arg *AddUserToChannelParams) error {
	_, err := q.db.ExecContext(ctx, addUserToChannel, arg.ChannelID, arg.UserID)
	return err
}

const createChannel = `-- name: CreateChannel :one
INSERT INTO channels (id, name, description) VALUES (?, ?, ?) RETURNING id, name, description, created_at, updated_at
`

type CreateChannelParams struct {
	ID          string
	Name        string
	Description string
}

func (q *Queries) CreateChannel(ctx context.Context, arg *CreateChannelParams) (*Channel, error) {
	row := q.db.QueryRowContext(ctx, createChannel, arg.ID, arg.Name, arg.Description)
	var i Channel
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Description,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return &i, err
}

const createMessage = `-- name: CreateMessage :one
INSERT INTO messages (id, message, channel_id, user_id) VALUES (?, ?, ?, ?) RETURNING id, message, user_id, channel_id, created_at
`

type CreateMessageParams struct {
	ID        string
	Message   string
	ChannelID string
	UserID    string
}

func (q *Queries) CreateMessage(ctx context.Context, arg *CreateMessageParams) (*Message, error) {
	row := q.db.QueryRowContext(ctx, createMessage,
		arg.ID,
		arg.Message,
		arg.ChannelID,
		arg.UserID,
	)
	var i Message
	err := row.Scan(
		&i.ID,
		&i.Message,
		&i.UserID,
		&i.ChannelID,
		&i.CreatedAt,
	)
	return &i, err
}

const getChannel = `-- name: GetChannel :one
SELECT id, name, description, created_at, updated_at FROM
channels
WHERE id = (
  SELECT channel_id FROM channel_users WHERE user_id = ? AND channel_id = ?
)
LIMIT 1
`

type GetChannelParams struct {
	UserID    string
	ChannelID string
}

func (q *Queries) GetChannel(ctx context.Context, arg *GetChannelParams) (*Channel, error) {
	row := q.db.QueryRowContext(ctx, getChannel, arg.UserID, arg.ChannelID)
	var i Channel
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Description,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return &i, err
}

const getChannelByName = `-- name: GetChannelByName :one
SELECT id, name, description, created_at, updated_at FROM
channels
WHERE name = ?
AND id IN (
  SELECT channel_id FROM channel_users WHERE user_id = ?
)
LIMIT 1
`

type GetChannelByNameParams struct {
	Name   string
	UserID string
}

func (q *Queries) GetChannelByName(ctx context.Context, arg *GetChannelByNameParams) (*Channel, error) {
	row := q.db.QueryRowContext(ctx, getChannelByName, arg.Name, arg.UserID)
	var i Channel
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Description,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return &i, err
}

const getChannelMessage = `-- name: GetChannelMessage :one
SELECT id, message, user_id, channel_id, created_at FROM
messages
WHERE channel_id = ? AND id = ?
LIMIT 1
`

type GetChannelMessageParams struct {
	ChannelID string
	ID        string
}

func (q *Queries) GetChannelMessage(ctx context.Context, arg *GetChannelMessageParams) (*Message, error) {
	row := q.db.QueryRowContext(ctx, getChannelMessage, arg.ChannelID, arg.ID)
	var i Message
	err := row.Scan(
		&i.ID,
		&i.Message,
		&i.UserID,
		&i.ChannelID,
		&i.CreatedAt,
	)
	return &i, err
}

const getChannelMessageCount = `-- name: GetChannelMessageCount :one
SELECT COUNT(*) FROM
messages
WHERE channel_id = ?
`

func (q *Queries) GetChannelMessageCount(ctx context.Context, channelID string) (int64, error) {
	row := q.db.QueryRowContext(ctx, getChannelMessageCount, channelID)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const getChannelMessages = `-- name: GetChannelMessages :many
SELECT id, message, user_id, channel_id, created_at FROM
messages
WHERE channel_id = ?
ORDER BY created_at DESC
LIMIT 100
`

func (q *Queries) GetChannelMessages(ctx context.Context, channelID string) ([]*Message, error) {
	rows, err := q.db.QueryContext(ctx, getChannelMessages, channelID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []*Message
	for rows.Next() {
		var i Message
		if err := rows.Scan(
			&i.ID,
			&i.Message,
			&i.UserID,
			&i.ChannelID,
			&i.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, &i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getChannelUser = `-- name: GetChannelUser :one
SELECT user_id, channel_id, created_at FROM
channel_users
WHERE channel_id = ? AND user_id = ?
LIMIT 1
`

type GetChannelUserParams struct {
	ChannelID string
	UserID    string
}

func (q *Queries) GetChannelUser(ctx context.Context, arg *GetChannelUserParams) (*ChannelUser, error) {
	row := q.db.QueryRowContext(ctx, getChannelUser, arg.ChannelID, arg.UserID)
	var i ChannelUser
	err := row.Scan(&i.UserID, &i.ChannelID, &i.CreatedAt)
	return &i, err
}

const getChannelUserCount = `-- name: GetChannelUserCount :one
SELECT COUNT(*) FROM
channel_users
WHERE channel_id = ?
`

func (q *Queries) GetChannelUserCount(ctx context.Context, channelID string) (int64, error) {
	row := q.db.QueryRowContext(ctx, getChannelUserCount, channelID)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const getChannelUsers = `-- name: GetChannelUsers :many
SELECT id, email, full_name, is_admin, verify_code, verify_expires_at, created_at, updated_at FROM
users
WHERE id IN (
  SELECT user_id FROM channel_users WHERE channel_id = ?
)
ORDER BY created_at DESC
`

func (q *Queries) GetChannelUsers(ctx context.Context, channelID string) ([]*User, error) {
	rows, err := q.db.QueryContext(ctx, getChannelUsers, channelID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []*User
	for rows.Next() {
		var i User
		if err := rows.Scan(
			&i.ID,
			&i.Email,
			&i.FullName,
			&i.IsAdmin,
			&i.VerifyCode,
			&i.VerifyExpiresAt,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, &i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listChannels = `-- name: ListChannels :many
SELECT id, name, description, created_at, updated_at FROM
channels
WHERE id IN (
  SELECT channel_id FROM channel_users WHERE user_id = ?
)
ORDER BY name ASC
`

func (q *Queries) ListChannels(ctx context.Context, userID string) ([]*Channel, error) {
	rows, err := q.db.QueryContext(ctx, listChannels, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []*Channel
	for rows.Next() {
		var i Channel
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Description,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, &i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const removeUserFromChannel = `-- name: RemoveUserFromChannel :exec
DELETE FROM channel_users WHERE channel_id = ? AND user_id = ?
`

type RemoveUserFromChannelParams struct {
	ChannelID string
	UserID    string
}

func (q *Queries) RemoveUserFromChannel(ctx context.Context, arg *RemoveUserFromChannelParams) error {
	_, err := q.db.ExecContext(ctx, removeUserFromChannel, arg.ChannelID, arg.UserID)
	return err
}
