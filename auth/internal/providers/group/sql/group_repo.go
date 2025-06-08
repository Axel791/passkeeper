package sql

import (
	"context"
	"fmt"
	groupdomain "github.com/Axel791/auth/internal/domains/group"
	userdomain "github.com/Axel791/auth/internal/domains/user"
	"github.com/Axel791/auth/internal/providers/group/sql/model"
	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
)

type GroupRepository struct {
	db *sqlx.DB
}

func NewGroupRepository(db *sqlx.DB) *GroupRepository {
	return &GroupRepository{
		db: db,
	}
}

// GetUserGroups возвращает все группы, в которых состоит пользователь.
func (r *GroupRepository) GetUserGroups(
	ctx context.Context,
	userID userdomain.UserID,
) ([]groupdomain.Group, error) {
	builder := sq.
		StatementBuilder.
		PlaceholderFormat(sq.Dollar).
		Select(
			"g.id",
			"g.name",
			"g.description",
			"g.created_at",
		).
		From("groups g").
		Join("user_groups ug ON ug.group_id = g.id").
		Where(sq.Eq{"ug.user_id": userID})

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("build user groups sql: %w", err)
	}

	var groups []model.Group
	if err = r.db.SelectContext(ctx, &groups, query, args...); err != nil {
		return nil, fmt.Errorf("select user groups: %w", err)
	}

	return toGroupDomain(groups)
}

func toGroupDomain(groups []model.Group) ([]groupdomain.Group, error) {
	var result []groupdomain.Group
	for _, group := range groups {
		groupID, err := groupdomain.NewGroupID(group.ID)
		if err != nil {
			return nil, fmt.Errorf("error generating group id: %w", err)
		}
		groupDomain := groupdomain.NewGroup(
			groupID,
			group.Name,
			group.Description,
			group.CreatedAt,
		)
		result = append(result, groupDomain)
	}
	return result, nil
}
