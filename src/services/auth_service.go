package services

import (
	"bank-service/src/models"
	"bank-service/src/repositories"
	"errors"
	"time"
	
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"github.com/sirupsen/logrus"
)

type AuthService struct {
	userRepo  *repositories.UserRepository
	jwtSecret string
	logger    *logrus.Logger
}

func NewAuthService(
	userRepo *repositories.UserRepository, 
	jwtSecret string,
	logger *logrus.Logger,
) *AuthService {
	return &AuthService{
		userRepo:  userRepo,
		jwtSecret: jwtSecret,
		logger:    logger,
	}
}

func (s *AuthService) Register(user *models.User) error {
	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(user.PasswordHash), bcrypt.DefaultCost)
	if err != nil {
		s.logger.WithError(err).Error("Failed to hash password")
		return err
	}
	
	user.PasswordHash = string(hashedPassword)
	return s.userRepo.Create(user)
}

func (s *AuthService) Login(email, password string) (string, error) {
	user, err := s.userRepo.GetByEmail(email)
	if err != nil {
		s.logger.WithError(err).Error("Login failed - user not found")
		return "", errors.New("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword(
		[]byte(user.PasswordHash), []byte(password)); err != nil {
		s.logger.WithError(err).Error("Login failed - invalid password")
		return "", errors.New("invalid credentials")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(24 * time.Hour).Unix(),
	})

	tokenString, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil {
		s.logger.WithError(err).Error("Failed to generate JWT")
		return "", err
	}

	return tokenString, nil
}