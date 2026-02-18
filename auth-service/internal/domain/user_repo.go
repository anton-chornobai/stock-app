package domain

import "context"

type UserRepository interface {
	Login(ctx context.Context, email, password string) (*User, error)
	Signup(ctx context.Context, user *User) error
	FindUserByEmail(ctx context.Context, email string) (*User, error)
	UpdateBalance(ctx context.Context, id string, amount int) (int, error)
	DeleteByID(ctx context.Context, id string) error
}
