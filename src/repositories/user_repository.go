package repositories

import (
	"bank-service/src/models"
	"database/sql"
	"errors"
	
	"github.com/sirupsen/logrus"
)

var (
	ErrUserExists = errors.New("user already exists")
)

type UserRepository struct {
	db     *sql.DB
	logger *logrus.Logger
}

func NewUserRepository(db *sql.DB, logger *logrus.Logger) *UserRepository {
	return &UserRepository{db: db, logger: logger}
}

func (r *UserRepository) Create(user *models.User) error {
	query := `
		INSERT INTO users (username, email, password_hash) 
		VALUES ($1, $2, $3) 
		RETURNING id, created_at`
		
	err := r.db.QueryRow(query, user.Username, user.Email, user.PasswordHash).Scan(
		&user.ID, &user.CreatedAt)
		
	if err != nil {
		r.logger.WithError(err).Error("Failed to create user")
		if err.Error() == "pq: duplicate key value violates unique constraint" {
			return ErrUserExists
		}
	}
	return err
}

func (r *UserRepository) GetByEmail(email string) (*models.User, error) {
	user := &models.User{}
	query := `
		SELECT id, username, email, password_hash, created_at 
		FROM users WHERE email = $1`
		
	err := r.db.QueryRow(query, email).Scan(
		&user.ID, &user.Username, &user.Email, 
		&user.PasswordHash, &user.CreatedAt)
		
	if err != nil {
		r.logger.WithError(err).Error("User not found")
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) GetByID(id uint) (*models.User, error) {
	user := &models.User{}
	query := `
		SELECT id, username, email, password_hash, created_at 
		FROM users WHERE id = $1`
		
	err := r.db.QueryRow(query, id).Scan(
		&user.ID, &user.Username, &user.Email, 
		&user.PasswordHash, &user.CreatedAt)
		
	if err != nil {
		if err == sql.ErrNoRows {
			r.logger.Warnf("User  with ID %d not found", id)
			return nil, nil 
		}
		r.logger.WithError(err).Error("Failed to get user by ID")
		return nil, err
	}
	return user, nil
}
