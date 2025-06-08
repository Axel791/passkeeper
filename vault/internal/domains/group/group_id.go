package domains

import "errors"

var (
	ErrInvalidGroupID = errors.New("invalid group ID")
)

type GroupID struct {
	value int64
}

func NewGroupID(value int64) (GroupID, error) {
	if value <= 0 {
		return GroupID{}, ErrInvalidGroupID
	}
	return GroupID{value: value}, nil
}

func (g GroupID) ToInt64() int64 {
	return g.value
}
