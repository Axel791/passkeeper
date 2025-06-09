package services

import (
	"context"

	userdomain "github.com/Axel791/vault/internal/domains/user"
)

type Services interface {
	AuthenticateToken(ctx context.Context, token string) (userdomain.UserID, error)
}
