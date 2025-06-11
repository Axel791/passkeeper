package model

import "time"

type VaultItemResponse struct {
	ID            int64     `json:"id"`
	GroupID       int64     `json:"group_id"`
	DataType      int64     `json:"data_type"`
	EncryptedBlob []byte    `json:"encrypted_blob"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
