-- name: ListTeams :many
SELECT
    t.code,
    COALESCE(tt.name, t.name)::varchar AS name,
    t.dissolution_date,
    t.confederation_code,
    t.federation_name,
    t.federation_code
FROM teams t
LEFT JOIN team_translations tt
    ON tt.team_code = t.code
    AND tt.language = sqlc.arg(language)
WHERE
    (sqlc.arg(name_filter)::text = '' OR LOWER(COALESCE(tt.name, t.name)) LIKE '%' || LOWER(sqlc.arg(name_filter)) || '%')
    AND (sqlc.arg(confederation_code)::text = '' OR LOWER(t.confederation_code) = LOWER(sqlc.arg(confederation_code)))
    AND (sqlc.arg(federation_name)::text = '' OR LOWER(t.federation_name) LIKE '%' || LOWER(sqlc.arg(federation_name)) || '%')
    AND (sqlc.arg(federation_code)::text = '' OR LOWER(t.federation_code) = LOWER(sqlc.arg(federation_code)))
    AND (sqlc.arg(include_dissolved)::boolean OR t.dissolution_date IS NULL)
ORDER BY COALESCE(tt.name, t.name) ASC
LIMIT sqlc.arg(limit_value) OFFSET sqlc.arg(offset_value);

-- name: CountTeams :one
SELECT COUNT(*)
FROM teams t
LEFT JOIN team_translations tt
    ON tt.team_code = t.code
    AND tt.language = sqlc.arg(language)
WHERE
    (sqlc.arg(name_filter)::text = '' OR LOWER(COALESCE(tt.name, t.name)) LIKE '%' || LOWER(sqlc.arg(name_filter)) || '%')
    AND (sqlc.arg(confederation_code)::text = '' OR LOWER(t.confederation_code) = LOWER(sqlc.arg(confederation_code)))
    AND (sqlc.arg(federation_name)::text = '' OR LOWER(t.federation_name) LIKE '%' || LOWER(sqlc.arg(federation_name)) || '%')
    AND (sqlc.arg(federation_code)::text = '' OR LOWER(t.federation_code) = LOWER(sqlc.arg(federation_code)))
    AND (sqlc.arg(include_dissolved)::boolean OR t.dissolution_date IS NULL);

-- name: GetTeamByCode :one
SELECT
    t.code,
    COALESCE(tt.name, t.name)::varchar AS name,
    t.dissolution_date,
    t.confederation_code,
    t.federation_name,
    t.federation_code
FROM teams t
LEFT JOIN team_translations tt
    ON tt.team_code = t.code
    AND tt.language = sqlc.arg(language)
WHERE LOWER(t.code) = LOWER(sqlc.arg(code));
