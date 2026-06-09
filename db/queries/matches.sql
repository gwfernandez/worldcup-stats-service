-- name: ListMatchesByYear :many
SELECT
    id,
    year,
    stage::text AS stage,
    stage_type::text AS stage_type,
    group_code,
    replayed,
    replay_of,
    match_date,
    match_time,
    stadium_id,
    home_team_code,
    COALESCE(htt.name, ht.name)::varchar AS home_team_name,
    away_team_code,
    COALESCE(att.name, at.name)::varchar AS away_team_name,
    home_team_score,
    away_team_score,
    extra_time,
    penalty_shootout,
    home_team_score_penalties,
    away_team_score_penalties,
    home_team_win,
    away_team_win,
    draw,
    ref_id
FROM matches
INNER JOIN teams ht ON ht.code = matches.home_team_code
INNER JOIN teams at ON at.code = matches.away_team_code
LEFT JOIN team_translations htt
    ON htt.team_code = ht.code
    AND htt.language = sqlc.arg(language)
LEFT JOIN team_translations att
    ON att.team_code = at.code
    AND att.language = sqlc.arg(language)
WHERE year = $1
ORDER BY matches.stage, group_code, match_date, match_time;
