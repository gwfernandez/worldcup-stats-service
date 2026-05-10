-- name: ListConfederations :many
SELECT id, code, name FROM confederations ORDER BY id;

-- name: GetConfederation :one
SELECT id, code, name FROM confederations WHERE id = $1;
