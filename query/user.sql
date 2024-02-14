-- name: CreateUser :one
INSERT INTO users (
  id, email, full_name, is_admin, verify_code, verify_expires_at
) VALUES (
  @id,
  @email,
  sqlc.narg('full_name'),
  coalesce(sqlc.narg('is_admin'), 0),
  sqlc.narg('verify_code'),
  sqlc.narg('verify_expires_at')
) RETURNING *;

-- name: UpdateUserById :one
UPDATE users
SET
    email = coalesce(sqlc.narg('email'), email),
    full_name = coalesce(sqlc.narg('full_name'), full_name),
    is_admin = coalesce(sqlc.narg('is_admin'), is_admin),
    verify_code = coalesce(sqlc.narg('verify_code'), verify_code),
    verify_expires_at = coalesce(sqlc.narg('verify_expires_at'), verify_expires_at),
    updated_at = (strftime('%F %R:%f'))
WHERE
    id = @id
RETURNING *;

-- name: UpdateUserByEmail :one
UPDATE users
SET
    full_name = coalesce(sqlc.narg('full_name'), full_name),
    is_admin = coalesce(sqlc.narg('is_admin'), is_admin),
    verify_code = coalesce(sqlc.narg('verify_code'), verify_code),
    verify_expires_at = coalesce(sqlc.narg('verify_expires_at'), verify_expires_at),
    updated_at = (strftime('%F %R:%f'))
WHERE
    email = @email
RETURNING *;

-- name: DeleteUser :execrows
DELETE FROM users
WHERE id = @id;

-- name: ListUsers :many
SELECT * FROM users
ORDER BY created_at DESC;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE
    email = @email
LIMIT 1;

-- name: ClearAndGetVerifiedUser :one
UPDATE users
SET
    verify_code = NULL,
    verify_expires_at = NULL
WHERE
    email = @email AND
    verify_code = @verify_code AND
    verify_expires_at > (strftime('%F %R:%f'))
RETURNING *;
