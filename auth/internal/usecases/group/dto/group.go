package dto

import "time"

type Group struct {
	ID          int64
	Name        string
	Description string
	CreatedAt   time.Time
}
