package domain

import "errors"

var (
	ErrUserNotFound = errors.New("user not found")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrShortPassword = errors.New("password must be at least 8 characters")
	ErrInvalidEmail  = errors.New("invalid email")

)