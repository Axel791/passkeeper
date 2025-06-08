package providers

import (
	"context"

	groupdomain "github.com/Axel791/auth/internal/domains/group"
	userdomain "github.com/Axel791/auth/internal/domains/user"
)

type GroupRepository interface {
	GetUserGroups(ctx context.Context, userID userdomain.UserID) ([]groupdomain.Group, error)
}
