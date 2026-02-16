package domain

import (
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
)

var (
	ErrShortPassword = errors.New("password must be at least 8 characters")
	ErrInvalidEmail  = errors.New("invalid email")
)

type User struct {
	ID        string
	Email     string
	Name      string
	Password  string
	Role      string
	Balance   int
	CreatedAt int
}

func NewUser(email, password string) (*User, error) {
	if !strings.Contains(email, "@") {
		return nil, ErrInvalidEmail
	}
	return &User{
		ID:        uuid.NewString(),
		Email:     email,
		Password:  password,
		Role: "user",
		Balance:   0,
		CreatedAt: int(time.Now().Unix()),
	}, nil
}

func (u *User) Signup() {

}

func (u *User) Login() {

}

func (u *User) Logout() {

}

func (u *User) Delete() {

}

func (u *User) Update() {

}
