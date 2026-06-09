-- name: ListStandings :many
SELECT
    s.team_code,
    COALESCE(tt.name, t.name)::varchar AS name,
    s.matches_played,
    s.wins,
    s.draws,
    s.losses,
    s.goals_for,
    s.goals_against,
    s.goal_difference,
    s.points,
    s.unified_points,
    s.position,
    s.unified_position
FROM standings s
INNER JOIN teams t ON t.code = s.team_code
LEFT JOIN team_translations tt
    ON tt.team_code = t.code
    AND tt.language = sqlc.arg(language)
WHERE
    (sqlc.arg(name_filter)::text = '' OR LOWER(COALESCE(tt.name, t.name)) LIKE '%' || LOWER(sqlc.arg(name_filter)) || '%')
    AND (sqlc.arg(confederation_code)::text = '' OR t.confederation_code = sqlc.arg(confederation_code))
ORDER BY s.position ASC
LIMIT sqlc.arg(limit_value) OFFSET sqlc.arg(offset_value);

-- name: CountStandings :one
SELECT COUNT(*)
FROM standings s
INNER JOIN teams t ON t.code = s.team_code
LEFT JOIN team_translations tt
    ON tt.team_code = t.code
    AND tt.language = sqlc.arg(language)
WHERE
    (sqlc.arg(name_filter)::text = '' OR LOWER(COALESCE(tt.name, t.name)) LIKE '%' || LOWER(sqlc.arg(name_filter)) || '%')
    AND (sqlc.arg(confederation_code)::text = '' OR t.confederation_code = sqlc.arg(confederation_code));
