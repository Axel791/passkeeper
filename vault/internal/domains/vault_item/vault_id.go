package domains

import "errors"

var ErrInvalidVaultID = errors.New("invalid vault ID")

type VaultID struct {
	value int64
}

func NewVaultID(value int64) (VaultID, error) {
	if value <= 0 {
		return VaultID{}, ErrInvalidVaultID
	}
	return VaultID{value: value}, nil
}
