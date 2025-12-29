package services

import (
	"errors"
	"pollvoting/pkg/models"
	"pollvoting/pkg/repositories"
)

type PollService interface {
	CreatePoll(poll *models.Poll) (*models.Poll, error)
	GetPolls() ([]models.Poll, error)
	GetPollDetails(id int64) (*models.Poll, error)
	UpdatePoll(id int64, poll *models.Poll) (*models.Poll, error)
	DeletePoll(id int64) error
}

type pollService struct {
	reposi repositories.PollRepository
}

func NewPollService(repo repositories.PollRepository) PollService {
	return &pollService{reposi: repo}
}

func (s *pollService) CreatePoll(poll *models.Poll) (*models.Poll, error) {
	return s.reposi.CreatePoll(poll)
}
func (s *pollService) GetPolls() ([]models.Poll, error) {
	return s.reposi.GetPolls()
}
func (s *pollService) GetPollDetails(id int64) (*models.Poll, error) {
	isExist, err := s.reposi.GetPollDetails(id)
	if err != nil {
		return nil, err
	}

	if isExist == nil {
		return nil, errors.New("this polls is not Created")
	}

	return s.reposi.GetPollDetails(id)
}

func (s *pollService) UpdatePoll(id int64, poll *models.Poll) (*models.Poll, error) {
	isExist, err := s.reposi.GetPollDetails(id)
	if err != nil {
		return nil, err
	}

	if isExist == nil{
		return nil, errors.New("this polls is not Created")
	}

	if poll.Title != "" {
		isExist.Title = poll.Title
	}
	if poll.Description != "" {
		isExist.Description = poll.Description
	}
	if poll.IsActive {
		isExist.IsActive = poll.IsActive
	}
	if poll.ExpiresAt !=nil{
		isExist.ExpiresAt = poll.ExpiresAt
	}

	return s.reposi.UpdatePoll(id, isExist)

}
func (s *pollService) DeletePoll(id int64) error {
	isExist, err := s.reposi.GetPollDetails(id)
	if err != nil {
		return err
	}

	if isExist == nil {
		return errors.New("this polls is not Found to Delete")
	}

	return s.reposi.DeletePoll(id)
}
