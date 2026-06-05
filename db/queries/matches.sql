-- name: ListMatchesByYear :many
SELECT
    id,
    year,
    stage,
    stage_type,
    group_code,
    replayed,
    replay_of,
    match_date,
    match_time,
    stadium_id,
    home_team_code,
    away_team_code,
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
WHERE year = $1
ORDER BY stage, group_code, match_date, match_time;
