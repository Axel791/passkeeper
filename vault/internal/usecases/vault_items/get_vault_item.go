package vault_items

import (
	"context"

	groupdomain "github.com/Axel791/vault/internal/domains/group"
	userdomain "github.com/Axel791/vault/internal/domains/user"
	"github.com/Axel791/vault/internal/usecases/vault_items/dto"
)

type GetVaultItems interface {
	Execute(ctx context.Context, userID userdomain.UserID, groupID groupdomain.GroupID) ([]dto.VaultOutput, error)
}
