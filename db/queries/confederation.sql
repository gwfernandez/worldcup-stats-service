-- name: ListConfederations :many
SELECT
    c.code,
    COALESCE(ct.name, c.name)::varchar AS name
FROM confederations c
LEFT JOIN confederation_translations ct
    ON ct.confederation_code = c.code
    AND ct.language = sqlc.arg(language)
ORDER BY c.code;

-- name: GetConfederationByCode :one
SELECT
    c.code,
    COALESCE(ct.name, c.name)::varchar AS name
FROM confederations c
LEFT JOIN confederation_translations ct
    ON ct.confederation_code = c.code
    AND ct.language = sqlc.arg(language)
WHERE lower(c.code) = lower(sqlc.arg(code));
