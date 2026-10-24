-- name: GetUser :one
SELECT
    *
FROM
    users
WHERE
    id = $1
LIMIT 1;

-- name: GetUserByEmail :one
SELECT
    *
FROM
    users
WHERE
    email = $1;

-- name: ListUsers :many
SELECT
    *
FROM
    users
ORDER BY
    email
LIMIT $1 OFFSET $2;

-- name: CreateUser :one
INSERT INTO users(email, phone_number, token, expired_time, password)
    VALUES ($1, $2, $3, $4, $5)
RETURNING
    *;

-- name: UpdateUser :one
UPDATE
    users
SET
    email = $2,
    phone_number = $3,
    password = $4
WHERE
    id = $1
RETURNING
    *;

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1;

-- name: ActivateUserByToken :execrows
UPDATE
    users
SET
    is_active = TRUE
WHERE
    token = $1
    AND expired_time > NOW();

