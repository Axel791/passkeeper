package sql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/Axel791/auth/internal/providers/user/sql/model"
	"github.com/Axel791/auth/internal/usecases/user/providers"
	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"

	userdomain "github.com/Axel791/auth/internal/domains/user"
)

type UserRepository struct {
	db *sqlx.DB
}

func NewSqlUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

// CreateUser - создание пользователя
func (r *UserRepository) CreateUser(ctx context.Context, user userdomain.User) error {
	query, args, err := sq.StatementBuilder.
		PlaceholderFormat(sq.Dollar).
		Insert("users").
		Columns("email", "pwhash").
		Values(
			user.Email(),
			user.PwHash(),
		).
		Suffix("RETURNING id").
		ToSql()

	if err != nil {
		return fmt.Errorf("build insert user query: %w", err)
	}

	var dummy int64
	if err = r.db.QueryRowxContext(ctx, query, args...).Scan(&dummy); err != nil {
		return fmt.Errorf("insert user: %w", err)
	}

	return nil
}

// GetUserByEmail - получение пользователя по его Email
func (r *UserRepository) GetUserByEmail(ctx context.Context, email string) (userdomain.User, error) {
	query, args, err := sq.StatementBuilder.
		PlaceholderFormat(sq.Dollar).
		Select("id", "email", "pwhash", "created_at", "disabled").
		From("users").
		Where(sq.Eq{"email": email}).
		Limit(1).
		ToSql()
	if err != nil {
		return userdomain.User{}, fmt.Errorf("build select user by email query: %w", err)
	}

	var userDB model.User
	if err = r.db.GetContext(ctx, &userDB, query, args...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return userdomain.User{}, providers.ErrUserNotFound
		}
		return userdomain.User{}, fmt.Errorf("query user by email: %w", err)
	}

	return toDomainUser(userDB)
}

func toDomainUser(userDB model.User) (userdomain.User, error) {
	userID, err := userdomain.NewUserID(userDB.ID)
	if err != nil {
		return userdomain.User{}, fmt.Errorf("invalid user ID: %w", err)
	}
	user := userdomain.NewUser(
		userID,
		userDB.Email,
		userDB.PwHash,
		userDB.CreatedAt,
		userDB.Disabled,
	)
	return user, nil
}
