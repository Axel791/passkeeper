package model

type VaultItemUpdateRequest struct {
	ID            int64  `json:"id"`
	GroupID       int64  `json:"group_id"`
	DataType      int64  `json:"data_type"`
	EncryptedBlob []byte `json:"encrypted_blob"`
}
