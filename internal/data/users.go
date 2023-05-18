package data

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/sulavmhrzn/internal/validator"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrDuplicateEmail = errors.New("dupliicate email")
	ErrNoRows         = errors.New("invalid email or password")
)

type User struct {
	ID        int      `json:"id"`
	Email     string   `json:"email"`
	Password  password `json:"-"`
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

type UserModel struct {
	DB *sql.DB
}

func (m UserModel) Insert(u *User) error {
	query := `INSERT INTO users (email, password, activated) VALUES ($1, $2, $3)`
	args := []interface{}{u.Email, u.Password.hash, u.Activated}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_, err := m.DB.ExecContext(ctx, query, args...)
	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "users_email_key"`:
			return ErrDuplicateEmail
		default:
			return err
		}
	}
	return nil
}

func (m UserModel) GetByEmail(email string) (*User, error) {
	query := `SELECT id, email, password, activated FROM users
	WHERE email = $1`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	var user User
	err := m.DB.QueryRowContext(ctx, query, email).Scan(&user.ID, &user.Email, &user.Password.hash, &user.Activated)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrNoRows
		default:
			return nil, err
		}
	}
	return &user, nil
}
