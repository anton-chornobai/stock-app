package application

import (
	"context"
	"fmt"
	"stock-app/internal/services/auth/domain"

	"golang.org/x/crypto/bcrypt"
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

func (a *AuthService) Signup(ctx context.Context, req SignupRequest) error {

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)

	if err != nil {
		return fmt.Errorf("couldnt generate hash for password %w", err)
	}

	user, err := domain.NewUser(req.Email, string(hashedPassword))

	if err != nil {
		return fmt.Errorf("failed to create new user: %w", err)
	}
	err = a.repo.Signup(ctx, user) 

	if err != nil {
		return fmt.Errorf("failed to signup user in repo: %w", err)
	}

	return nil
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
