-- name: GetCategory :one
SELECT
    c.id,
    c.name,
    c.description,
    c.parent,
    c.created_at,
    p.name AS parent_name
FROM
    categories AS c
    LEFT JOIN categories AS p ON p.id = c.parent
WHERE
    c.id = $1
LIMIT 1;

-- name: ListCategories :many
SELECT
    c.id,
    c.name,
    c.description,
    c.parent,
    c.created_at,
    p.name AS parent_name
FROM
    categories AS c
    LEFT JOIN categories AS p ON p.id = c.parent
ORDER BY
    c.created_at
LIMIT $1 OFFSET $2;

-- name: CreateCategory :one
INSERT INTO categories(name, parent, description)
    VALUES ($1, $2, $3)
RETURNING
    *;

-- name: UpdateCategory :one
UPDATE
    categories
SET
    name = $2,
    parent = $3,
    description = $4
WHERE
    id = $1
RETURNING
    *;

-- name: DeleteCategory :exec
DELETE FROM categories
WHERE id = $1;

