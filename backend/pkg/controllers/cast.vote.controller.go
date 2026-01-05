package controllers

import (
	"errors"
	"net/http"
	"strconv"

	"pollvoting/pkg/middleware"
	"pollvoting/pkg/services"
	"pollvoting/pkg/utils"
)

type CastVoteController interface {
	CastVote(w http.ResponseWriter, r *http.Request)
	PollResults(w http.ResponseWriter, r *http.Request)
}

type castVoteController struct {
	service services.CastVoteService
}

type castVoteRequest struct {
	PollID   int64 `json:"poll_id"`
	OptionID int64 `json:"option_id"`
}

func NewCastVoteController(service services.CastVoteService) CastVoteController {
	return &castVoteController{service: service}
}

func (c *castVoteController) CastVote(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())

	req := &castVoteRequest{}
	utils.ParseBody(r, req)

	if req.PollID <= 0 || req.OptionID <= 0 {
		utils.ErrorJSON(w, http.StatusBadRequest, errors.New("invalid poll_id or option_id"))
		return
	}

	if err := c.service.CastVote(r.Context(), req.PollID, req.OptionID, userID); err != nil {
		utils.ErrorJSON(w, http.StatusBadRequest, err)
		return
	}

	utils.JSONResponse(w, http.StatusOK, "Vote cast successfully", nil)
}

func (c *castVoteController) PollResults(w http.ResponseWriter, r *http.Request) {
	pollID, err := strconv.ParseInt(r.PathValue("pollID"), 10, 64)
	if err != nil || pollID <= 0 {
		utils.ErrorJSON(w, http.StatusBadRequest, errors.New("invalid poll id"))
		return
	}

	result, err := c.service.PollResults(r.Context(), pollID)
	if err != nil {
		utils.ErrorJSON(w, http.StatusInternalServerError, err)
		return
	}

	utils.JSONResponse(w, http.StatusOK, "Poll results fetched", result)
}
