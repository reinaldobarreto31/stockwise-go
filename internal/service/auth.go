package service

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/reinaldobarreto31/stockwise-go/internal/middleware"
	"github.com/reinaldobarreto31/stockwise-go/internal/model"
	"github.com/reinaldobarreto31/stockwise-go/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrEmailTaken     = errors.New("email already registered")
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrUserNotFound   = errors.New("user not found")
)

// AuthService handles registration, login, and token generation.
type AuthService struct {
	userRepo *repository.UserRepository
}

// NewAuthService creates a new AuthService.
func NewAuthService(userRepo *repository.UserRepository) *AuthService {
	return &AuthService{userRepo: userRepo}
}

// Register creates a new user account and returns a signed JWT.
func (s *AuthService) Register(input model.RegisterInput) (*model.LoginResponse, error) {
	existing, err := s.userRepo.GetByEmail(input.Email)
	if err != nil {
		return nil, fmt.Errorf("checking existing user: %w", err)
	}
	if existing != nil {
		return nil, ErrEmailTaken
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("hashing password: %w", err)
	}

	user, err := s.userRepo.Create(input.Name, input.Email, string(hash), "operator")
	if err != nil {
		return nil, fmt.Errorf("creating user: %w", err)
	}

	token, err := generateToken(user)
	if err != nil {
		return nil, err
	}

	return &model.LoginResponse{Token: token, User: *user}, nil
}

// Login validates credentials and returns a signed JWT.
func (s *AuthService) Login(input model.LoginInput) (*model.LoginResponse, error) {
	user, err := s.userRepo.GetByEmail(input.Email)
	if err != nil {
		return nil, fmt.Errorf("fetching user: %w", err)
	}
	if user == nil {
		return nil, ErrInvalidCredentials
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.Password)); err != nil {
		return nil, ErrInvalidCredentials
	}

	token, err := generateToken(user)
	if err != nil {
		return nil, err
	}

	return &model.LoginResponse{Token: token, User: *user}, nil
}

// generateToken creates a signed JWT for the given user.
func generateToken(user *model.User) (string, error) {
	secret := []byte(os.Getenv("JWT_SECRET"))

	expirationHours := 24
	if h, err := strconv.Atoi(os.Getenv("JWT_EXPIRATION_HOURS")); err == nil && h > 0 {
		expirationHours = h
	}

	claims := &middleware.Claims{
		UserID: user.ID,
		Role:   user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   strconv.FormatInt(user.ID, 10),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(expirationHours) * time.Hour)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString(secret)
	if err != nil {
		return "", fmt.Errorf("signing token: %w", err)
	}
	return signed, nil
}
