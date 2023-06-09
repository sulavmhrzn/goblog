package data

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"errors"
	"time"

	"github.com/sulavmhrzn/goblog/internal/validator"
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

var AnonymousUser = &User{}

func (u *User) IsAnonymous() bool {
	return u == AnonymousUser
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
	query := `INSERT INTO users (email, password, activated) VALUES ($1, $2, $3) RETURNING id`
	args := []interface{}{u.Email, u.Password.hash, u.Activated}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&u.ID)
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

func (m UserModel) GetForToken(tokenScope, tokenPlaintext string) (*User, error) {
	tokenHash := sha256.Sum256([]byte(tokenPlaintext))
	query := `
	SELECT users.id, users.email, users.password, users.activated
	FROM users
	INNER JOIN tokens
	ON users.id = tokens.user_id
	WHERE tokens.hash = $1
	AND tokens.scope = $2
	AND tokens.expiry>$3`
	args := []interface{}{tokenHash[:], tokenScope, time.Now()}

	var user User

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, args...).Scan(
		&user.ID,
		&user.Email,
		&user.Password.hash,
		&user.Activated,
	)
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

func (m UserModel) Update(user *User) error {
	query := `
	UPDATE users
	SET email=$1, password = $2, activated=$3
	WHERE id=$4`
	args := []interface{}{user.Email, user.Password.hash, user.Activated, user.ID}
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

type UserDashboardDetails struct {
	User  User
	Blogs []Blog
}

func (m UserModel) DashboardDetails(userID int) (*UserDashboardDetails, error) {
	blogQuery := `
	SELECT
	blogs.id as blog_id, blogs.title as blog_title, blogs.content as blog_content, 
	blogs.created_at as blog_created_at, blogs.user_id as  blog_user_id, blogs.slug as blog_slug
	FROM blogs
	WHERE blogs.user_id = $1
	`
	userQuery := `
	SELECT id, email, activated FROM users WHERE id = $1`
	var dashboard UserDashboardDetails

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, userQuery, userID).Scan(&dashboard.User.ID, &dashboard.User.Email, &dashboard.User.Activated)
	if err != nil {
		return nil, err
	}

	rows, err := m.DB.QueryContext(ctx, blogQuery, userID)
	if err != nil {
		return nil, err
	}
	var blogs []Blog
	for rows.Next() {
		var blog Blog
		err := rows.Scan(
			&blog.ID,
			&blog.Title,
			&blog.Content,
			&blog.CreatedAt,
			&blog.UserID,
			&blog.Slug,
		)
		if err != nil {
			return nil, err
		}
		blogs = append(blogs, blog)
	}
	err = rows.Close()
	if err != nil {
		return nil, err
	}
	dashboard.Blogs = blogs
	return &dashboard, nil

}
