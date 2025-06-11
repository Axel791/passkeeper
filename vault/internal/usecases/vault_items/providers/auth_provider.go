package providers

import (
	"context"

	groupdomain "github.com/Axel791/vault/internal/domains/group"
	userdomain "github.com/Axel791/vault/internal/domains/user"
)

type AuthGroupRetriever interface {
	GetUserGroups(ctx context.Context, id userdomain.UserID) ([]groupdomain.Group, error)
}
