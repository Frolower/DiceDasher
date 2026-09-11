package handler

type createUserRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type createUserResponse struct {
	Status string `json:"status"`
}
