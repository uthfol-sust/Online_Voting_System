package controllers

import (
	"net/http"
	dto "pollvoting/pkg/DTO"
	"pollvoting/pkg/models"
	"pollvoting/pkg/services"
	"pollvoting/pkg/utils"
	"strconv"
)

type OptionController interface {
	CreateOption(w http.ResponseWriter, r *http.Request)
	DeleteOption(w http.ResponseWriter, r *http.Request)
}

type optionController struct {
	service services.OptionService
}

func NewOptionController(service services.OptionService) OptionController {
	return &optionController{service: service}
}

func (c *optionController) CreateOption(w http.ResponseWriter, r *http.Request) {
	option := &models.PollOption{}
	utils.ParseBody(r, option)

	option, err := c.service.Create(option)
	if err != nil {
		utils.ErrorJSON(w, 500, err)
		return
	}

	res := &dto.OptionResponse{
		ID: option.ID,
		Score: option.Score,
		OptionText: option.OptionText,
	}

	utils.JSONResponse(w,201,"New Option Created",res)
}

func (c *optionController) DeleteOption(w http.ResponseWriter, r *http.Request) {
	ID, _ := strconv.Atoi(r.PathValue("id"))

	err := c.service.Delete(int64(ID))
	if err != nil {
		utils.ErrorJSON(w, 500, err)
		return
	}

	res := &dto.DeleteResponse{
		ID:      int64(ID),
		Message: "Deleted",
	}

	utils.JSONResponse(w, 200, "Option Successfully Deleted", res)
}
