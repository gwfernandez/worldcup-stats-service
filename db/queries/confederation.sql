-- name: ListConfederations :many
SELECT code, name FROM confederations ORDER BY code;

-- name: GetConfederationByCode :one
SELECT code, name FROM confederations WHERE lower(code) = lower($1);
