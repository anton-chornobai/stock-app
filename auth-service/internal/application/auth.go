package application

import (
	"auth-service/internal/domain"
	tokenmanager "auth-service/internal/lib/jwt"
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrSignatureNotConfigured = errors.New("jwt secret not configured")
)

type AuthService struct {
	repo   domain.UserRepository
	secret []byte
}

type SignupRequest struct {
	Email    string `json:"email"`
	Password string	`json:"password"`
}

type LoginRequest struct {
	Email    string
	Password string
}

func NewAuthService(repo domain.UserRepository, secret []byte) *AuthService {
	return &AuthService{
		repo:   repo,
		secret: secret,
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

	signedToken, err := tokenmanager.GenerateUserToken(user, a.secret)

	if err != nil {
		return "", fmt.Errorf("failed to generate token: %w", err)
	}

	return signedToken, nil
}

func (a *AuthService) Login(ctx context.Context, req LoginRequest) (string, error) {
	user, err := a.repo.Login(ctx, req.Email, req.Password)

	if err != nil {
		return "", fmt.Errorf("couldnt login user in db %w", err)
	}

	claims := jwt.MapClaims{
		"sub":  user.ID,
		"exp":  time.Now().Add(time.Minute * 15).Unix(),
		"role": user.Role,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString(a.secret)
	if err != nil {
		return "", fmt.Errorf("failed to sign JWT: %w", err)
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
