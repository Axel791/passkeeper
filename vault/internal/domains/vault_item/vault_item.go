package domains

import (
	"time"

	groupdomain "github.com/Axel791/vault/internal/domains/group"
)

// VaultItem представляет записанный защищённый элемент (секрет) в «хранилище»
// Содержит зашифрованные данные и метаданные для управления
type VaultItem struct {
	// id — уникальный идентификатор элемента хранилища
	id VaultID
	// dataType — тип данных (может указывать на формат или категорию секрета)
	dataType int64
	// groupId — ссылка на группу (владельца или категорию доступа)
	groupId groupdomain.GroupID
	// encryptedBlob — зашифрованный бинарный контент секрета
	encryptedBlob []byte
	// createdAt — момент создания записи секрета
	createdAt time.Time
	// updatedAt — момент последнего обновления (пересохранения) секрета
	updatedAt time.Time
}

func NewVaultItem(
	id VaultID,
	dataType int64,
	groupId groupdomain.GroupID,
	encryptedBlob []byte,
	createdAt time.Time,
	updatedAt time.Time,
) VaultItem {
	return VaultItem{
		id:            id,
		dataType:      dataType,
		groupId:       groupId,
		encryptedBlob: encryptedBlob,
		createdAt:     createdAt,
		updatedAt:     updatedAt,
	}
}

func NewCreateVaultItem(
	dataType int64,
	encryptedBlob []byte,
	groupId groupdomain.GroupID,
) VaultItem {
	return VaultItem{
		dataType:      dataType,
		groupId:       groupId,
		encryptedBlob: encryptedBlob,
		createdAt:     time.Now(),
	}
}

func NewUpdateVaultItem(
	id VaultID,
	dataType int64,
	encryptedBlob []byte,
	groupId groupdomain.GroupID,
) VaultItem {
	return VaultItem{
		id:            id,
		dataType:      dataType,
		groupId:       groupId,
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

func (v VaultItem) GroupID() groupdomain.GroupID {
	return v.groupId
}
