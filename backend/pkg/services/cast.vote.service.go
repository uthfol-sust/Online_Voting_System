package services

import (
	"context"

	dto "pollvoting/pkg/DTO"
	"pollvoting/pkg/repositories"
)

type CastVoteService interface {
	CastVote(ctx context.Context, pollID, optionID, userID int64) error
	PollResults(ctx context.Context, pollID int64) (*dto.PollResultResponse, error)
}

type castVoteService struct {
	repository repositories.CastVoteRepository
}

func NewCastVoteService(repo repositories.CastVoteRepository) CastVoteService {
	return &castVoteService{repository: repo}
}

func (s *castVoteService) CastVote(ctx context.Context, pollID, optionID, userID int64) error {
	return s.repository.CastVote(ctx, pollID, optionID, userID)
}

func (s *castVoteService) PollResults(ctx context.Context, pollID int64) (*dto.PollResultResponse, error) {
	if res, err := s.repository.PollResultsFromRedis(ctx, pollID); err == nil && res != nil {
		return res, err
	}

	return s.repository.PollResultsFromDB(ctx, pollID)
}
