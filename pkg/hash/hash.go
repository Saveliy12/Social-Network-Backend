package hash

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type HasherInterface interface {
	HashPassword(string) (string, error)
	Compare(string, string) error
}

type Hasher struct {
	salt string
}

func NewHasher(salt string) *Hasher {
	return &Hasher{
		salt: salt,
	}
}

func (h *Hasher) HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hashedPassword), nil
}

func (h *Hasher) Compare(password string, credentials string) (int, error) {
	if err := bcrypt.CompareHashAndPassword([]byte(password), []byte(credentials)); err != nil {
		return 0, errors.New("invalid login or password")
	}

	return 1, nil
}
