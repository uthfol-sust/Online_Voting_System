package controllers

import (
	"net/http"
	dto "pollvoting/pkg/DTO"
	"pollvoting/pkg/models"
	"pollvoting/pkg/services"
	"pollvoting/pkg/utils"
	"strconv"
)

type PollController interface {
	CreatePoll(w http.ResponseWriter, r *http.Request)
	GetPolls(w http.ResponseWriter, r *http.Request)
	PollWithOptions(w http.ResponseWriter, r *http.Request)
	UpdatePoll(w http.ResponseWriter, r *http.Request)
	DeletePoll(w http.ResponseWriter, r *http.Request)
}

type pollController struct {
	service services.PollService
}

func NewPollControllers(service services.PollService) PollController {
	return &pollController{service: service}
}

func (c *pollController) CreatePoll(w http.ResponseWriter, r *http.Request) {
	poll := &models.Poll{}
	utils.ParseBody(r, poll)

	newPoll, err := c.service.CreatePoll(poll)
	if err != nil {
		utils.ErrorJSON(w, 500, err)
		return
	}

	res := dto.PollDetails{
		ID:          newPoll.ID,
		Title:       newPoll.Title,
		Description: newPoll.Description,
		CreatedBy:   newPoll.CreatedBy,
		IsActive:    newPoll.IsActive,
		ExpiresAt:   *newPoll.ExpiresAt,
		CreatedAt:   newPoll.CreatedAt,
	}

	utils.JSONResponse(w, 201, "Poll created successfully!", res)
}

func (c *pollController) GetPolls(w http.ResponseWriter, r *http.Request) {
	polls, err := c.service.GetPolls()

	if err != nil {
		utils.ErrorJSON(w, 500, err)
		return
	}

	var pollList []dto.PollDetails
	for _, ele := range polls {
		pollList = append(pollList, dto.PollDetails{
			ID:          ele.ID,
			Title:       ele.Title,
			Description: ele.Description,
			IsActive:    ele.IsActive,
		})
	}

	utils.JSONResponse(w, 200, "Load all the polls successfully", pollList)
}

func (c *pollController) PollWithOptions(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.PathValue("id"))

	poll, err := c.service.PollWithOptions(int64(id))
	if err != nil {
		utils.ErrorJSON(w, 500, err)
		return
	}

	utils.JSONResponse(w, 200, "Poll is found successfully!", poll)
}

func (c *pollController) UpdatePoll(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.PathValue("id"))

	toUpdate := &models.Poll{}
	utils.ParseBody(r, toUpdate)

	poll, err := c.service.UpdatePoll(int64(id), toUpdate)

	if err != nil {
		utils.ErrorJSON(w, 500, err)
		return
	}

	res := dto.PollDetails{
		ID:          poll.ID,
		Title:       poll.Title,
		Description: poll.Description,
		CreatedBy:   poll.CreatedBy,
		IsActive:    poll.IsActive,
		ExpiresAt:   *poll.ExpiresAt,
		CreatedAt:   poll.CreatedAt,
	}

	utils.JSONResponse(w, 201, "Poll is updated successfully!", res)
}

func (c *pollController) DeletePoll(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.PathValue("id"))

	err := c.service.DeletePoll(int64(id))

	if err != nil {
		utils.ErrorJSON(w, 500, err)
		return
	}

	utils.JSONResponse(w, 201, "Poll is deleted successfully!", map[string]interface{}{
		"ID":     id,
		"status": "Deleted",
	})
}
