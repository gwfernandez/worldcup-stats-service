-- name: ListGoalsByPlayer :many
SELECT
    g.year,
    m.match_date,
    (CASE
        WHEN g.team_condition = 'home' THEN m.away_team_code
        ELSE m.home_team_code
    END)::text AS opponent_team_code,
    g.minute_regular,
    g.penalty,
    COALESCE(m.stage::text, '') AS stage
FROM goals g
INNER JOIN matches m ON g.match_id = m.id
WHERE g.player_id = sqlc.arg(player_id)
    AND g.own_goal = FALSE
    AND (sqlc.arg(year)::int = 0 OR g.year = sqlc.arg(year))
ORDER BY m.match_date ASC NULLS LAST, g.minute_regular ASC, g.id ASC
LIMIT sqlc.arg(limit_value) OFFSET sqlc.arg(offset_value);

-- name: CountGoalsByPlayer :one
SELECT COUNT(*)
FROM goals g
INNER JOIN matches m ON g.match_id = m.id
WHERE g.player_id = sqlc.arg(player_id)
    AND g.own_goal = FALSE
    AND (sqlc.arg(year)::int = 0 OR g.year = sqlc.arg(year));
