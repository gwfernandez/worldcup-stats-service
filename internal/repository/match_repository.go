package repository

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgtype"

	"github.com/jendrix/worldcup-stats-service/db/sqlc"
	"github.com/jendrix/worldcup-stats-service/internal/domain"
)

// matchRepository implements MatchRepository using sqlc-generated queries.
type matchRepository struct {
	queries *sqlc.Queries
}

// NewMatchRepository creates a new repository backed by sqlc.
func NewMatchRepository(db sqlc.DBTX) MatchRepository {
	return &matchRepository{queries: sqlc.New(db)}
}

func (r *matchRepository) ListByYear(ctx context.Context, year int, language string) ([]domain.FixtureMatchRecord, error) {
	rows, err := r.queries.ListMatchesByYear(ctx, int32(year))
	if err != nil {
		return nil, err
	}

	matches := make([]domain.FixtureMatchRecord, len(rows))
	for i, row := range rows {
		matches[i] = toFixtureMatchRecord(row)
	}

	return matches, nil
}

func toFixtureMatchRecord(row sqlc.ListMatchesByYearRow) domain.FixtureMatchRecord {
	return domain.FixtureMatchRecord{
		Stage:     enumString(row.Stage),
		GroupCode: strings.ToUpper(textString(row.GroupCode)),
		Match: domain.FixtureMatch{
			ID:                     row.ID,
			StageType:              enumString(row.StageType),
			Replayed:               row.Replayed,
			ReplayOf:               int8Ptr(row.ReplayOf),
			MatchDate:              datePtr(row.MatchDate),
			MatchTime:              timePtr(row.MatchTime),
			StadiumID:              int8Ptr(row.StadiumID),
			HomeTeam:               domain.SimpleTeam{Code: strings.ToUpper(row.HomeTeamCode)},
			AwayTeam:               domain.SimpleTeam{Code: strings.ToUpper(row.AwayTeamCode)},
			HomeTeamScore:          int4Ptr(row.HomeTeamScore),
			AwayTeamScore:          int4Ptr(row.AwayTeamScore),
			ExtraTime:              row.ExtraTime,
			PenaltyShootout:        row.PenaltyShootout,
			HomeTeamScorePenalties: int4Ptr(row.HomeTeamScorePenalties),
			AwayTeamScorePenalties: int4Ptr(row.AwayTeamScorePenalties),
			HomeTeamWin:            boolPtr(row.HomeTeamWin),
			AwayTeamWin:            boolPtr(row.AwayTeamWin),
			Draw:                   boolPtr(row.Draw),
			RefID:                  textPtr(row.RefID),
		},
	}
}

func enumString(value interface{}) string {
	if value == nil {
		return ""
	}
	return fmt.Sprint(value)
}

func textString(value pgtype.Text) string {
	if !value.Valid {
		return ""
	}
	return value.String
}

func textPtr(value pgtype.Text) *string {
	if !value.Valid {
		return nil
	}
	return &value.String
}

func int4Ptr(value pgtype.Int4) *int32 {
	if !value.Valid {
		return nil
	}
	return &value.Int32
}

func int8Ptr(value pgtype.Int8) *int64 {
	if !value.Valid {
		return nil
	}
	return &value.Int64
}

func boolPtr(value pgtype.Bool) *bool {
	if !value.Valid {
		return nil
	}
	return &value.Bool
}

func datePtr(value pgtype.Date) *string {
	if !value.Valid {
		return nil
	}
	formatted := value.Time.Format("2006-01-02")
	return &formatted
}

func timePtr(value pgtype.Time) *string {
	if !value.Valid {
		return nil
	}
	t := time.UnixMicro(value.Microseconds).UTC()
	formatted := t.Format("15:04:05")
	return &formatted
}
