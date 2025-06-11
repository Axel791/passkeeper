package config

import (
	groupprovider "github.com/Axel791/auth/internal/usecases/group/providers"
	userprovider "github.com/Axel791/auth/internal/usecases/user/providers"
	"github.com/jmoiron/sqlx"

	groupsql "github.com/Axel791/auth/internal/providers/group/sql"
	usersql "github.com/Axel791/auth/internal/providers/user/sql"
)

func newGroupProvider(db *sqlx.DB) groupprovider.GroupRepository {
	group := groupsql.NewGroupRepository(db)
	return group
}

func newUserProvider(db *sqlx.DB) userprovider.UserRepository {
	user := usersql.NewSqlUserRepository(db)
	return user
}

type Providers struct {
	GroupProvider groupprovider.GroupRepository
	UserProvider  userprovider.UserRepository
}

func NewProviders(db *sqlx.DB) Providers {
	return Providers{
		GroupProvider: newGroupProvider(db),
		UserProvider:  newUserProvider(db),
	}
}
