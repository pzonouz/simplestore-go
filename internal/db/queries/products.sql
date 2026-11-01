-- name: GetProduct :one
SELECT
    p.id,
    p.name,
    p.description,
    p.category,
    p.slug,
    p.image_url,
    p.count,
    p.price,
    p.packing_type,
    p.packing_coefficient_types,
    p.created_at,
    c.name AS category_name
FROM
    products AS p
    LEFT JOIN categories AS c ON p.category = c.id
WHERE
    p.id = $1
LIMIT 1;

-- name: GetProductBySlug :one
SELECT
    p.id,
    p.name,
    p.description,
    p.category,
    p.slug,
    p.image_url,
    p.count,
    p.price,
    p.packing_type,
    p.packing_coefficient_types,
    p.created_at,
    c.name AS category_name
FROM
    products AS p
    LEFT JOIN categories AS c ON p.category = c.id
WHERE
    p.slug = $1
LIMIT 1;

-- name: ListProducts :many
SELECT
    p.id,
    p.name,
    p.description,
    p.category,
    p.slug,
    p.image_url,
    p.count,
    p.price,
    p.packing_type,
    p.packing_coefficient_types,
    p.created_at,
    c.name AS category_name
FROM
    products AS p
    LEFT JOIN categories AS c ON p.category = c.id
ORDER BY
    c.name
LIMIT $1 OFFSET $2;

-- name: CreateProduct :one
INSERT INTO products(name, description, category, image_url, price, count, discount, slug, packing_type, packing_coefficient_types)
    VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
RETURNING
    *;

-- name: UpdateProduct :one
UPDATE
    products
SET
    name = $2,
    description = $3,
    category = $4,
    image_url = $5,
    price = $6,
    count = $7,
    discount = $8,
    slug = $9,
    packing_type = $10,
    packing_coefficient_types = $11
WHERE
    id = $1
RETURNING
    *;

-- name: DeleteProduct :exec
DELETE FROM products
WHERE id = $1;

