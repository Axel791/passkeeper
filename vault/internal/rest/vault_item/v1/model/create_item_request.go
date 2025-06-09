package model

type VaultItemCreateRequest struct {
	GroupID       int64  `json:"group_id"`
	DataType      int64  `json:"data_type"`
	EncryptedBlob []byte `json:"encrypted_blob"`
}
