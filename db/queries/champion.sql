-- name: ListChampions :many
SELECT
    c.champion_code AS team_code,
    t.name,
    c.wins,
    c.years
FROM (
    SELECT
        cs.champion_code,
        COUNT(1) AS wins,
        ARRAY_AGG(cs.year ORDER BY cs.year ASC)::integer[] AS years
    FROM championships_stats cs
    GROUP BY cs.champion_code
) c
INNER JOIN teams t ON t.code = c.champion_code
ORDER BY c.wins DESC, t.name ASC
LIMIT $1 OFFSET $2;

-- name: CountChampions :one
SELECT COUNT(1)
FROM (
    SELECT cs.champion_code
    FROM championships_stats cs
    GROUP BY cs.champion_code
) c;
