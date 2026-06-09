package service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/jendrix/worldcup-stats-service/internal/domain"
	"github.com/jendrix/worldcup-stats-service/internal/service"
)

type mockMatchRepository struct {
	mock.Mock
}

func (m *mockMatchRepository) ListByYear(ctx context.Context, year int, language string) ([]domain.FixtureMatchRecord, error) {
	args := m.Called(ctx, year, language)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.FixtureMatchRecord), args.Error(1)
}

type mockGroupStatsRepository struct {
	mock.Mock
}

func (m *mockGroupStatsRepository) ListByYear(ctx context.Context, year int, language string) ([]domain.GroupStandingRecord, error) {
	args := m.Called(ctx, year, language)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.GroupStandingRecord), args.Error(1)
}

func TestFixtureService_GetByYear(t *testing.T) {
	ctx := context.Background()

	t.Run("world cup with groups and knockout", func(t *testing.T) {
		matchRepo := new(mockMatchRepository)
		statsRepo := new(mockGroupStatsRepository)
		svc := service.NewFixtureService(matchRepo, statsRepo)

		matches := []domain.FixtureMatchRecord{
			fixtureMatchRecord("group_stage", "A", 1, "group"),
			fixtureMatchRecord("group_stage", "A", 2, "group"),
			fixtureMatchRecord("group_stage", "B", 3, "group"),
			fixtureMatchRecord("group_stage", "C", 4, "group"),
			fixtureMatchRecord("final", "", 5, "knockout"),
		}
		standings := []domain.GroupStandingRecord{
			groupStandingRecord("group_stage", "A", "ARG", 1),
			groupStandingRecord("group_stage", "A", "ITA", 2),
			groupStandingRecord("group_stage", "B", "BRA", 1),
		}

		matchRepo.On("ListByYear", ctx, 1978, "en").Return(matches, nil)
		statsRepo.On("ListByYear", ctx, 1978, "en").Return(standings, nil)

		result, err := svc.GetByYear(ctx, 1978, "en")
		assert.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, 1978, result.Year)
		require.Len(t, result.Stages, 2)
		assert.Equal(t, "group_stage", result.Stages[0].Stage)
		require.Len(t, result.Stages[0].Groups, 3)
		assert.Equal(t, "A", result.Stages[0].Groups[0].GroupCode)
		assert.Len(t, result.Stages[0].Groups[0].Matches, 2)
		assert.Len(t, result.Stages[0].Groups[0].Standings, 2)
		assert.NotNil(t, result.Stages[0].Groups[2].Standings)
		assert.Empty(t, result.Stages[0].Groups[2].Standings)
		assert.Equal(t, "final", result.Stages[1].Stage)
		assert.Len(t, result.Stages[1].Matches, 1)
		assert.Empty(t, result.Stages[1].Groups)
		matchRepo.AssertExpectations(t)
		statsRepo.AssertExpectations(t)
	})

	t.Run("world cup with only knockout stages and replay", func(t *testing.T) {
		matchRepo := new(mockMatchRepository)
		statsRepo := new(mockGroupStatsRepository)
		svc := service.NewFixtureService(matchRepo, statsRepo)
		replayOf := int64(10)

		matches := []domain.FixtureMatchRecord{
			fixtureMatchRecord("round_of_16", "", 10, "knockout"),
			{
				Stage:     "round_of_16",
				GroupCode: "",
				Match: domain.FixtureMatch{
					ID:        11,
					StageType: "knockout",
					ReplayOf:  &replayOf,
				},
			},
		}

		matchRepo.On("ListByYear", ctx, 1938, "es").Return(matches, nil)
		statsRepo.On("ListByYear", ctx, 1938, "es").Return([]domain.GroupStandingRecord{}, nil)

		result, err := svc.GetByYear(ctx, 1938, "es")
		assert.NoError(t, err)
		require.NotNil(t, result)
		require.Len(t, result.Stages, 1)
		assert.Empty(t, result.Stages[0].Groups)
		require.Len(t, result.Stages[0].Matches, 2)
		require.NotNil(t, result.Stages[0].Matches[1].ReplayOf)
		assert.Equal(t, int64(10), *result.Stages[0].Matches[1].ReplayOf)
		matchRepo.AssertExpectations(t)
		statsRepo.AssertExpectations(t)
	})

	t.Run("championship not found", func(t *testing.T) {
		matchRepo := new(mockMatchRepository)
		statsRepo := new(mockGroupStatsRepository)
		svc := service.NewFixtureService(matchRepo, statsRepo)

		matchRepo.On("ListByYear", ctx, 2023, "es").Return([]domain.FixtureMatchRecord{}, nil)

		result, err := svc.GetByYear(ctx, 2023, "es")
		assert.Error(t, err)
		assert.True(t, errors.Is(err, domain.ErrNotFound))
		assert.Nil(t, result)
		statsRepo.AssertNotCalled(t, "ListByYear")
		matchRepo.AssertExpectations(t)
	})

	t.Run("group stats repository error", func(t *testing.T) {
		matchRepo := new(mockMatchRepository)
		statsRepo := new(mockGroupStatsRepository)
		svc := service.NewFixtureService(matchRepo, statsRepo)

		matchRepo.On("ListByYear", ctx, 1978, "es").Return([]domain.FixtureMatchRecord{
			fixtureMatchRecord("group_stage", "A", 1, "group"),
		}, nil)
		statsRepo.On("ListByYear", ctx, 1978, "es").Return(nil, errors.New("db error"))

		result, err := svc.GetByYear(ctx, 1978, "es")
		assert.Error(t, err)
		assert.Nil(t, result)
		matchRepo.AssertExpectations(t)
		statsRepo.AssertExpectations(t)
	})
}

func fixtureMatchRecord(stage, groupCode string, id int64, stageType string) domain.FixtureMatchRecord {
	return domain.FixtureMatchRecord{
		Stage:     stage,
		GroupCode: groupCode,
		Match: domain.FixtureMatch{
			ID:        id,
			StageType: stageType,
		},
	}
}

func groupStandingRecord(stage, groupCode, teamCode string, position int32) domain.GroupStandingRecord {
	return domain.GroupStandingRecord{
		Stage:     stage,
		GroupCode: groupCode,
		Standing: domain.GroupStanding{
			TeamCode: teamCode,
			Position: &position,
		},
	}
}
