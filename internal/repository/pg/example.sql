-- name: GetExampleByID :one
SELECT * FROM examples WHERE id = $1;

-- name: CreateExample :one
INSERT INTO examples (title, content) VALUES ($1, $2) RETURNING id;
