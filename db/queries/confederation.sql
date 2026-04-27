-- name: ListConfederations :many
SELECT id, code, name FROM confederations ORDER BY id;

-- name: GetConfederation :one
SELECT id, code, name FROM confederations WHERE id = $1;

-- name: CreateConfederation :one
INSERT INTO confederations (code, name) VALUES ($1, $2) RETURNING id, code, name;

-- name: UpdateConfederation :one
UPDATE confederations SET code = $2, name = $3 WHERE id = $1 RETURNING id, code, name;

-- name: DeleteConfederation :exec
DELETE FROM confederations WHERE id = $1;
