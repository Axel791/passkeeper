package domains

import "errors"

var (
	ErrInvalidUserID = errors.New("invalid user ID")
)

type UserID struct {
	value int64
}

func NewUserID(value int64) (UserID, error) {
	if value <= 0 {
		return UserID{}, ErrInvalidUserID
	}
	return UserID{value: value}, nil
}

func (u UserID) ToInt64() int64 {
	return u.value
}
