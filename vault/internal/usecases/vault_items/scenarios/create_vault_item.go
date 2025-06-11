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

// CreateVaultItem - создает запись в хранилище
type CreateVaultItem struct {
	logger    *log.Logger
	groupRepo providers.AuthGroupRetriever
	vaultRepo providers.VaultItemRepository
}

func NewCreateVaultItem(
	logger *log.Logger,
	groupRepo providers.AuthGroupRetriever,
	vaultRepo providers.VaultItemRepository,
) *CreateVaultItem {
	return &CreateVaultItem{
		logger:    logger,
		groupRepo: groupRepo,
		vaultRepo: vaultRepo,
	}
}

func (s *CreateVaultItem) Execute(
	ctx context.Context,
	userID userdomain.UserID,
	data dto.VaultItemInput,
) error {
	groups, err := s.groupRepo.GetUserGroups(ctx, userID)
	if err != nil {
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

	item := vaultdomain.NewCreateVaultItem(
		data.DataType,
		data.EncryptedBlob,
		data.GroupID,
	)

	if err = s.vaultRepo.CreateVault(ctx, item); err != nil {
		s.logger.WithError(err).Error("creating vault")

		return appkit.Wrap(appkit.Unknown, "error creating vault: %w", err)
	}

	s.logger.Printf("[vault] item %d created for group %d", item.ID().ToInt64(), item.GroupID().ToInt64())
	return nil
}
