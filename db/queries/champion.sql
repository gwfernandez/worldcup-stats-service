-- name: ListChampions :many
SELECT
    c.unified_code AS team_code,
    COALESCE(tt.name, t.name)::varchar AS name,
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
LEFT JOIN team_translations tt
    ON tt.team_code = t.code
    AND tt.language = sqlc.arg(language)
ORDER BY c.wins DESC, COALESCE(tt.name, t.name) ASC
LIMIT sqlc.arg(limit_value) OFFSET sqlc.arg(offset_value);

-- name: CountChampions :one
SELECT COUNT(1)
FROM (
    SELECT t.unified_code
    FROM championships_stats cs
    INNER JOIN teams t ON t.code = cs.champion_code
    GROUP BY t.unified_code
) c;
