package ports

import "context"

type UserAuthRepository interface {
	GetUserCredentialsByEmail(ctx context.Context, email string) (userID string, passwordHash string, active bool, err error)
	GetRolesByUserID(ctx context.Context, userID string) ([]string, error)
}