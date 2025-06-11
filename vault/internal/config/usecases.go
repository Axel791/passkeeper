package config

import (
	vaultusecases "github.com/Axel791/vault/internal/usecases/vault_items"
	vaultscenarios "github.com/Axel791/vault/internal/usecases/vault_items/scenarios"
	log "github.com/sirupsen/logrus"
)

func newCreateVaultItem(
	logger *log.Logger,
	providers Providers,
) vaultusecases.CreateVaultItem {
	createVault := vaultscenarios.NewCreateVaultItem(
		logger,
		providers.GroupProvider,
		providers.VaultRepo,
	)
	return createVault
}

func newUpdateVaultItem(
	logger *log.Logger,
	providers Providers,
) vaultusecases.UpdateVaultItem {
	updateVault := vaultscenarios.NewUpdateVaultItem(
		logger,
		providers.GroupProvider,
		providers.VaultRepo,
	)
	return updateVault
}

func newGetVaultItems(
	logger *log.Logger,
	providers Providers,
) vaultusecases.GetVaultItems {
	getVaultItems := vaultscenarios.NewGetVaultItems(
		logger,
		providers.GroupProvider,
		providers.VaultRepo,
	)
	return getVaultItems
}

type UseCases struct {
	CreateVaultItem vaultusecases.CreateVaultItem
	GetVaultItems   vaultusecases.GetVaultItems
	UpdateVaultItem vaultusecases.UpdateVaultItem
}

func NewUseCases(
	logger *log.Logger,
	providers Providers,
) UseCases {
	return UseCases{
		CreateVaultItem: newCreateVaultItem(logger, providers),
		GetVaultItems:   newGetVaultItems(logger, providers),
		UpdateVaultItem: newUpdateVaultItem(logger, providers),
	}
}
