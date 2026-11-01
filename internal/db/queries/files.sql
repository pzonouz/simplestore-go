-- name: GetFile :one
SELECT
    f.id,
    f.name,
    f.url,
    f.created_at
FROM
    files AS f
WHERE
    f.id = $1
LIMIT 1;

-- name: ListFiles :many
SELECT
    f.id,
    f.name,
    f.url,
    f.created_at
FROM
    files AS f
ORDER BY
    f.created_at
LIMIT $1 OFFSET $2;

-- name: CreateFile :one
INSERT INTO files(name, url)
    VALUES ($1, $2)
RETURNING
    *;

-- name: UpdateFile :one
UPDATE
    files
SET
    name = $2,
    url = $3
WHERE
    id = $1
RETURNING
    *;

-- name: DeleteFile :exec
DELETE FROM files
WHERE id = $1;

