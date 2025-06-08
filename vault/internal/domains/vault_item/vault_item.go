package domains

import "time"

type VaultItem struct {
	id            VaultID
	dataType      int64
	encryptedBlob []byte
	createdAt     time.Time
	updatedAt     time.Time
}

func NewVaultItem(
	id VaultID,
	dataType int64,
	encryptedBlob []byte,
	createdAt time.Time,
	updatedAt time.Time,
) VaultItem {
	return VaultItem{
		id:            id,
		dataType:      dataType,
		encryptedBlob: encryptedBlob,
		createdAt:     createdAt,
		updatedAt:     updatedAt,
	}
}

func NewCreateVaultItem(
	dataType int64,
	encryptedBlob []byte,
) VaultItem {
	return VaultItem{
		dataType:      dataType,
		encryptedBlob: encryptedBlob,
		createdAt:     time.Now(),
	}
}

func (v VaultItem) ID() VaultID {
	return v.id
}

func (v VaultItem) DataType() int64 {
	return v.dataType
}

func (v VaultItem) EncryptedBlob() []byte {
	return v.encryptedBlob
}

func (v VaultItem) CreatedAt() time.Time {
	return v.createdAt
}

func (v VaultItem) UpdatedAt() time.Time {
	return v.updatedAt
}
