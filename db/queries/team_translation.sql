-- name: ListTeamTranslations :many
SELECT
    team_code,
    language,
    name
FROM team_translations
ORDER BY language ASC, team_code ASC;
