package services

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"golang-postgres-docker/models"
	"golang-postgres-docker/repository"
	"golang-postgres-docker/utils"

	"github.com/jackc/pgx/v5"
)

type AuthService struct {
	userRepo  *repository.UserRepository
	jwtSecret string
}

func NewAuthService(
	userRepo *repository.UserRepository,
	jwtSecret string,
) *AuthService {
	return &AuthService{
		userRepo:  userRepo,
		jwtSecret: jwtSecret,
	}
}

func (s *AuthService) Register(
	ctx context.Context,
	name string,
	email string,
	password string,
) (*models.User, error) {
	name = strings.TrimSpace(name)
	email = strings.TrimSpace(strings.ToLower(email))

	if name == "" {
		return nil, fmt.Errorf("name is required")
	}

	if email == "" {
		return nil, fmt.Errorf("email is required")
	}

	if password == "" {
		return nil, fmt.Errorf("password is required")
	}

	if len(password) < 6 {
		return nil, fmt.Errorf("password must be at least 6 characters")
	}

	existingUser, err := s.userRepo.FindUserByEmail(ctx, email)

	if err == nil && existingUser != nil {
		return nil, fmt.Errorf("email already exists")
	}

	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return nil, fmt.Errorf("failed to check email: %w", err)
	}

	passwordHash, err := utils.HashPassword(password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	user := &models.User{
		Name:         name,
		Email:        email,
		PasswordHash: passwordHash,
	}

	if err := s.userRepo.CreateUser(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return user, nil
}

type LoginResult struct {
	User  *models.User
	Token string
}

func (s *AuthService) Login(
	ctx context.Context,
	email string,
	password string,
) (*LoginResult, error) {
	email = strings.TrimSpace(strings.ToLower(email))

	if email == "" {
		return nil, fmt.Errorf("email is required")
	}

	if password == "" {
		return nil, fmt.Errorf("password is required")
	}

	user, err := s.userRepo.FindUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("invalid email or password")
		}

		return nil, fmt.Errorf("failed to find user: %w", err)
	}

	if err := utils.CheckPassword(password, user.PasswordHash); err != nil {
		return nil, fmt.Errorf("invalid email or password")
	}

	token, err := utils.GenerateToken(user.ID, s.jwtSecret)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	return &LoginResult{
		User:  user,
		Token: token,
	}, nil
}
