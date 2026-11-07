package user

import "github.com/google/uuid"

type UserUpdateRecord struct {
	ID       uuid.UUID
	Username string
	Email    string
}

type PasswordUpdateRecord struct {
	ID               uuid.UUID
	CurrentSessionID uuid.UUID
	OldPassword      string
	NewPassword      string
}

type UserProfile struct {
	ID       uuid.UUID
	Username string
	Email    string
}
