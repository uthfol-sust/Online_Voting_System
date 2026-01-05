package controllers

import (
	"errors"
	"net/http"
	dto "pollvoting/pkg/DTO"
	"pollvoting/pkg/models"
	"pollvoting/pkg/services"
	"pollvoting/pkg/utils"
	"strconv"
)

type userLogin struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserController interface {
	SingUp(w http.ResponseWriter, r *http.Request)
	Login(w http.ResponseWriter, r *http.Request)
	GetAll(w http.ResponseWriter, r *http.Request)
	GetByID(w http.ResponseWriter, r *http.Request)
	Update(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
	RefreshToken(w http.ResponseWriter, r *http.Request)
}

type userController struct {
	userService services.UserService
}

func NewUserController(service services.UserService) UserController {
	return &userController{userService: service}
}

func (c *userController) SingUp(w http.ResponseWriter, r *http.Request) {
	user := &models.User{}
	utils.ParseBody(r, user)

	newUser, err := c.userService.SingUp(user)
	if err != nil {
		utils.ErrorJSON(w, 500, err)
		return
	}

	res := dto.UserPublicResponse{
		ID:    newUser.ID,
		Name:  newUser.Name,
		Image: newUser.Image,
	}

	utils.JSONResponse(w, 201, "User Created Successfully", res)
}

func (c *userController) Login(w http.ResponseWriter, r *http.Request) {
	user := &userLogin{}
	utils.ParseBody(r, user)

	loggedUser, accessToken, refreshToken, err := c.userService.Login(user.Email, user.Password)
	if err != nil {
		utils.ErrorJSON(w, 500, err)
		return
	}

	err = c.userService.SetToken(loggedUser.ID, refreshToken)
	if err != nil {
		utils.ErrorJSON(w, 500, err)
		return
	}

	utils.JSONResponse(w, 200, "User Login Successfully", map[string]interface{}{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
		"user_email":    user.Email,
	})
}

func (c *userController) GetAll(w http.ResponseWriter, r *http.Request) {
	users, err := c.userService.GetAll()
	if err != nil {
		utils.ErrorJSON(w, 500, err)
		return
	}

	if len(users) == 0 {
		utils.JSONResponse(w, 200, "Empty User List", []dto.UserPublicResponse{})
		return
	}

	var userList []dto.UserPublicResponse
	for _, ele := range users {
		userList = append(userList, dto.UserPublicResponse{
			ID:    ele.ID,
			Name:  ele.Name,
			Image: ele.Image,
		})
	}

	utils.JSONResponse(w, 200, "Users Found Successfully", userList)
}

func (c *userController) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		utils.ErrorJSON(w, 400, err)
		return
	}

	user, err := c.userService.GetByID(int64(id))
	if err != nil {
		utils.ErrorJSON(w, 500, err)
		return
	}

	if user == nil {
		utils.ErrorJSON(w, 404, errors.New("user not found"))
		return
	}

	res := dto.UserPublicResponse{
		ID:    user.ID,
		Name:  user.Name,
		Image: user.Image,
	}

	utils.JSONResponse(w, 200, "User Found Successfully", res)
}

func (c *userController) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		utils.ErrorJSON(w, http.StatusBadRequest, err)
		return
	}

	err = c.userService.Delete(int64(id))
	if err != nil {
		utils.ErrorJSON(w, http.StatusInternalServerError, err)
		return
	}

	res := dto.UserDeleteResponse{
		ID:      int64(id),
		Message: "User deleted successfully",
	}

	utils.JSONResponse(w, http.StatusOK, "Delete successful", res)
}

func (c *userController) Update(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.PathValue("id"))

	var user models.User
	utils.ParseBody(r, &user)

	updatedUser, err := c.userService.Update(int64(id), &user)
	if err != nil {
		utils.ErrorJSON(w, 500, err)
		return
	}
	if updatedUser == nil {
		utils.ErrorJSON(w, 404, errors.New("user not found"))
		return
	}
	res := dto.UserPublicResponse{
		ID:    updatedUser.ID,
		Name:  updatedUser.Name,
		Image: updatedUser.Image,
	}

	utils.JSONResponse(w, 200, "User Updated Successfully", res)
}

func (c *userController) RefreshToken(w http.ResponseWriter, r *http.Request) {
	var body struct {
		RefreshToken string `json:"refresh_token"`
	}
	utils.ParseBody(r, &body)

	access, refresh, err := c.userService.Refresh(body.RefreshToken)
	if err != nil {
		utils.ErrorJSON(w, 401, err)
		return
	}

	utils.JSONResponse(w, 200, "Token refreshed", map[string]string{
		"access_token":  access,
		"refresh_token": refresh,
	})
}
