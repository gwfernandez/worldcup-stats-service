-- name: ListChampionships :many
SELECT
    year,
    start_date,
    end_date,
    host_nation_codes,
    champion_code
FROM championships c
WHERE
    ($1::integer = 0 OR c.year = $1)
    AND ($2::text = '' OR EXISTS (
        SELECT 1
        FROM national_teams t
        WHERE t.code = ANY(c.host_nation_codes)
          AND LOWER(t.name) LIKE '%' || LOWER($2) || '%'
    ))
    AND ($3::text = '' OR EXISTS (
        SELECT 1
        FROM national_teams t
        WHERE t.code = ANY(c.host_nation_codes)
          AND LOWER(t.confederation_code) = LOWER($3)
    ))
ORDER BY c.year ASC
LIMIT $4 OFFSET $5;

-- name: CountChampionships :one
SELECT COUNT(*)
FROM championships c
WHERE
    ($1::integer = 0 OR c.year = $1)
    AND ($2::text = '' OR EXISTS (
        SELECT 1
        FROM national_teams t
        WHERE t.code = ANY(c.host_nation_codes)
          AND LOWER(t.name) LIKE '%' || LOWER($2) || '%'
    ))
    AND ($3::text = '' OR EXISTS (
        SELECT 1
        FROM national_teams t
        WHERE t.code = ANY(c.host_nation_codes)
          AND LOWER(t.confederation_code) = LOWER($3)
    ));

-- name: GetChampionshipByYear :one
SELECT 
    c.year,
    c.start_date,
    c.end_date,
    c.host_nation_codes,
    c.champion_code,
    s.total_teams,
    s.total_matches,
    s.total_stadiums,
    s.total_players,
    s.total_goals,
    s.champion_code AS stats_champion_code,
    s.runner_up_code AS stats_runner_up_code,
    s.third_place_code AS stats_third_place_code,
    s.fourth_place_code AS stats_fourth_place_code,
    s.top_scorer_ids,
    s.top_scorer_goals
FROM championships c
LEFT JOIN championship_stats s ON s.year = c.year
WHERE c.year = $1;
