package config

import (
	"github.com/Axel791/passkeeper_grpc/pb"
	vaultgrpc "github.com/Axel791/vault/internal/providers/auth/grpc"
	vaultsql "github.com/Axel791/vault/internal/providers/vault_item/sql"
	authprovider "github.com/Axel791/vault/internal/services/providers"
	vaultprovider "github.com/Axel791/vault/internal/usecases/vault_items/providers"
	"github.com/jmoiron/sqlx"
)

type Providers struct {
	VaultRepo     vaultprovider.VaultItemRepository
	GroupProvider vaultprovider.AuthGroupRetriever
	AuthProvider  authprovider.AuthValidator
}

func newVaultRepo(db *sqlx.DB) vaultprovider.VaultItemRepository {
	vault := vaultsql.NewVaultItemRepo(db)
	return vault
}

func newGroupProvider(client pb.AuthServiceClient) vaultprovider.AuthGroupRetriever {
	group := vaultgrpc.NewAuthRepo(client)
	return group
}

func newAuthProvider(client pb.AuthServiceClient) authprovider.AuthValidator {
	auth := vaultgrpc.NewAuthRepo(client)
	return auth
}

func NewProviders(db *sqlx.DB, client pb.AuthServiceClient) Providers {
	return Providers{
		VaultRepo:     newVaultRepo(db),
		GroupProvider: newGroupProvider(client),
		AuthProvider:  newAuthProvider(client),
	}
}
