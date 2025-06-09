package providers

import (
	"context"

	vaultdomain "github.com/Axel791/vault/internal/domains/vault_item"
)

type VaultItemRepository interface {
	CreateVault(ctx context.Context, vaultItem vaultdomain.VaultItem) error
	UpdateVaultItem(ctx context.Context, vaultItem vaultdomain.VaultItem) error
	GetListVaultItems(ctx context.Context, groupID vaultdomain.GroupID) ([]vaultdomain.VaultItem, error)
}
