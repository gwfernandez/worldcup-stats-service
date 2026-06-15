-- name: ListChampions :many
SELECT
    c.unified_code AS team_code,
    c.wins,
    c.years
FROM (
    SELECT
        t.unified_code,
        COUNT(1) AS wins,
        ARRAY_AGG(cs.year ORDER BY cs.year ASC)::integer[] AS years
    FROM championships_stats cs
    INNER JOIN teams t ON t.code = cs.champion_code
    GROUP BY t.unified_code
) c
ORDER BY c.wins DESC, c.unified_code ASC
LIMIT sqlc.arg(limit_value) OFFSET sqlc.arg(offset_value);

-- name: CountChampions :one
SELECT COUNT(1)
FROM (
    SELECT t.unified_code
    FROM championships_stats cs
    INNER JOIN teams t ON t.code = cs.champion_code
    GROUP BY t.unified_code
) c;
