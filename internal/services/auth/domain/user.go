package domain

import (
	"fmt"

	"github.com/google/uuid"
)

type User struct {
	ID      string
	Email   string
	Name    string
	Password string
	Balance int	
	CreatedAt int
}

func NewUser(email, name string, balance int) (*User, error) {
	if email == "" {
		return nil, fmt.Errorf("email cannot be empty")
	}
	if name == "" {
		return nil, fmt.Errorf("name cannot be empty")
	}
	if balance < 0 {
		return nil, fmt.Errorf("balance cannot be negative")
	}

	return &User{
		ID:      uuid.NewString(),
		Email:   email,
		Balance: balance,
		Name:    name,
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
