package domain

type UserRepository interface {
	// Create() (error)
	FindUserByEmail(email string) (error)
	
}