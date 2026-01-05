package models

import "time"

type User struct {
	ID        int64     `json:"user_id"`
	Name      string    `json:"user_name"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	Image     string    `json:"image"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
}

type Poll struct {
	ID          int64      `json:"poll_id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	CreatedBy   int64      `json:"created_by"`
	IsActive    bool       `json:"is_active"`
	ExpiresAt   *time.Time `json:"expires_at"`
	CreatedAt   time.Time  `json:"created_at"`
}

type PollOption struct {
	ID         int64    `json:"option_id"`
	PollID     int64    `json:"poll_id"`
	Image      *string `json:"option_image"`
	Score      int    `json:"score"`
	OptionText string `json:"option_text"`
}
