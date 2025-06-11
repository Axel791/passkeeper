package providers

import (
	"context"
	"errors"

	userdomain "github.com/Axel791/auth/internal/domains/user"
)

var ErrUserNotFound = errors.New("user not found")

type UserRepository interface {
	CreateUser(ctx context.Context, user userdomain.User) error
	GetUserByEmail(ctx context.Context, email string) (userdomain.User, error)
}
