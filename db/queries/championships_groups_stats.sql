-- name: ListGroupsStatsByYear :many
SELECT
    championships_groups_stats.year,
    championships_groups_stats.stage::text AS stage,
    championships_groups_stats.group_code,
    championships_groups_stats.team_code,
    championships_groups_stats.matches_played,
    championships_groups_stats.wins,
    championships_groups_stats.draws,
    championships_groups_stats.losses,
    championships_groups_stats.goals_for,
    championships_groups_stats.goals_against,
    championships_groups_stats.goal_difference,
    championships_groups_stats.points,
    championships_groups_stats.unified_points,
    championships_groups_stats.position
FROM championships_groups_stats
WHERE year = $1
ORDER BY championships_groups_stats.stage, group_code, position;
