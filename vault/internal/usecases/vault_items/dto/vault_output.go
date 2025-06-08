package dto

import "time"

type VaultOutput struct {
	ID            int64
	DataType      int64
	EncryptedBlob []byte
	GroupID       int64
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
