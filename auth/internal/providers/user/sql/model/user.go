package model

import "time"

type User struct {
	ID        int64     `db:"id"`
	Email     string    `db:"email"`
	PwHash    string    `db:"pwhash"`
	CreatedAt time.Time `db:"created_at"`
	Disabled  bool      `db:"disabled"`
}
