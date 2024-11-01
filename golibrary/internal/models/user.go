package models

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID          int64    `json:"id"`
	Name        string   `json:"name"`
	Email       string   `json:"email"`
	RentedBooks []Book   `json:"rented_books,omitempty"`
	Password    Password `json:"-"`
	Deleted     bool     `json:"deleted"`
	Version     int      `json:"-"`
}

type Password struct {
	Plaintext *string
	Hash      []byte
}

func (p *Password) Set(plaintextPassword string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(plaintextPassword), 12)
	if err != nil {
		return err
	}

	p.Plaintext = &plaintextPassword
	p.Hash = hash
	return nil
}

func (p *Password) Matches(plaintextPassword string) (bool, error) {
	err := bcrypt.CompareHashAndPassword(p.Hash, []byte(plaintextPassword))
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return false, nil
		default:
			return false, err
		}
	}

	return true, nil
}
