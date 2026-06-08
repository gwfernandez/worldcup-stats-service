-- name: ListStandings :many
SELECT
    s.team_code,
    t.name,
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
WHERE
    ($1::text = '' OR LOWER(t.name) LIKE '%' || LOWER($1) || '%')
    AND ($2::text = '' OR t.confederation_code = $2)
ORDER BY s.position ASC
LIMIT $3 OFFSET $4;

-- name: CountStandings :one
SELECT COUNT(*)
FROM standings s
INNER JOIN teams t ON t.code = s.team_code
WHERE
    ($1::text = '' OR LOWER(t.name) LIKE '%' || LOWER($1) || '%')
    AND ($2::text = '' OR t.confederation_code = $2);
