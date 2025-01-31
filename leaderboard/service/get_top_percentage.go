package service

import (
	"context"
	"math"

	"github.com/topfreegames/podium/leaderboard/model"
)

const getTopPercentageServiceLabel = "get top percentage"

// GetTopPercentage retrieves top x% members from the leaderboard.
func (s *Service) GetTopPercentage(ctx context.Context, leaderboardID string, pageSize, amount, maxMembers int, order string) ([]*model.Member, error) {
	if amount < 1 || amount > 100 {
		return nil, NewGeneralError(getTopPercentageServiceLabel, "Percentage must be a valid integer between 1 and 100")
	}

	if order != "desc" && order != "asc" {
		return nil, NewInvalidOrderError(order)
	}

	amountInPercentage := float64(amount) / 100.0
	totalNumberMembers, err := s.Database.GetTotalMembers(ctx, leaderboardID)
	if err != nil {
		return nil, NewGeneralError(getTopPercentageServiceLabel, err.Error())
	}

	numberMembersToReturn := int(math.Floor(float64(totalNumberMembers) * amountInPercentage))

	if numberMembersToReturn < 1 {
		numberMembersToReturn = 1
	}

	if numberMembersToReturn > maxMembers {
		numberMembersToReturn = maxMembers
	}

	databaseMembers, err := s.Database.GetOrderedMembers(ctx, leaderboardID, 0, numberMembersToReturn-1, order)
	if err != nil {
		return nil, NewGeneralError(getTopPercentageServiceLabel, err.Error())
	}

	members := convertDatabaseMembersIntoModelMembers(databaseMembers)
	return members, nil
}
