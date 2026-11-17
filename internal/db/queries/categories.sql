-- name: GetCategory :one
SELECT
    c.id,
    c.name,
    c.description,
    c.parent,
    c.slug,
    c.priority,
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
    c.slug,
    c.priority,
    c.created_at,
    p.name AS parent_name
FROM
    categories AS c
    LEFT JOIN categories AS p ON p.id = c.parent
ORDER BY
    c.priority,c.created_at
LIMIT $1 OFFSET $2;


-- name: ListSubCategoriesBySlug :many
SELECT
    c.id,
    c.name,
    c.description,
    c.parent,
    c.slug,
    c.priority,
    c.created_at,
    p.name AS parent_name
FROM
    categories AS c
    LEFT JOIN categories AS p ON p.id = c.parent
    WHERE p.slug=$1
ORDER BY
    c.priority,c.created_at;

-- name: ListCategoriesTree :one
SELECT tree FROM category_tree_view;

-- name: NoChildrenCategories :many
SELECT
    id,
    name,
    description,
    parent,
    slug,
    priority,
    created_at
FROM
    categories c
WHERE
    NOT EXISTS (
        SELECT
            1
        FROM
            categories child
        WHERE
            child.parent = c.id)
ORDER BY
    priority,created_at;

-- name: CreateCategory :one
INSERT INTO categories(name, parent, description,slug, priority)
    VALUES ($1, $2, $3, $4, $5)
RETURNING
    *;

-- name: UpdateCategory :one
UPDATE
    categories
SET
    name = $2,
    parent = $3,
    description = $4,
    slug=$5,
    priority=$6
WHERE
    id = $1
RETURNING
    *;

-- name: DeleteCategory :exec
DELETE FROM categories
WHERE id = $1;

