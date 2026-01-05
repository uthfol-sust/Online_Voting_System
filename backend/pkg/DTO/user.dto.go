package dto

import "time"

type UserPublicResponse struct {
	ID    int64  `json:"id"`
	Name  string `json:"name"`
	Image string `json:"image"`
}

type UserProfileResponse struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Image     string    `json:"image"`
	CreatedAt time.Time `json:"created_at"`
}

type DeleteResponse struct {
    ID      int64  `json:"id"`
    Message string `json:"message"`
}
