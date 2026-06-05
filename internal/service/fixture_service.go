package service

import (
	"context"
	"fmt"

	"github.com/jendrix/worldcup-stats-service/internal/domain"
	"github.com/jendrix/worldcup-stats-service/internal/repository"
)

const groupStageType = "group"

// fixtureService implements FixtureService with fixture assembly logic.
type fixtureService struct {
	matchRepo      repository.MatchRepository
	groupStatsRepo repository.GroupStatsRepository
}

// NewFixtureService creates a new FixtureService.
func NewFixtureService(matchRepo repository.MatchRepository, groupStatsRepo repository.GroupStatsRepository) FixtureService {
	return &fixtureService{
		matchRepo:      matchRepo,
		groupStatsRepo: groupStatsRepo,
	}
}

// GetByYear returns the full fixture for a championship year.
func (s *fixtureService) GetByYear(ctx context.Context, year int) (*domain.Fixture, error) {
	matches, err := s.matchRepo.ListByYear(ctx, year)
	if err != nil {
		return nil, err
	}
	if len(matches) == 0 {
		return nil, fmt.Errorf("%w: championship not found", domain.ErrNotFound)
	}

	standings, err := s.groupStatsRepo.ListByYear(ctx, year)
	if err != nil {
		return nil, err
	}

	standingsByGroup := groupStandings(standings)
	stages := assembleStages(matches, standingsByGroup)

	return &domain.Fixture{
		Year:   year,
		Stages: stages,
	}, nil
}

func groupStandings(rows []domain.GroupStandingRecord) map[string][]domain.GroupStanding {
	standings := make(map[string][]domain.GroupStanding)
	for _, row := range rows {
		key := fixtureGroupKey(row.Stage, row.GroupCode)
		standings[key] = append(standings[key], row.Standing)
	}
	return standings
}

func assembleStages(matches []domain.FixtureMatchRecord, standingsByGroup map[string][]domain.GroupStanding) []domain.FixtureStage {
	stages := make([]domain.FixtureStage, 0)
	stageIndexes := make(map[string]int)
	groupIndexes := make(map[string]int)

	for _, record := range matches {
		stageIndex, ok := stageIndexes[record.Stage]
		if !ok {
			stages = append(stages, domain.FixtureStage{Stage: record.Stage})
			stageIndex = len(stages) - 1
			stageIndexes[record.Stage] = stageIndex
		}

		if record.Match.StageType != groupStageType {
			stages[stageIndex].Matches = append(stages[stageIndex].Matches, record.Match)
			continue
		}

		groupKey := fixtureGroupKey(record.Stage, record.GroupCode)
		groupIndex, ok := groupIndexes[groupKey]
		if !ok {
			stages[stageIndex].Groups = append(stages[stageIndex].Groups, domain.FixtureGroup{
				GroupCode: record.GroupCode,
				Matches:   make([]domain.FixtureMatch, 0),
				Standings: standingsForGroup(standingsByGroup, groupKey),
			})
			groupIndex = len(stages[stageIndex].Groups) - 1
			groupIndexes[groupKey] = groupIndex
		}

		stages[stageIndex].Groups[groupIndex].Matches = append(stages[stageIndex].Groups[groupIndex].Matches, record.Match)
	}

	return stages
}

func fixtureGroupKey(stage, groupCode string) string {
	return stage + "\x00" + groupCode
}

func standingsForGroup(standingsByGroup map[string][]domain.GroupStanding, groupKey string) []domain.GroupStanding {
	standings, ok := standingsByGroup[groupKey]
	if !ok {
		return []domain.GroupStanding{}
	}
	return standings
}
