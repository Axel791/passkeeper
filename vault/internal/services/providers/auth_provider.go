package providers

import (
	"context"

	userdomain "github.com/Axel791/vault/internal/domains/user"
)

type AuthValidator interface {
	ValidateToken(ctx context.Context, token string) (userdomain.UserID, error)
}
