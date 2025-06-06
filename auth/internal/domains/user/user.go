package domains

import "time"

type User struct {
	id        UserID
	email     string
	pwHash    string
	createdAt time.Time
	disabled  bool
}

func NewUser(
	id UserID,
	email string,
	pwHash string,
	createdAt time.Time,
	disabled bool,
) User {
	return User{
		id:        id,
		email:     email,
		pwHash:    pwHash,
		createdAt: createdAt,
		disabled:  disabled,
	}
}

func (u User) ID() UserID {
	return u.id
}

func (u User) Email() string {
	return u.email
}

func (u User) PwHash() string {
	return u.pwHash
}

func (u User) CreatedAt() time.Time {
	return u.createdAt
}

func (u User) Disabled() bool {
	return u.disabled
}
