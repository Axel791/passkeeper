package group

import (
	"context"

	userdomain "github.com/Axel791/auth/internal/domains/user"
	"github.com/Axel791/auth/internal/usecases/group/dto"
)

type GetGroupByUserID interface {
	Execute(ctx context.Context, userID userdomain.UserID) ([]dto.Group, error)
}
