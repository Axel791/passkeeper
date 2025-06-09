package domains

import "time"

// User представляет доменную сущность «пользователь»
// Пользователь хранит данные для аутентификации и состояния учётной записи
type User struct {
	// id — уникальный идентификатор пользователя
	id UserID
	// email — адрес электронной почты пользователя (логин)
	email string
	// pwHash — хеш пароля для проверки подлинности
	pwHash string
	// createdAt — момент создания учётной записи
	createdAt time.Time
	// disabled — флаг, указывающий, заблокирована ли учётная запись
	disabled bool
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

func NewCreateUser(
	email string,
	pwHash string,
) User {
	return User{
		email:  email,
		pwHash: pwHash,
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
