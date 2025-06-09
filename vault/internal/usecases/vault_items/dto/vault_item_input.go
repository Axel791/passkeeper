package dto

import groupdomain "github.com/Axel791/vault/internal/domains/group"

type VaultItemInput struct {
	GroupID       groupdomain.GroupID
	DataType      int64
	EncryptedBlob []byte
}
