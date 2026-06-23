package repository

import (
	"context"
	"errors"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"

	"github.com/jendrix/worldcup-stats-service/db/sqlc"
	"github.com/jendrix/worldcup-stats-service/internal/domain"
)

// championshipRepository implements ChampionshipRepository using sqlc.
type championshipRepository struct {
	queries *sqlc.Queries
}

// NewChampionshipRepository creates a new ChampionshipRepository.
func NewChampionshipRepository(db sqlc.DBTX) ChampionshipRepository {
	return &championshipRepository{queries: sqlc.New(db)}
}

// List retrieves a paginated list of championships based on the given filters.
func (r *championshipRepository) List(ctx context.Context, filter domain.ChampionshipFilter) ([]domain.Championship, int64, error) {
	total, err := r.countChampionships(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	limit := int32(filter.Size)
	offset := int32((filter.Page - 1) * filter.Size)

	championships, err := r.listChampionships(ctx, filter, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	return championships, total, nil
}

func (r *championshipRepository) countChampionships(ctx context.Context, filter domain.ChampionshipFilter) (int64, error) {
	if filter.Host == "" {
		return r.queries.CountChampionshipsWithoutHostFilter(ctx, sqlc.CountChampionshipsWithoutHostFilterParams{
			Column1: int32(filter.Year),
			Column2: filter.ConfederationCode,
		})
	}
	return r.queries.CountChampionships(ctx, sqlc.CountChampionshipsParams{
		Column1:  int32(filter.Year),
		Column2:  filter.Host,
		Column3:  filter.ConfederationCode,
		Language: filter.Language,
	})
}

func (r *championshipRepository) listChampionships(ctx context.Context, filter domain.ChampionshipFilter, limit int32, offset int32) ([]domain.Championship, error) {
	if filter.Host == "" {
		rows, err := r.queries.ListChampionshipsWithoutHostFilter(ctx, sqlc.ListChampionshipsWithoutHostFilterParams{
			Column1: int32(filter.Year),
			Column2: filter.ConfederationCode,
			Limit:   limit,
			Offset:  offset,
		})
		if err != nil {
			return nil, err
		}
		championships := make([]domain.Championship, len(rows))
		for i, row := range rows {
			championships[i] = toChampionshipDomainFromFields(
				row.Year,
				row.StartDate,
				row.EndDate,
				row.HostCodes,
				row.ConfederationCodes,
				row.ChampionCode,
			)
		}
		return championships, nil
	}

	rows, err := r.queries.ListChampionships(ctx, sqlc.ListChampionshipsParams{
		Column1:  int32(filter.Year),
		Column2:  filter.Host,
		Column3:  filter.ConfederationCode,
		Limit:    limit,
		Offset:   offset,
		Language: filter.Language,
	})
	if err != nil {
		return nil, err
	}
	championships := make([]domain.Championship, len(rows))
	for i, row := range rows {
		championships[i] = toChampionshipDomain(row)
	}
	return championships, nil
}

// GetByYear retrieves a championship and its stats by year.
func (r *championshipRepository) GetByYear(ctx context.Context, year int) (*domain.Championship, error) {
	row, err := r.queries.GetChampionshipByYear(ctx, int32(year))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	championship := toChampionshipDetailDomain(row)
	return &championship, nil
}

// ListTeamTranslations retrieves all team translations used to hydrate championship hosts.
func (r *championshipRepository) ListTeamTranslations(ctx context.Context) ([]domain.TeamTranslation, error) {
	rows, err := r.queries.ListTeamTranslations(ctx)
	if err != nil {
		return nil, err
	}

	translations := make([]domain.TeamTranslation, len(rows))
	for i, row := range rows {
		translations[i] = domain.TeamTranslation{
			TeamCode: strings.ToUpper(row.TeamCode),
			Language: row.Language,
			Name:     row.Name,
		}
	}

	return translations, nil
}

// ListTeamsByYear retrieves a paginated list of teams that participated in a championship year.
func (r *championshipRepository) ListTeamsByYear(ctx context.Context, filter domain.ChampionshipTeamFilter) ([]domain.ChampionshipTeam, int64, error) {
	total, err := r.countChampionshipTeamsByYear(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	limit := int32(filter.Size)
	offset := int32((filter.Page - 1) * filter.Size)

	teams, err := r.listChampionshipTeamsByYear(ctx, filter, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	return teams, total, nil
}

func (r *championshipRepository) countChampionshipTeamsByYear(ctx context.Context, filter domain.ChampionshipTeamFilter) (int64, error) {
	if filter.Name == "" {
		return r.queries.CountChampionshipTeamsByYearWithoutNameFilter(ctx, sqlc.CountChampionshipTeamsByYearWithoutNameFilterParams{
			Year:    int32(filter.Year),
			Column2: filter.ConfederationCode,
			Column3: filter.GroupCode,
		})
	}
	return r.queries.CountChampionshipTeamsByYear(ctx, sqlc.CountChampionshipTeamsByYearParams{
		Year:     int32(filter.Year),
		Column2:  filter.Name,
		Column3:  filter.ConfederationCode,
		Column4:  filter.GroupCode,
		Language: filter.Language,
	})
}

func (r *championshipRepository) listChampionshipTeamsByYear(ctx context.Context, filter domain.ChampionshipTeamFilter, limit int32, offset int32) ([]domain.ChampionshipTeam, error) {
	if filter.Name == "" {
		rows, err := r.queries.ListChampionshipTeamsByYearWithoutNameFilter(ctx, sqlc.ListChampionshipTeamsByYearWithoutNameFilterParams{
			Year:    int32(filter.Year),
			Column2: filter.ConfederationCode,
			Column3: filter.GroupCode,
			Limit:   limit,
			Offset:  offset,
		})
		if err != nil {
			return nil, err
		}
		teams := make([]domain.ChampionshipTeam, len(rows))
		for i, row := range rows {
			teams[i] = toChampionshipTeamDomainFromFields(
				row.Year,
				row.TeamCode,
				row.ConfederationCode,
				row.GroupCode,
				row.StageReached,
				row.Managers,
			)
		}
		return teams, nil
	}

	rows, err := r.queries.ListChampionshipTeamsByYear(ctx, sqlc.ListChampionshipTeamsByYearParams{
		Year:     int32(filter.Year),
		Column2:  filter.Name,
		Column3:  filter.ConfederationCode,
		Column4:  filter.GroupCode,
		Limit:    limit,
		Offset:   offset,
		Language: filter.Language,
	})
	if err != nil {
		return nil, err
	}
	teams := make([]domain.ChampionshipTeam, len(rows))
	for i, row := range rows {
		teams[i] = toChampionshipTeamDomain(row)
	}
	return teams, nil
}

// ListStadiumsByYear retrieves a paginated list of stadiums used in a championship year.
func (r *championshipRepository) ListStadiumsByYear(ctx context.Context, filter domain.ChampionshipStadiumFilter) ([]domain.ChampionshipStadium, int64, error) {
	total, err := r.queries.CountChampionshipStadiumsByYear(ctx, sqlc.CountChampionshipStadiumsByYearParams{
		Year:    int32(filter.Year),
		Column2: filter.Name,
	})
	if err != nil {
		return nil, 0, err
	}

	limit := int32(filter.Size)
	offset := int32((filter.Page - 1) * filter.Size)

	rows, err := r.queries.ListChampionshipStadiumsByYear(ctx, sqlc.ListChampionshipStadiumsByYearParams{
		Year:    int32(filter.Year),
		Column2: filter.Name,
		Limit:   limit,
		Offset:  offset,
	})
	if err != nil {
		return nil, 0, err
	}

	stadiums := make([]domain.ChampionshipStadium, len(rows))
	for i, row := range rows {
		stadiums[i] = toChampionshipStadiumDomain(row)
	}

	return stadiums, total, nil
}

// ListScorersByYear retrieves a paginated list of scorers for a championship year.
func (r *championshipRepository) ListScorersByYear(ctx context.Context, filter domain.ChampionshipScorerFilter) ([]domain.ChampionshipScorer, int64, error) {
	total, err := r.queries.CountChampionshipScorersByYear(ctx, sqlc.CountChampionshipScorersByYearParams{
		Year:    int32(filter.Year),
		Column2: filter.Name,
		Column3: filter.TeamCode,
	})
	if err != nil {
		return nil, 0, err
	}

	limit := int32(filter.Size)
	offset := int32((filter.Page - 1) * filter.Size)

	rows, err := r.queries.ListChampionshipScorersByYear(ctx, sqlc.ListChampionshipScorersByYearParams{
		Year:        int32(filter.Year),
		Column2:     filter.Name,
		Column3:     filter.TeamCode,
		OffsetValue: offset,
		LimitValue:  limit,
	})
	if err != nil {
		return nil, 0, err
	}

	scorers := make([]domain.ChampionshipScorer, len(rows))
	for i, row := range rows {
		scorers[i] = toChampionshipScorerDomain(row)
	}

	return scorers, total, nil
}

// ListSquadByYearAndTeam retrieves a paginated list of squad players for a team in a championship year.
func (r *championshipRepository) ListSquadByYearAndTeam(ctx context.Context, filter domain.ChampionshipSquadFilter) ([]domain.ChampionshipSquadPlayer, int64, error) {
	total, err := r.queries.CountChampionshipSquadByYearAndTeam(ctx, sqlc.CountChampionshipSquadByYearAndTeamParams{
		Year:     int32(filter.Year),
		TeamCode: filter.TeamCode,
	})
	if err != nil {
		return nil, 0, err
	}

	limit := int32(filter.Size)
	offset := int32((filter.Page - 1) * filter.Size)

	rows, err := r.queries.ListChampionshipSquadByYearAndTeam(ctx, sqlc.ListChampionshipSquadByYearAndTeamParams{
		Year:        int32(filter.Year),
		TeamCode:    filter.TeamCode,
		OffsetValue: offset,
		LimitValue:  limit,
	})
	if err != nil {
		return nil, 0, err
	}

	players := make([]domain.ChampionshipSquadPlayer, len(rows))
	for i, row := range rows {
		players[i] = toChampionshipSquadPlayerDomain(row)
	}

	return players, total, nil
}

// ListStandingsByYear retrieves a paginated list of standings for a championship year.
func (r *championshipRepository) ListStandingsByYear(ctx context.Context, filter domain.ChampionshipStandingFilter) ([]domain.ChampionshipStanding, int64, error) {
	total, err := r.queries.CountChampionshipStandingsByYear(ctx, int32(filter.Year))
	if err != nil {
		return nil, 0, err
	}

	limit := int32(filter.Size)
	offset := int32((filter.Page - 1) * filter.Size)

	rows, err := r.queries.ListChampionshipStandingsByYear(ctx, sqlc.ListChampionshipStandingsByYearParams{
		Year:   int32(filter.Year),
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return nil, 0, err
	}

	standings := make([]domain.ChampionshipStanding, len(rows))
	for i, row := range rows {
		standings[i] = toChampionshipStandingDomain(row)
	}

	return standings, total, nil
}

func toChampionshipDomain(row sqlc.ListChampionshipsRow) domain.Championship {
	return toChampionshipDomainFromFields(
		row.Year,
		row.StartDate,
		row.EndDate,
		row.HostCodes,
		row.ConfederationCodes,
		row.ChampionCode,
	)
}

func toChampionshipDomainFromFields(year int32, startDate pgtype.Date, endDate pgtype.Date, hostCodes []string, confederationCodes []string, championCodeValue pgtype.Text) domain.Championship {
	var championCode *string
	if championCodeValue.Valid {
		val := strings.ToUpper(championCodeValue.String)
		championCode = &val
	}

	return domain.Championship{
		Year:               int(year),
		StartDate:          dateToString(startDate),
		EndDate:            dateToString(endDate),
		HostCodes:          uppercaseSlice(hostCodes),
		ConfederationCodes: uppercaseSlice(confederationCodes),
		ChampionCode:       championCode,
	}
}

func toChampionshipDetailDomain(row sqlc.GetChampionshipByYearRow) domain.Championship {
	var championCode *string
	if row.ChampionCode.Valid {
		val := strings.ToUpper(row.ChampionCode.String)
		championCode = &val
	}

	c := domain.Championship{
		Year:               int(row.Year),
		StartDate:          dateToString(row.StartDate),
		EndDate:            dateToString(row.EndDate),
		HostCodes:          uppercaseSlice(row.HostCodes),
		ConfederationCodes: uppercaseSlice(row.ConfederationCodes),
		ChampionCode:       championCode,
	}

	// If TotalTeams is valid, it means we have stats in the DB
	if row.TotalTeams.Valid {
		var runnerUpCode string
		if row.StatsRunnerUpCode.Valid {
			runnerUpCode = strings.ToUpper(row.StatsRunnerUpCode.String)
		}

		var thirdPlaceCode string
		if row.StatsThirdPlaceCode.Valid {
			thirdPlaceCode = strings.ToUpper(row.StatsThirdPlaceCode.String)
		}

		var fourthPlaceCode string
		if row.StatsFourthPlaceCode.Valid {
			fourthPlaceCode = strings.ToUpper(row.StatsFourthPlaceCode.String)
		}

		c.Stats = &domain.ChampionshipsStats{
			TotalTeams:      row.TotalTeams.Int32,
			TotalMatches:    row.TotalMatches.Int32,
			TotalStadiums:   row.TotalStadiums.Int32,
			TotalPlayers:    row.TotalPlayers.Int32,
			TotalGoals:      row.TotalGoals.Int32,
			RunnerUpCode:    runnerUpCode,
			ThirdPlaceCode:  thirdPlaceCode,
			FourthPlaceCode: fourthPlaceCode,
			TopScorers:      make([]domain.TopScorer, 0), // Default empty slice until players table exists
			TopScorerGoals:  row.TopScorerGoals.Int32,
		}
	}

	return c
}

func toChampionshipTeamDomain(row sqlc.ListChampionshipTeamsByYearRow) domain.ChampionshipTeam {
	return toChampionshipTeamDomainFromFields(
		row.Year,
		row.TeamCode,
		row.ConfederationCode,
		row.GroupCode,
		row.StageReached,
		row.Managers,
	)
}

func toChampionshipTeamDomainFromFields(year int32, teamCode string, confederationCode string, groupCodeValue pgtype.Text, stageReached string, managers string) domain.ChampionshipTeam {
	groupCode := ""
	if groupCodeValue.Valid {
		groupCode = groupCodeValue.String
	}
	return domain.ChampionshipTeam{
		Year:              int(year),
		Team:              domain.SimpleTeam{Code: strings.ToUpper(teamCode)},
		ConfederationCode: strings.ToUpper(confederationCode),
		GroupCode:         strings.ToUpper(groupCode),
		StageReached:      stageReached,
		Managers:          managers,
	}
}

func toChampionshipStadiumDomain(row sqlc.ListChampionshipStadiumsByYearRow) domain.ChampionshipStadium {
	return domain.ChampionshipStadium{
		Year:          int(row.Year),
		ID:            row.ID,
		Name:          row.Name,
		CityName:      row.CityName,
		Capacity:      row.Capacity,
		MatchesPlayed: row.MatchesPlayed,
	}
}

func toChampionshipScorerDomain(row sqlc.ListChampionshipScorersByYearRow) domain.ChampionshipScorer {
	return domain.ChampionshipScorer{
		PlayerID: row.PlayerID,
		FullName: row.FullName,
		Team:     domain.SimpleTeam{Code: strings.ToUpper(row.TeamCode)},
		Goals:    row.Goals,
	}
}

func toChampionshipSquadPlayerDomain(row sqlc.ListChampionshipSquadByYearAndTeamRow) domain.ChampionshipSquadPlayer {
	var position *string
	if row.Position.Valid {
		position = &row.Position.String
	}

	var shirtNumber *int32
	if row.ShirtNumber.Valid {
		shirtNumber = &row.ShirtNumber.Int32
	}

	return domain.ChampionshipSquadPlayer{
		PlayerID:    row.PlayerID,
		FirstName:   row.FirstName,
		LastName:    row.LastName,
		Position:    position,
		ShirtNumber: shirtNumber,
	}
}

func toChampionshipStandingDomain(row sqlc.ListChampionshipStandingsByYearRow) domain.ChampionshipStanding {
	return domain.ChampionshipStanding{
		Team:           domain.SimpleTeam{Code: strings.ToUpper(row.TeamCode)},
		GroupCode:      strings.ToUpper(row.GroupCode),
		MatchesPlayed:  row.MatchesPlayed,
		Wins:           row.Wins,
		Draws:          row.Draws,
		Losses:         row.Losses,
		GoalsFor:       row.GoalsFor,
		GoalsAgainst:   row.GoalsAgainst,
		GoalDifference: row.GoalDifference.Int32,
		Points:         row.Points,
		UnifiedPoints:  row.UnifiedPoints,
		Position:       row.Position.Int32,
		Performance:    row.Performance,
	}
}

func dateToString(date pgtype.Date) string {
	if !date.Valid {
		return ""
	}
	return date.Time.Format("2006-01-02")
}

func uppercaseSlice(slice []string) []string {
	if slice == nil {
		return nil
	}
	res := make([]string, len(slice))
	for i, val := range slice {
		res[i] = strings.ToUpper(val)
	}
	return res
}
