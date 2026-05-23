-- name: ListNationalTeams :many
SELECT
    id,
    name,
    code,
    dissolution_date,
    confederation_id,
    federation_name,
    federation_code
FROM national_teams
WHERE
    ($1::text = '' OR LOWER(name) LIKE '%' || LOWER($1) || '%')
    AND ($2::bigint = 0 OR confederation_id = $2)
    AND ($3::text = '' OR LOWER(federation_name) LIKE '%' || LOWER($3) || '%')
    AND ($4::text = '' OR LOWER(federation_code) = LOWER($4))
    AND ($5::boolean OR dissolution_date IS NULL)
ORDER BY name ASC
LIMIT $6 OFFSET $7;

-- name: CountNationalTeams :one
SELECT COUNT(*)
FROM national_teams
WHERE
    ($1::text = '' OR LOWER(name) LIKE '%' || LOWER($1) || '%')
    AND ($2::bigint = 0 OR confederation_id = $2)
    AND ($3::text = '' OR LOWER(federation_name) LIKE '%' || LOWER($3) || '%')
    AND ($4::text = '' OR LOWER(federation_code) = LOWER($4))
    AND ($5::boolean OR dissolution_date IS NULL);

-- name: GetNationalTeamByID :one
SELECT
    id,
    name,
    code,
    dissolution_date,
    confederation_id,
    federation_name,
    federation_code
FROM national_teams
WHERE id = $1;

-- name: GetNationalTeamByCode :one
SELECT
    id,
    name,
    code,
    dissolution_date,
    confederation_id,
    federation_name,
    federation_code
FROM national_teams
WHERE LOWER(code) = LOWER($1);
