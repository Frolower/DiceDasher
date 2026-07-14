package handler

import "github.com/google/uuid"

type createUserRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type createUserResponse struct {
	ID uuid.UUID `json:"id"`
}
