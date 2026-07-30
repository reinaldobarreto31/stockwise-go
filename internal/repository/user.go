package repository

import (
	"database/sql"

	"github.com/reinaldobarreto31/stockwise-go/internal/model"
)

// UserRepository handles user persistence.
type UserRepository struct {
	db *sql.DB
}

// NewUserRepository creates a new UserRepository.
func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

// GetByEmail returns a user by email, or nil if not found.
func (r *UserRepository) GetByEmail(email string) (*model.User, error) {
	var u model.User
	err := r.db.QueryRow(`
		SELECT id, name, email, password_hash, role, created_at
		FROM users WHERE email = $1
	`, email).Scan(&u.ID, &u.Name, &u.Email, &u.PasswordHash, &u.Role, &u.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &u, nil
}

// GetByID returns a user by ID, or nil if not found.
func (r *UserRepository) GetByID(id int64) (*model.User, error) {
	var u model.User
	err := r.db.QueryRow(`
		SELECT id, name, email, password_hash, role, created_at
		FROM users WHERE id = $1
	`, id).Scan(&u.ID, &u.Name, &u.Email, &u.PasswordHash, &u.Role, &u.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &u, nil
}

// Create inserts a new user and returns the created record.
func (r *UserRepository) Create(name, email, passwordHash, role string) (*model.User, error) {
	var u model.User
	err := r.db.QueryRow(`
		INSERT INTO users (name, email, password_hash, role)
		VALUES ($1, $2, $3, $4)
		RETURNING id, name, email, password_hash, role, created_at
	`, name, email, passwordHash, role).Scan(
		&u.ID, &u.Name, &u.Email, &u.PasswordHash, &u.Role, &u.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &u, nil
}
