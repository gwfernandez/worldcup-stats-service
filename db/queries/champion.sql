-- name: ListChampions :many
SELECT
    c.unified_code AS team_code,
    c.wins,
    c.years,
    c.confederation_code
FROM (
    SELECT
        t.unified_code,
        t.confederation_code,
        COUNT(1) AS wins,
        ARRAY_AGG(cs.year ORDER BY cs.year ASC)::integer[] AS years
    FROM championships_stats cs
    INNER JOIN teams t ON t.code = cs.champion_code
    GROUP BY t.unified_code, t.confederation_code
) c
ORDER BY c.wins DESC, c.unified_code ASC
LIMIT sqlc.arg(limit_value) OFFSET sqlc.arg(offset_value);

-- name: CountChampions :one
SELECT COUNT(1)
FROM (
    SELECT t.unified_code, t.confederation_code
    FROM championships_stats cs
    INNER JOIN teams t ON t.code = cs.champion_code
    GROUP BY t.unified_code, t.confederation_code
) c;

-- name: ListFinalsWonByTeam :many
SELECT
    c.year,
    c.host_codes,
    m.match_date,
    m.match_time,
    m.home_team_code,
    m.home_team_score,
    m.home_team_score_penalties,
    m.away_team_code,
    m.away_team_score,
    m.away_team_score_penalties
FROM teams t
INNER JOIN championships c
    ON c.champion_code = t.code
    AND t.unified_code = sqlc.arg(team_code)
INNER JOIN matches m
    ON m.year = c.year
    AND (
        (m.home_team_code = t.code AND m.home_team_win)
        OR (m.away_team_code = t.code AND m.away_team_win)
    )
WHERE m.stage = 'final'
    OR (m.year = 1950 AND m.id = 75)
ORDER BY c.year ASC
LIMIT sqlc.arg(limit_value) OFFSET sqlc.arg(offset_value);

-- name: CountFinalsWonByTeam :one
SELECT COUNT(1)
FROM teams t
INNER JOIN championships c
    ON c.champion_code = t.code
    AND t.unified_code = sqlc.arg(team_code)
INNER JOIN matches m
    ON m.year = c.year
    AND (
        (m.home_team_code = t.code AND m.home_team_win)
        OR (m.away_team_code = t.code AND m.away_team_win)
    )
WHERE m.stage = 'final'
    OR (m.year = 1950 AND m.id = 75);
