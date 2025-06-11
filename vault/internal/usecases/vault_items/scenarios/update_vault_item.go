package scenarios

import (
	"context"

	"github.com/Axel791/appkit"
	userdomain "github.com/Axel791/vault/internal/domains/user"
	vaultdomain "github.com/Axel791/vault/internal/domains/vault_item"
	"github.com/Axel791/vault/internal/usecases/vault_items/dto"
	"github.com/Axel791/vault/internal/usecases/vault_items/providers"
	log "github.com/sirupsen/logrus"
)

type UpdateVaultItem struct {
	logger    *log.Logger
	groupRepo providers.AuthGroupRetriever
	vaultRepo providers.VaultItemRepository
}

func NewUpdateVaultItem(
	logger *log.Logger,
	groupRepo providers.AuthGroupRetriever,
	vaultRepo providers.VaultItemRepository,
) *UpdateVaultItem {
	return &UpdateVaultItem{
		logger:    logger,
		groupRepo: groupRepo,
		vaultRepo: vaultRepo,
	}
}

func (s *UpdateVaultItem) Execute(
	ctx context.Context,
	userID userdomain.UserID,
	data dto.VaultItemUpdateInput,
) error {
	groups, err := s.groupRepo.GetUserGroups(ctx, userID)
	if err != nil {
		s.logger.WithError(err).Error("failed to fetch groups")

		return appkit.Wrap(appkit.Unknown, "error getting user groups", err)
	}

	inGroup := false
	for _, g := range groups {
		if g.ID().ToInt64() == data.GroupID.ToInt64() {
			inGroup = true
			break
		}
	}
	if !inGroup {
		return appkit.ForbiddenError("user group not found")
	}

	item := vaultdomain.NewUpdateVaultItem(
		data.VaultID,
		data.DataType,
		data.EncryptedBlob,
		data.GroupID,
	)

	if err = s.vaultRepo.UpdateVaultItem(ctx, item); err != nil {
		s.logger.Printf("[vault] update error: %v", err)
		return appkit.Wrap(appkit.Unknown, "error updating vault item: %w", err)
	}

	s.logger.Printf("[vault] item %d updated in group %d", item.ID().ToInt64(), item.GroupID().ToInt64())
	return nil
}
