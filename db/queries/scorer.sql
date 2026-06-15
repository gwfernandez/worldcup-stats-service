-- name: ListScorers :many
SELECT
    TRIM(CONCAT_WS(' ', NULLIF(p.first_name, ''), NULLIF(p.last_name, '')))::text AS full_name,
    t.unified_code AS team_code,
    ps.goals,
    p.list_teams,
    t.confederation_code
FROM players_stats ps
INNER JOIN players p ON p.id = ps.id
INNER JOIN teams t ON cardinality(p.list_teams) > 0 AND p.list_teams[1] = t.code
WHERE ps.goals > 0
    AND (
        $1::text = ''
        OR LOWER(p.first_name) LIKE '%' || LOWER($1) || '%'
        OR LOWER(p.last_name) LIKE '%' || LOWER($1) || '%'
    )
    AND ($2::text = '' OR t.unified_code = $2)
    AND ($3::text = '' OR t.confederation_code = $3)
ORDER BY ps.goals DESC, full_name ASC
LIMIT sqlc.arg(limit_value) OFFSET sqlc.arg(offset_value);

-- name: CountScorers :one
SELECT COUNT(*)
FROM players_stats ps
INNER JOIN players p ON p.id = ps.id
INNER JOIN teams t ON cardinality(p.list_teams) > 0 AND p.list_teams[1] = t.code
WHERE ps.goals > 0
    AND (
        $1::text = ''
        OR LOWER(p.first_name) LIKE '%' || LOWER($1) || '%'
        OR LOWER(p.last_name) LIKE '%' || LOWER($1) || '%'
    )
    AND ($2::text = '' OR t.unified_code = $2)
    AND ($3::text = '' OR t.confederation_code = $3);
