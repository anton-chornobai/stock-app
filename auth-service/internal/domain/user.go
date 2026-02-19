package domain

import (
	"strings"

	"github.com/google/uuid"
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
		ID:       uuid.NewString(),
		Email:    email,
		Password: password,
		Role:     "user",
		Balance:  0,
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
