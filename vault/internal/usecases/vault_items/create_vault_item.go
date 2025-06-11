package vault_items

import (
	"context"

	userdomain "github.com/Axel791/vault/internal/domains/user"
	"github.com/Axel791/vault/internal/usecases/vault_items/dto"
)

type CreateVaultItem interface {
	Execute(ctx context.Context, userID userdomain.UserID, data dto.VaultItemInput) error
}
