package services

import (
	"errors"
	"fmt"
	"pollvoting/pkg/models"
	"pollvoting/pkg/repositories"
	"pollvoting/pkg/utils"
)

type UserService interface {
	SingUp(user *models.User) (*models.User, error)
	Login(email, password string) (*models.User, string, string, error)
	GetAll() ([]models.User, error)
	GetByID(id int64) (*models.User, error)
	Update(id int64, user *models.User) (*models.User, error)
	Delete(id int64) error
	Refresh(refreshToken string) (string, string, error)
	SetToken(userID int64, refresh_token string) error
}

type userService struct {
	userRepo repositories.UserRepository
}

func NewUserService(repo repositories.UserRepository) UserService {
	return &userService{userRepo: repo}
}

func (s *userService) SingUp(user *models.User) (*models.User, error) {
	user.Password, _ = utils.HashingValue(user.Password)

	return s.userRepo.Create(user)
}

func (s *userService) Login(email, password string) (*models.User, string, string, error) {
	user, err := s.userRepo.GetByEmail(email)
	if err != nil {
		return nil, "", "", err
	}
	if user == nil {
		return nil, "", "", errors.New("user not found!")
	}

	fmt.Println(user.Password)

	isValid := utils.ChcekPassword(user.Password, password)
	if isValid == false {
		return nil, "", "", errors.New("Wrong Password")
	}

	accessToken, err := utils.GenerateAccessToken(user)
	if err != nil {
		return nil, "", "", err
	}
	refreshToken, err := utils.GenerateRefreshToken(user)
	if err != nil {
		return nil, "", "", err
	}

	return user, accessToken, refreshToken, nil
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
		hashed, _ := utils.HashingValue(user.Password)
		existingUser.Password = hashed
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

func (s *userService) SetToken(userID int64, refresh_token string) error {
	return s.userRepo.TokenSetToRDB(userID, refresh_token)
}

func (s *userService) Refresh(refreshToken string) (string, string, error) {

	claims, err := utils.VerifyRefreshToken(refreshToken)
	if err != nil {
		return "", "", errors.New("invalid refresh token")
	}

	fmt.Println(claims.ID)

	storedToken, err := s.userRepo.GetTokenFromRDB(claims.ID)
	if err != nil || storedToken == "" {
		return "", "", errors.New("refresh token not found or revoked")
	}

	if storedToken != refreshToken {
		return "", "", errors.New("token mismatch detected")
	}

	user, _ := s.userRepo.GetByID(claims.ID)

	newAccess, _ := utils.GenerateAccessToken(user)
	newRefresh, _ := utils.GenerateRefreshToken(user)

	s.userRepo.TokenSetToRDB(user.ID, newRefresh)

	return newAccess, newRefresh, nil
}

