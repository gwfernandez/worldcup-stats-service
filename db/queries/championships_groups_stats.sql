-- name: ListGroupsStatsByYear :many
SELECT
    year,
    stage,
    group_code,
    team_code,
    matches_played,
    wins,
    draws,
    losses,
    goals_for,
    goals_against,
    goal_difference,
    points,
    unified_points,
    position
FROM championships_groups_stats
WHERE year = $1
ORDER BY stage, group_code, position;
