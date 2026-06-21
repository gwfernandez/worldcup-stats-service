-- name: ListScorers :many
SELECT
    p.id AS player_id,
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

-- name: GetScorerByPlayerID :one
SELECT
    p.id,
    COALESCE(p.first_name, '')::text AS first_name,
    COALESCE(p.last_name, '')::text AS last_name,
    COALESCE(p.position::text, '') AS position,
    p.list_championships,
    p.list_teams
FROM players p
WHERE p.id = sqlc.arg(player_id);

-- name: ListScorerGoalsByPlayer :many
SELECT
    g.year,
    c.host_codes,
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
INNER JOIN championships c ON c.year = g.year
WHERE g.player_id = sqlc.arg(player_id)
    AND g.own_goal = FALSE
ORDER BY m.match_date ASC NULLS LAST, g.minute_regular ASC, g.id ASC;
