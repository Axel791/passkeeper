package user

import (
	"context"

	"github.com/Axel791/auth/internal/usecases/user/dto"
)

type Validate interface {
	Execute(ctx context.Context, token string) (dto.User, error)
}
