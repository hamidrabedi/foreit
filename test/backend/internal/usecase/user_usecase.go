package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/etemademan/backend/internal/domain"
	"github.com/etemademan/backend/pkg/utils"
	"github.com/google/uuid"
)

type userUsecase struct {
	userRepo       domain.UserRepository
	jwtSecret      string
	contextTimeout time.Duration
}

func NewUserUsecase(userRepo domain.UserRepository, jwtSecret string, timeout time.Duration) domain.UserUsecase {
	return &userUsecase{
		userRepo:       userRepo,
		jwtSecret:      jwtSecret,
		contextTimeout: timeout,
	}
}

func (u *userUsecase) Register(ctx context.Context, name, email, password string) (*domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	existingUser, _ := u.userRepo.GetByEmail(ctx, email)
	if existingUser != nil {
		return nil, errors.New("user already exists")
	}

	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return nil, err
	}

	user := &domain.User{
		ID:           uuid.New(),
		Name:         name,
		Email:        email,
		PasswordHash: hashedPassword,
		CreatedAt:    time.Now(),
	}

	err = u.userRepo.Create(ctx, user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (u *userUsecase) Login(ctx context.Context, email, password string) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	user, err := u.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	if !utils.CheckPasswordHash(password, user.PasswordHash) {
		return "", errors.New("invalid credentials")
	}

	token, err := utils.GenerateJWT(user.ID, u.jwtSecret)
	if err != nil {
		return "", err
	}

	return token, nil
}
