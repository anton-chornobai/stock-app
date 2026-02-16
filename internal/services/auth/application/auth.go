package application

import (
	"context"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"os"
	"stock-app/internal/services/auth/domain"
	"time"
)

var (
	ErrSignatureNotConfigured = errors.New("jwt secret not configured")
)

type AuthService struct {
	repo domain.UserRepository
}

type SignupRequest struct {
	Email    string
	Password string
}

func NewAuthService(repo domain.UserRepository) *AuthService {
	return &AuthService{
		repo: repo,
	}
}

func (a *AuthService) Signup(ctx context.Context, req SignupRequest) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)

	if err != nil {
		return "", fmt.Errorf("couldnt generate hash for password %w", err)
	}

	user, err := domain.NewUser(req.Email, string(hashedPassword))

	if err != nil {
		return "", fmt.Errorf("failed to create new user: %w", err)
	}

	err = a.repo.Signup(ctx, user)

	if err != nil {
		return "", fmt.Errorf("failed to signup user in repo: %w", err)
	}

	// assigning user hashed password
	claims := jwt.MapClaims{
		"sub":  user.ID,
		"exp":  time.Now().Add(time.Minute * 15).Unix(),
		"role": user.Role,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	secret := os.Getenv("SECRET")
	if secret == "" {
		return "", ErrSignatureNotConfigured
	}
	signedToken, err := token.SignedString([]byte(secret))

	if err != nil {
		return "", fmt.Errorf("failed to sign jwt: %w", err)
	}

	return signedToken, nil
}

func (a *AuthService) Login(ctx context.Context) {

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

func (a *AuthService) DeleteByID(ctx context.Context) {}

func (a *AuthService) ValidateToken() {

}
func (a *AuthService) UpdateProfile() {

}
func (a *AuthService) CheckAdmin() {

}
