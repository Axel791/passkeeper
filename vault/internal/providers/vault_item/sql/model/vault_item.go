package model

import "time"

type VaultItem struct {
	ID            int64     `db:"id"`
	DataType      int64     `db:"data_type"`
	EncryptedBlob []byte    `db:"encrypted_blob"`
	GroupID       int64     `db:"group_id"`
	CreatedAt     time.Time `db:"created_at"`
	UpdatedAt     time.Time `db:"updated_at"`
}
