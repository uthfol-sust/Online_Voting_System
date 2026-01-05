package dto

import (
	"pollvoting/pkg/models"
	"time"
)

type PollDetails struct {
	ID          int64     `json:"poll_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	CreatedBy   int64     `json:"created_by"`
	IsActive    bool      `json:"is_active"`
	ExpiresAt   time.Time `json:"expires_at"`
	CreatedAt   time.Time `json:"created_at"`
}

type PollResponse struct {
	ID          int64                 `json:"id"`
	Title       string              `json:"title"`
	Description string              `json:"description"`
	IsActive    bool                `json:"is_active"`
	ExpiresAt   time.Time           `json:"expires_at"`
	Options     []models.PollOption `json:"options"`
}

type OptionResponse struct {
	ID         int64    `json:"option_id"`
	Score      int    `json:"score"`
	OptionText string `json:"option_text"`
}

type PollResultResponse struct {
	PollID  int64              `json:"poll_id"`
	Results []PollOptionResult `json:"results"`
}

type PollOptionResult struct {
	OptionID int64 `json:"option_id"`
	Votes    int64 `json:"votes"`
}

