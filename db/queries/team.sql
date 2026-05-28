-- name: ListTeams :many
SELECT
    code,
    name,
    dissolution_date,
    confederation_code,
    federation_name,
    federation_code
FROM teams
WHERE
    ($1::text = '' OR LOWER(name) LIKE '%' || LOWER($1) || '%')
    AND ($2::text = '' OR LOWER(confederation_code) = LOWER($2))
    AND ($3::text = '' OR LOWER(federation_name) LIKE '%' || LOWER($3) || '%')
    AND ($4::text = '' OR LOWER(federation_code) = LOWER($4))
    AND ($5::boolean OR dissolution_date IS NULL)
ORDER BY name ASC
LIMIT $6 OFFSET $7;

-- name: CountTeams :one
SELECT COUNT(*)
FROM teams
WHERE
    ($1::text = '' OR LOWER(name) LIKE '%' || LOWER($1) || '%')
    AND ($2::text = '' OR LOWER(confederation_code) = LOWER($2))
    AND ($3::text = '' OR LOWER(federation_name) LIKE '%' || LOWER($3) || '%')
    AND ($4::text = '' OR LOWER(federation_code) = LOWER($4))
    AND ($5::boolean OR dissolution_date IS NULL);

-- name: GetTeamByCode :one
SELECT
    code,
    name,
    dissolution_date,
    confederation_code,
    federation_name,
    federation_code
FROM teams
WHERE LOWER(code) = LOWER($1);
