package scenarios

import (
	"context"
	"fmt"

	"github.com/Axel791/appkit"
	groupdomain "github.com/Axel791/vault/internal/domains/group"
	userdomain "github.com/Axel791/vault/internal/domains/user"
	"github.com/Axel791/vault/internal/usecases/vault_items/dto"
	"github.com/Axel791/vault/internal/usecases/vault_items/providers"
	log "github.com/sirupsen/logrus"
)

// GetVaultItems возвращает все записи сейфа для указанной группы,
// попутно убеждаясь, что userID принадлежит этой группе.
type GetVaultItems struct {
	logger    *log.Logger
	groupRepo providers.AuthGroupRetriever
	vaultRepo providers.VaultItemRepository
}

func NewGetVaultItems(
	logger *log.Logger,
	groupRepo providers.AuthGroupRetriever,
	vaultRepo providers.VaultItemRepository,
) *GetVaultItems {
	return &GetVaultItems{
		logger:    logger,
		groupRepo: groupRepo,
		vaultRepo: vaultRepo,
	}
}

func (s *GetVaultItems) Execute(
	ctx context.Context,
	userID userdomain.UserID,
	groupID groupdomain.GroupID,
) ([]dto.VaultOutput, error) {
	groups, err := s.groupRepo.GetUserGroups(ctx, userID)
	if err != nil {
		s.logger.WithError(err).Error("failed to retrieve groups")

		return nil, fmt.Errorf("fetching user groups: %w", err)
	}

	isMember := false
	for _, g := range groups {
		if g.ID().ToInt64() == groupID.ToInt64() {
			isMember = true
			break
		}
	}
	if !isMember {
		return nil, appkit.ForbiddenError("user group not found")
	}

	items, err := s.vaultRepo.GetListVaultItems(ctx, groupID)
	if err != nil {
		s.logger.WithError(err).Error("failed to retrieve items")

		return nil, fmt.Errorf("loading vault items: %w", err)
	}

	out := make([]dto.VaultOutput, 0, len(items))
	for _, it := range items {
		out = append(out, dto.VaultOutput{
			ID:            it.ID().ToInt64(),
			DataType:      it.DataType(),
			EncryptedBlob: it.EncryptedBlob(),
			GroupID:       it.GroupID().ToInt64(),
			CreatedAt:     it.CreatedAt(),
			UpdatedAt:     it.UpdatedAt(),
		})
	}

	s.logger.Printf(
		"[vault] %d items fetched for group %s by user %s",
		len(out), groupID, userID,
	)

	return out, nil
}
