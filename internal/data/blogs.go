package data

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/sulavmhrzn/internal/validator"
)

type Blog struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UserID    int       `json:"-"`
	Slug      string    `json:"slug"`
}

func ValidateBlog(v *validator.Validator, blog *Blog) {
	v.Check(blog.Title != "", "title", "must be provided")
	v.Check(len(blog.Title) >= 2, "title", "must be greater than 2 characters")
	v.Check(blog.Content != "", "content", "must be provided")
	v.Check(len(blog.Content) >= 5, "content", "must be greater than 5 characters")

}

type BlogModel struct {
	DB *sql.DB
}

func (m BlogModel) Insert(b *Blog) error {
	query := `
	INSERT INTO blogs (title, content, created_at, user_id, slug)
	VALUES ($1, $2, $3, $4, $5) RETURNING id`
	args := []interface{}{b.Title, b.Content, b.CreatedAt, b.UserID, b.Slug}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&b.ID)

	if err != nil {
		return err
	}
	return nil
}

func (m BlogModel) List() ([]Blog, error) {
	query := `
	SELECT id, title, content, created_at, slug FROM blogs`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var blogs []Blog
	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var b Blog
		err := rows.Scan(&b.ID, &b.Title, &b.Content, &b.CreatedAt, &b.Slug)
		if err != nil {
			return nil, err
		}
		blogs = append(blogs, b)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return blogs, nil
}

func (m BlogModel) Get(id int) (*Blog, error) {
	query := `SELECT id, title, content, created_at, user_id FROM blogs WHERE id = $1`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	var blog Blog
	err := m.DB.QueryRowContext(ctx, query, id).Scan(&blog.ID, &blog.Title, &blog.Content, &blog.CreatedAt, &blog.UserID)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrNoRows
		default:
			return nil, err
		}
	}
	return &blog, nil
}
