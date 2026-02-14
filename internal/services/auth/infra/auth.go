package infra

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"stock-app/internal/services/auth/domain"

	"golang.org/x/crypto/bcrypt"
)

type UserRepo struct {
	DB *sql.DB
}

func (u *UserRepo) Signup(ctx context.Context, user *domain.User) error {
	_, err := u.DB.Exec(`
		INSERT INTO users (id, name, email, password, balance, created_at)
		VALUES ($1, $2, $3, $4, $5, $6);
	 `, user.ID, user.Name, user.Email, user.Password, user.Balance, user.CreatedAt)

	if err != nil {
		return fmt.Errorf("couldnt insert new user: %w", err)
	}

	return nil
}

func (u *UserRepo) Login(ctx context.Context, email, password string) error {
	var storedPassowrd string

	row := u.DB.QueryRow(`SELECT password FROM users WHERE email=$1`, email)

	err := row.Scan(
		&storedPassowrd,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("invalid credentials")
		}
		return fmt.Errorf("db error: %w", err)
	}

	err = bcrypt.CompareHashAndPassword(
		[]byte(password),
		[]byte(storedPassowrd),
	)

	if err != nil {
		return fmt.Errorf("invalid credentials")
	}

	return nil
}

func (u *UserRepo) FindUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	var user domain.User
	row := u.DB.QueryRow(`
		SELECT id, name, email, balance, created_at FROM users WHERE email=$1;
	`, email)

	err := row.Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Balance,
		&user.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("row not found")
		}

		return nil, fmt.Errorf("error scanning: %s", err)
	}

	return &user, nil
}

func (u *UserRepo) UpdateBalance(ctx context.Context, id string, amount int) (int, error) {

	var newBalance int

	err := u.DB.QueryRowContext(ctx, `
		UPDATE users
		SET balance = balance + $1
		WHERE id = $2
		AND balance + $1 >= 0
		RETURNING balance
	`, amount, id).Scan(&newBalance)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, fmt.Errorf("insufficient balance or user not found")
		}
		return 0, err
	}

	return newBalance, nil
}

func (u *UserRepo) DeleteByID(ctx context.Context, id string) error {
	res, err := u.DB.ExecContext(ctx, `
		DELETE FROM users WHERE id = $1
	`, id)

	if err != nil {
		return fmt.Errorf("failder to delete user: %w", err)
	}
	// need to check if rows were affected cuz if user is not found then postge doesnt throw 
	rowsAffected, err := res.RowsAffected()

	if err != nil {
		return fmt.Errorf("failed to check result: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}
