package application

import (
	"auth-service/internal/domain"
	tokenmanager "auth-service/internal/lib/jwt"
	"context"
	"errors"
	"fmt"
	"log/slog"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrSignatureNotConfigured = errors.New("jwt secret not configured")
)

type AuthService struct {
	repo   domain.UserRepository
	secret []byte
	logger *slog.Logger
}

func NewAuthService(repo domain.UserRepository, secret []byte, logger *slog.Logger) *AuthService {
	return &AuthService{
		repo:   repo,
		secret: secret,
		logger: logger,
	}
}

func (a *AuthService) Signup(ctx context.Context, email, password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		return "", err
	}

	user, err := domain.NewUser(email, string(hashedPassword))

	if err != nil {
		return "", err
	}

	err = a.repo.Signup(ctx, user)

	if err != nil {
		return "", err
	}

	signedToken, err := tokenmanager.GenerateUserToken(user, a.secret)

	if err != nil {
		return "", fmt.Errorf("failed to generate token: %w", err)
	}

	return signedToken, nil
}

func (a *AuthService) Login(ctx context.Context, email, password string) (string, error) {
	user, err := a.repo.Login(ctx, email, password)

	if err != nil {
		return "", err
	}

	signedToken, err := tokenmanager.GenerateUserToken(user, a.secret)
	if err != nil {
		return "", fmt.Errorf("failed to generate JWT: %w", err)
	}

	return signedToken, nil
}

func (a *AuthService) Logout(ctx context.Context) {

}
func (g *AuthService) FindUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	user, err := g.repo.FindUserByEmail(ctx, email)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (a *AuthService) DeleteByID(ctx context.Context, id string) error {
	err := a.repo.DeleteByID(ctx, id)

	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	return nil
}

func (a *AuthService) ValidateToken() {

}

func (a *AuthService) UpdateProfile() {

}
func (a *AuthService) CheckAdmin() {

}
