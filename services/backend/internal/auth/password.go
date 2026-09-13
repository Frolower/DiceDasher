package auth

import "golang.org/x/crypto/bcrypt"

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
