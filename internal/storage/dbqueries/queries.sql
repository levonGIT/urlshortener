-- name: GetUrls :many
SELECT *
FROM urls;

-- name: GetUrlByAlias :one
SELECT *
FROM urls
WHERE alias = $1;

-- name: InsertUrl :one
INSERT INTO urls (url, alias)
VALUES ($1, $2)
RETURNING *;

-- name: UpdateUrl :execrows
UPDATE urls SET alias = $2
WHERE id = $1;

-- name: DeleteUrl :execrows
DELETE from urls
WHERE id = $1;

-- name: GetUrlById :one
SELECT *
FROM urls
WHERE id = $1;