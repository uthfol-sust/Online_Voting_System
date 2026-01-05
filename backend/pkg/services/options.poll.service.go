package services

import (
	"pollvoting/pkg/models"
	"pollvoting/pkg/repositories"
)

type OptionService interface {
	Create(option *models.PollOption) (*models.PollOption, error)
	//Update(options *models.PollOption) (*models.PollOption, error)
	Delete(ID int64) error
}

type optionService struct{
	optionRepo repositories.PollOptionsRepository
}

func NewOptionService(repo repositories.PollOptionsRepository) OptionService{
	return &optionService{optionRepo: repo}
}

func (s *optionService) Create(option *models.PollOption) (*models.PollOption, error){
  return s.optionRepo.Create(option)
}

func (s *optionService) Delete(ID int64) error{
	return s.optionRepo.Delete(ID)
}

