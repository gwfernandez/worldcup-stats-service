-- name: ListChampionships :many
SELECT
    year,
    start_date,
    end_date,
    host_codes,
    confederation_codes,
    champion_code
FROM championships c
WHERE
    ($1::integer = 0 OR c.year = $1)
    AND ($2::text = '' OR EXISTS (
        SELECT 1
        FROM teams t
        LEFT JOIN team_translations tt
            ON tt.team_code = t.code
            AND tt.language = sqlc.arg(language)
        WHERE t.code = ANY(c.host_codes)
          AND LOWER(COALESCE(tt.name, t.name)) LIKE '%' || LOWER($2) || '%'
    ))
    AND ($3::text = '' OR EXISTS (
        SELECT 1
        FROM teams t
        WHERE t.code = ANY(c.host_codes)
          AND LOWER(t.confederation_code) = LOWER($3)
    ))
ORDER BY c.year ASC
LIMIT $4 OFFSET $5;

-- name: ListChampionshipsWithoutHostFilter :many
SELECT
    year,
    start_date,
    end_date,
    host_codes,
    confederation_codes,
    champion_code
FROM championships c
WHERE
    ($1::integer = 0 OR c.year = $1)
    AND ($2::text = '' OR EXISTS (
        SELECT 1
        FROM teams t
        WHERE t.code = ANY(c.host_codes)
          AND LOWER(t.confederation_code) = LOWER($2)
    ))
ORDER BY c.year ASC
LIMIT $3 OFFSET $4;

-- name: CountChampionships :one
SELECT COUNT(*)
FROM championships c
WHERE
    ($1::integer = 0 OR c.year = $1)
    AND ($2::text = '' OR EXISTS (
        SELECT 1
        FROM teams t
        LEFT JOIN team_translations tt
            ON tt.team_code = t.code
            AND tt.language = sqlc.arg(language)
        WHERE t.code = ANY(c.host_codes)
          AND LOWER(COALESCE(tt.name, t.name)) LIKE '%' || LOWER($2) || '%'
    ))
    AND ($3::text = '' OR EXISTS (
        SELECT 1
        FROM teams t
        WHERE t.code = ANY(c.host_codes)
          AND LOWER(t.confederation_code) = LOWER($3)
    ));

-- name: CountChampionshipsWithoutHostFilter :one
SELECT COUNT(*)
FROM championships c
WHERE
    ($1::integer = 0 OR c.year = $1)
    AND ($2::text = '' OR EXISTS (
        SELECT 1
        FROM teams t
        WHERE t.code = ANY(c.host_codes)
          AND LOWER(t.confederation_code) = LOWER($2)
    ));

-- name: GetChampionshipByYear :one
SELECT 
    c.year,
    c.start_date,
    c.end_date,
    c.host_codes,
    c.confederation_codes,
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
LEFT JOIN championships_stats s ON s.year = c.year
WHERE c.year = $1;

-- name: ListChampionshipTeamsByYear :many
SELECT
    ct.year,
    ct.team_code,
    t.confederation_code,
    cts.group_code,
    COALESCE(CASE
        WHEN ct.team_code = cs.champion_code THEN 'champion'
        WHEN ct.team_code = cs.runner_up_code THEN 'runner_up'
        WHEN ct.team_code = cs.third_place_code THEN 'third_place'
        WHEN ct.team_code = cs.fourth_place_code THEN 'fourth_place'
        ELSE cts.stage_reached::text
    END, '')::text AS stage_reached,
    COALESCE(m.managers, '')::text AS managers
FROM championships_teams ct
INNER JOIN teams t ON t.code = ct.team_code
LEFT JOIN team_translations tt
    ON tt.team_code = t.code
    AND tt.language = sqlc.arg(language)
INNER JOIN championships_teams_stats cts ON ct.year = cts.year AND ct.team_code = cts.team_code
INNER JOIN championships_stats cs ON cs.year = ct.year
LEFT JOIN (
    SELECT
        cm.team_code,
        string_agg(NULLIF(TRIM(CONCAT_WS(' ', NULLIF(m.first_name, ''), NULLIF(m.last_name, ''))), ''), ', ') AS managers
    FROM championships_managers cm
    INNER JOIN managers m ON cm.manager_id = m.id
    WHERE cm.year = $1
    GROUP BY cm.team_code
) m ON m.team_code = ct.team_code
WHERE ct.year = $1
    AND ($2::text = '' OR LOWER(COALESCE(tt.name, t.name)) LIKE '%' || LOWER($2) || '%')
    AND ($3::text = '' OR t.confederation_code = $3)
    AND ($4::text = '' OR cts.group_code = $4)
ORDER BY cts.position ASC, cts.stage_reached DESC
LIMIT $5 OFFSET $6;

-- name: ListChampionshipTeamsByYearWithoutNameFilter :many
SELECT
    ct.year,
    ct.team_code,
    t.confederation_code,
    cts.group_code,
    COALESCE(CASE
        WHEN ct.team_code = cs.champion_code THEN 'champion'
        WHEN ct.team_code = cs.runner_up_code THEN 'runner_up'
        WHEN ct.team_code = cs.third_place_code THEN 'third_place'
        WHEN ct.team_code = cs.fourth_place_code THEN 'fourth_place'
        ELSE cts.stage_reached::text
    END, '')::text AS stage_reached,
    COALESCE(m.managers, '')::text AS managers
FROM championships_teams ct
INNER JOIN teams t ON t.code = ct.team_code
INNER JOIN championships_teams_stats cts ON ct.year = cts.year AND ct.team_code = cts.team_code
INNER JOIN championships_stats cs ON cs.year = ct.year
LEFT JOIN (
    SELECT
        cm.team_code,
        string_agg(NULLIF(TRIM(CONCAT_WS(' ', NULLIF(m.first_name, ''), NULLIF(m.last_name, ''))), ''), ', ') AS managers
    FROM championships_managers cm
    INNER JOIN managers m ON cm.manager_id = m.id
    WHERE cm.year = $1
    GROUP BY cm.team_code
) m ON m.team_code = ct.team_code
WHERE ct.year = $1
    AND ($2::text = '' OR t.confederation_code = $2)
    AND ($3::text = '' OR cts.group_code = $3)
ORDER BY cts.position ASC, cts.stage_reached DESC
LIMIT $4 OFFSET $5;

-- name: CountChampionshipTeamsByYear :one
SELECT COUNT(*)
FROM championships_teams ct
INNER JOIN teams t ON t.code = ct.team_code
LEFT JOIN team_translations tt
    ON tt.team_code = t.code
    AND tt.language = sqlc.arg(language)
INNER JOIN championships_teams_stats cts ON ct.year = cts.year AND ct.team_code = cts.team_code
WHERE ct.year = $1
    AND ($2::text = '' OR LOWER(COALESCE(tt.name, t.name)) LIKE '%' || LOWER($2) || '%')
    AND ($3::text = '' OR t.confederation_code = $3)
    AND ($4::text = '' OR cts.group_code = $4);

-- name: CountChampionshipTeamsByYearWithoutNameFilter :one
SELECT COUNT(*)
FROM championships_teams ct
INNER JOIN teams t ON t.code = ct.team_code
INNER JOIN championships_teams_stats cts ON ct.year = cts.year AND ct.team_code = cts.team_code
WHERE ct.year = $1
    AND ($2::text = '' OR t.confederation_code = $2)
    AND ($3::text = '' OR cts.group_code = $3);

-- name: ListChampionshipStandingsByYear :many
SELECT
    cts.team_code,
    COALESCE(cts.group_code, '')::text AS group_code,
    cts.matches_played,
    cts.wins,
    cts.draws,
    cts.losses,
    cts.goals_for,
    cts.goals_against,
    cts.goal_difference,
    cts.points,
    cts.unified_points,
    cts.position,
    COALESCE(CASE
        WHEN ct.team_code = cs.champion_code THEN 'champion'
        WHEN ct.team_code = cs.runner_up_code THEN 'runner_up'
        WHEN ct.team_code = cs.third_place_code THEN 'third_place'
        WHEN ct.team_code = cs.fourth_place_code THEN 'fourth_place'
        ELSE cts.stage_reached::text
    END, '')::text AS performance
FROM championships_teams ct
INNER JOIN championships_teams_stats cts ON ct.year = cts.year AND ct.team_code = cts.team_code
INNER JOIN championships_stats cs ON cs.year = ct.year
WHERE ct.year = $1
ORDER BY cts.position ASC, cts.stage_reached
LIMIT $2 OFFSET $3;

-- name: CountChampionshipStandingsByYear :one
SELECT COUNT(*)
FROM championships_teams ct
INNER JOIN championships_teams_stats cts ON ct.year = cts.year AND ct.team_code = cts.team_code
INNER JOIN championships_stats cs ON cs.year = ct.year
WHERE ct.year = $1;

-- name: ListChampionshipStadiumsByYear :many
SELECT
    css.year,
    s.id,
    s.name,
    COALESCE(s.city_name, '')::text AS city_name,
    COALESCE(s.capacity, 0)::integer AS capacity,
    css.matches_played
FROM championships_stadiums_stats css
INNER JOIN stadiums s ON s.id = css.stadium_id
WHERE css.year = $1
    AND ($2::text = '' OR LOWER(s.name) LIKE '%' || LOWER($2) || '%')
ORDER BY css.matches_played DESC, s.name ASC
LIMIT $3 OFFSET $4;

-- name: CountChampionshipStadiumsByYear :one
SELECT COUNT(*)
FROM championships_stadiums_stats css
INNER JOIN stadiums s ON s.id = css.stadium_id
WHERE css.year = $1
    AND ($2::text = '' OR LOWER(s.name) LIKE '%' || LOWER($2) || '%');

-- name: ListChampionshipScorersByYear :many
SELECT
    p.id AS player_id,
    TRIM(CONCAT_WS(' ', NULLIF(p.first_name, ''), NULLIF(p.last_name, '')))::text AS full_name,
    ss.team_code,
    ss.goals
FROM squads_stats ss
INNER JOIN players p ON p.id = ss.player_id
WHERE ss.year = $1
    AND ss.goals > 0
    AND (
        $2::text = ''
        OR LOWER(p.first_name) LIKE '%' || LOWER($2) || '%'
        OR LOWER(p.last_name) LIKE '%' || LOWER($2) || '%'
    )
    AND ($3::text = '' OR ss.team_code = $3)
ORDER BY ss.goals DESC, full_name ASC
LIMIT sqlc.arg(limit_value) OFFSET sqlc.arg(offset_value);

-- name: CountChampionshipScorersByYear :one
SELECT COUNT(*)
FROM squads_stats ss
INNER JOIN players p ON p.id = ss.player_id
WHERE ss.year = $1
    AND ss.goals > 0
    AND (
        $2::text = ''
        OR LOWER(p.first_name) LIKE '%' || LOWER($2) || '%'
        OR LOWER(p.last_name) LIKE '%' || LOWER($2) || '%'
    )
    AND ($3::text = '' OR ss.team_code = $3);
