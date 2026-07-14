package user

import "github.com/google/uuid"

type CreateInput struct {
	Username string
	Email    string
	Password string
}

type CreateRecord struct {
	Username     string
	Email        string
	PasswordHash string
}

type Created struct {
	ID uuid.UUID
}
