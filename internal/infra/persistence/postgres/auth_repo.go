package postgres

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"

	"wimed/internal/application/ports"
)

type AuthRepository struct{}

func (r AuthRepository) GetUserCredentialsByEmail(ctx context.Context, email string) (string, string, bool, error) {
	const q = `
SELECT id, password_hash, active
FROM users
WHERE email = $1
`
	var (
		id           string
		passwordHash string
		active       bool
	)
	err := dbQueryRow(ctx, q, email).Scan(&id, &passwordHash, &active)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", "", false, errors.New("invalid credentials")
		}
		return "", "", false, err
	}
	return id, passwordHash, active, nil
}

func (r AuthRepository) GetRolesByUserID(ctx context.Context, userID string) ([]string, error) {
	roles := make([]string, 0, 2)

	// DOCTOR?
	{
		const q = `SELECT 1 FROM doctors WHERE user_id = $1 LIMIT 1`
		var one int
		err := dbQueryRow(ctx, q, userID).Scan(&one)
		if err == nil {
			roles = append(roles, "DOCTOR")
		} else if !errors.Is(err, pgx.ErrNoRows) {
			return nil, err
		}
	}

	// PATIENT?
	{
		const q = `SELECT 1 FROM patients WHERE user_id = $1 LIMIT 1`
		var one int
		err := dbQueryRow(ctx, q, userID).Scan(&one)
		if err == nil {
			roles = append(roles, "PATIENT")
		} else if !errors.Is(err, pgx.ErrNoRows) {
			return nil, err
		}
	}

	return roles, nil
}

var _ ports.UserAuthRepository = AuthRepository{}
