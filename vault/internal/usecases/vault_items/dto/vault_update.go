package dto

import (
	groupdomain "github.com/Axel791/vault/internal/domains/group"
	vaultdomain "github.com/Axel791/vault/internal/domains/vault_item"
)

type VaultItemUpdateInput struct {
	VaultID       vaultdomain.VaultID
	GroupID       groupdomain.GroupID
	DataType      int64
	EncryptedBlob []byte
}
