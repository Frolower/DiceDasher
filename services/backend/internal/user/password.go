package user

import (
	"context"

	"golang.org/x/crypto/bcrypt"
)

type PasswordManager interface {
	Hash(password string) (string, error)
	Verify(hash, password string) error
}

type BcryptHasher struct {
	Cost int
}

func NewBcryptHasher() BcryptHasher {
	return BcryptHasher{Cost: bcrypt.DefaultCost}
}

func (h BcryptHasher) Hash(password string) (string, error) {
	cost := h.Cost
	if cost == 0 {
		cost = bcrypt.DefaultCost
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), cost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func (h BcryptHasher) Verify(hash, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

func (s *Service) UpdatePassword(ctx context.Context, rec PasswordUpdateRecord) error {
	if s.users == nil {
		return ErrStoreNotConfigured
	}

	err := validatePassword(rec.OldPassword)
	if err != nil {
		return InvalidInputError{Err: err}
	}
	err = validatePassword(rec.NewPassword)
	if err != nil {
		return InvalidInputError{Err: err}
	}

	oldHash, err := s.users.GetHashByID(ctx, rec.Id)
	if err != nil {
		return err
	}

	err = s.passwords.Verify(oldHash, rec.OldPassword)
	if err != nil {
		return ErrOldPasswordDoesntMatch
	}

	newHash, err := s.passwords.Hash(rec.NewPassword)
	if err != nil {
		return err
	}

	err = s.users.UpdatePassword(ctx, rec.Id, oldHash, newHash)
	if err != nil {
		return err
	}

	return nil
}
