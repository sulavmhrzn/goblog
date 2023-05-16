package data

import (
	"errors"

	"github.com/sulavmhrzn/internal/validator"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        int      `json:"id"`
	Email     string   `json:"email"`
	Password  password `json:"password"`
	Activated bool     `json:"activated"`
}

type password struct {
	plaintext *string
	hash      []byte
}

func (p *password) Set(plaintextPassword string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(plaintextPassword), 12)
	if err != nil {
		return err
	}
	p.hash = hash
	p.plaintext = &plaintextPassword
	return nil
}

func (p *password) Matches(plaintextPassword string) (bool, error) {
	err := bcrypt.CompareHashAndPassword(p.hash, []byte(plaintextPassword))
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

func ValidateEmail(v *validator.Validator, email string) {
	v.Check(email != "", "email", "must be provided")
	v.Check(v.ValidateEmail(email), "email", "must be a valid email address")
}

func ValidatePlaintextPassword(v *validator.Validator, password string) {
	v.Check(password != "", "password", "must be provided")
	v.Check(len(password) >= 8, "password", "must be greater than 8 characters")
	v.Check(len(password) <= 72, "password", "must be less than 72 characters")

}

func ValidateUser(v *validator.Validator, user *User) {
	ValidateEmail(v, user.Email)
	ValidatePlaintextPassword(v, *user.Password.plaintext)
}
