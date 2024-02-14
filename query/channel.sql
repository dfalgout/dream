-- name: ListChannels :many
SELECT * FROM
channels
WHERE id IN (
  SELECT channel_id FROM channel_users WHERE user_id = ?
)
ORDER BY name ASC;

-- name: GetChannel :one
SELECT * FROM
channels
WHERE id = (
  SELECT channel_id FROM channel_users WHERE user_id = ? AND channel_id = ?
)
LIMIT 1;

-- name: GetChannelByName :one
SELECT * FROM
channels
WHERE name = ?
AND id IN (
  SELECT channel_id FROM channel_users WHERE user_id = ?
)
LIMIT 1;

-- name: CreateChannel :one
INSERT INTO channels (id, name, description) VALUES (?, ?, ?) RETURNING *;

-- name: AddUserToChannel :exec
INSERT INTO channel_users (channel_id, user_id) VALUES (?, ?);

-- name: RemoveUserFromChannel :exec
DELETE FROM channel_users WHERE channel_id = ? AND user_id = ?;

-- name: GetChannelUsers :many
SELECT * FROM
users
WHERE id IN (
  SELECT user_id FROM channel_users WHERE channel_id = ?
)
ORDER BY created_at DESC;

-- name: GetChannelMessages :many
SELECT * FROM
messages
WHERE channel_id = ?
ORDER BY created_at DESC
LIMIT 100;

-- name: CreateMessage :one
INSERT INTO messages (id, message, channel_id, user_id) VALUES (?, ?, ?, ?) RETURNING *;

-- name: GetChannelUser :one
SELECT * FROM
channel_users
WHERE channel_id = ? AND user_id = ?
LIMIT 1;

-- name: GetChannelMessage :one
SELECT * FROM
messages
WHERE channel_id = ? AND id = ?
LIMIT 1;

-- name: GetChannelUserCount :one
SELECT COUNT(*) FROM
channel_users
WHERE channel_id = ?;

-- name: GetChannelMessageCount :one
SELECT COUNT(*) FROM
messages
WHERE channel_id = ?;
