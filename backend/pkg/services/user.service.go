package services

import (
	"errors"
	"pollvoting/pkg/models"
	"pollvoting/pkg/repositories"
)

type UserService interface {
	SingUp(user *models.User) (*models.User, error)
	GetAll() ([]models.User, error)
	GetByID(id int64) (*models.User, error)
	Update(id int64, user *models.User) (*models.User, error)
	Delete(id int64) error
}

type userService struct {
	userRepo repositories.UserRepository
}

func NewUserService(repo repositories.UserRepository) UserService {
	return &userService{userRepo: repo}
}

func (s *userService) SingUp(user *models.User) (*models.User, error) {
	return s.userRepo.Create(user)
}

func (s *userService) GetAll() ([]models.User, error) {
	return s.userRepo.GetAll()
}

func (s *userService) GetByID(id int64) (*models.User, error) {
	return s.userRepo.GetByID(id)
}

func (s *userService) Update(id int64, user *models.User) (*models.User, error) {
	existingUser, err := s.userRepo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if existingUser == nil {
		return nil, errors.New("user not found")
	}

	if user.Name != "" {
		existingUser.Name = user.Name
	}
	if user.Email != "" {
		existingUser.Email = user.Email
	}
	if user.Image != "" {
		existingUser.Image = user.Image
	}
	if user.Role != "" {
		existingUser.Role = user.Role
	}
	if user.Password != "" {
		existingUser.Password = user.Password
	}

	err = s.userRepo.Update(id, existingUser)
	if err != nil {
		return nil, err
	}

	return existingUser, nil
}
func (s *userService) Delete(id int64) error {
	existingUser, err := s.userRepo.GetByID(id)
	if err != nil {
		return err
	}
	if existingUser == nil {
		return errors.New("user not found")
	}

	return s.userRepo.Delete(id)
}
