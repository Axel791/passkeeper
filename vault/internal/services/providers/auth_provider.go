package providers

import "context"

type AuthValidator interface {
	ValidateToken(ctx context.Context, token string) error
}
