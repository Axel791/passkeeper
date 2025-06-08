package sql

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/Axel791/vault/internal/providers/vault_item/sql/model"
	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"

	vaultdomain "github.com/Axel791/vault/internal/domains/vault_item"
)

type VaultItemRepo struct {
	db *sqlx.DB
}

func NewVaultItemRepo(db *sqlx.DB) *VaultItemRepo {
	return &VaultItemRepo{
		db: db,
	}
}

func (r *VaultItemRepo) CreateVault(ctx context.Context, vaultItem vaultdomain.VaultItem) error {
	qb := sq.Insert("vault_items").
		Columns(
			"data_type",
			"encrypted_blob",
			"group_id",
			"created_at",
			"updated_at",
		).
		Values(
			vaultItem.DataType(),
			vaultItem.EncryptedBlob(),
			vaultItem.GroupID().ToInt64(),
			vaultItem.CreatedAt(),
			vaultItem.UpdatedAt(),
		).
		PlaceholderFormat(sq.Dollar)

	query, args, err := qb.ToSql()
	if err != nil {
		return fmt.Errorf("build insert vault_items sql: %w", err)
	}

	if _, err := r.db.ExecContext(ctx, query, args...); err != nil {
		return fmt.Errorf("exec insert vault_items: %w", err)
	}
	return nil
}

func (r *VaultItemRepo) UpdateVaultItem(ctx context.Context, vaultItem vaultdomain.VaultItem) error {
	qb := sq.Update("vault_items").
		Set("data_type", vaultItem.DataType()).
		Set("encrypted_blob", vaultItem.EncryptedBlob()).
		Set("group_id", vaultItem.GroupID().ToInt64()).
		Set("updated_at", vaultItem.UpdatedAt()).
		Where(sq.Eq{"id": vaultItem.ID()}).
		PlaceholderFormat(sq.Dollar)

	query, args, err := qb.ToSql()
	if err != nil {
		return fmt.Errorf("build update vault_items sql: %w", err)
	}

	res, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("exec update vault_items: %w", err)
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected: %w", err)
	}
	if rows == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (r *VaultItemRepo) GetListVaultItems(
	ctx context.Context,
	groupID vaultdomain.GroupID,
) ([]vaultdomain.VaultItem, error) {
	qb := sq.Select(
		"id",
		"data_type",
		"encrypted_blob",
		"group_id",
		"created_at",
		"updated_at",
	).
		From("vault_items").
		Where(sq.Eq{"group_id": groupID.ToInt64()}).
		OrderBy("created_at DESC").
		PlaceholderFormat(sq.Dollar)

	query, args, err := qb.ToSql()
	if err != nil {
		return nil, fmt.Errorf("build select vault_items sql: %w", err)
	}

	var vaults []model.VaultItem
	if err = r.db.SelectContext(ctx, &vaults, query, args...); err != nil {
		return nil, fmt.Errorf("select vault_items: %w", err)
	}

	return toVaultDomain(vaults)
}

func toVaultDomain(vaults []model.VaultItem) ([]vaultdomain.VaultItem, error) {
	var result []vaultdomain.VaultItem
	for _, v := range vaults {
		vaultID, err := vaultdomain.NewVaultID(v.ID)
		if err != nil {
			return nil, fmt.Errorf("error creating vaultID: %w", err)
		}

		groupID, err := vaultdomain.NewGroupID(v.GroupID)
		if err != nil {
			return nil, fmt.Errorf("error creating groupID: %w", err)
		}

		vault := vaultdomain.NewVaultItem(
			vaultID,
			v.DataType,
			groupID,
			v.EncryptedBlob,
			v.CreatedAt,
			v.UpdatedAt,
		)

		result = append(result, vault)
	}
	return result, nil
}
