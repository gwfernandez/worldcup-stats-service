-- name: ListChampions :many
SELECT
    c.unified_code AS team_code,
    t.name,
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
INNER JOIN teams t ON t.code = c.unified_code
ORDER BY c.wins DESC, t.name ASC
LIMIT $1 OFFSET $2;

-- name: CountChampions :one
SELECT COUNT(1)
FROM (
    SELECT t.unified_code
    FROM championships_stats cs
    INNER JOIN teams t ON t.code = cs.champion_code
    GROUP BY t.unified_code
) c;
