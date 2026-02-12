package infra

import (
	"database/sql"

)

type UserRepo struct {
	DB *sql.DB
}

func (u *UserRepo) Create() {
	
}

func (u *UserRepo) FindUserByEmail(s string) error {
	return nil
}

