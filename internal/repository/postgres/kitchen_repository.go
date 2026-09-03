package postgres

import (
	"context"
	"errors"
	"order-it-backend/internal/domain"
	"order-it-backend/internal/middleware"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type KitchenRepository struct {
	pool *pgxpool.Pool
}

func NewKitchenRepository(pool *pgxpool.Pool) domain.KitchenRepository {
	return &KitchenRepository{pool: pool}
}

func (k *KitchenRepository) db(ctx context.Context) domain.DBTX {
	if tx, ok := middleware.GetTx(ctx); ok {
		return tx
	}
	return k.pool
}

func (k *KitchenRepository) Create(ctx context.Context, kitchen *domain.Kitchen) error {
	query := `insert into kitchens (tenant_id, name) values ($1, $2) returning id, created_at`
	err := k.db(ctx).QueryRow(ctx, query, kitchen.TenantId, kitchen.Name).Scan(&kitchen.Id, &kitchen.CreatedAt)
	if err != nil {
		return err
	}
	return nil
}

func (k *KitchenRepository) GetById(ctx context.Context, id uuid.UUID) (*domain.Kitchen, error) {
	query := `SELECT id, tenant_id, name, created_at FROM kitchens WHERE id = $1;`
	var kitchen domain.Kitchen
	err := k.db(ctx).QueryRow(ctx, query, id).Scan(&kitchen.Id, &kitchen.TenantId, &kitchen.Name, &kitchen.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("kitchen not found")
		}
		return nil, err
	}
	return &kitchen, nil
}

func (k *KitchenRepository) GetAllByTenantId(ctx context.Context, tenantId uuid.UUID) ([]*domain.Kitchen, error) {
	query := `select id, tenant_id, name, created_at from kitchens where tenant_id = $1`
	rows, err := k.db(ctx).Query(ctx, query, tenantId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var kitchens []*domain.Kitchen
	for rows.Next() {
		var kitchen domain.Kitchen
		err := rows.Scan(&kitchen.Id, &kitchen.TenantId, &kitchen.Name, &kitchen.CreatedAt)
		if err != nil {
			return nil, err
		}
		kitchens = append(kitchens, &kitchen)
	}
	return kitchens, nil
}

func (k *KitchenRepository) GetAll(ctx context.Context) ([]*domain.Kitchen, error) {
	query := `select id, tenant_id, name, created_at from kitchens`
	rows, err := k.db(ctx).Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var kitchens []*domain.Kitchen
	for rows.Next() {
		var kitchen domain.Kitchen
		err := rows.Scan(&kitchen.Id, &kitchen.TenantId, &kitchen.Name, &kitchen.CreatedAt)
		if err != nil {
			return nil, err
		}
		kitchens = append(kitchens, &kitchen)
	}
	return kitchens, nil
}

func (k *KitchenRepository) GetByTenantIdAndName(ctx context.Context, tenantId uuid.UUID, name string) (*domain.Kitchen, error) {
	query := `SELECT id, tenant_id, name, created_at FROM kitchens WHERE tenant_id = $1 AND name = $2;`
	var kitchen domain.Kitchen
	err := k.db(ctx).QueryRow(ctx, query, tenantId, name).Scan(&kitchen.Id, &kitchen.TenantId, &kitchen.Name, &kitchen.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &kitchen, nil
}

func (k *KitchenRepository) Update(ctx context.Context, kitchen *domain.Kitchen) error {
	query := `update kitchens set name = $1 where id = $2 and tenant_id = $3`
	cmd, err := k.db(ctx).Exec(ctx, query, kitchen.Name, kitchen.Id, kitchen.TenantId)
	if err != nil {
		return err
	}
	if cmd.RowsAffected() == 0 {
		return errors.New("kitchen not found to update")
	}
	return nil
}

func (k *KitchenRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `delete from kitchens where id = $1`
	cmd, err := k.db(ctx).Exec(ctx, query, id)
	if err != nil {
		return err
	}
	if cmd.RowsAffected() == 0 {
		return errors.New("kitchen not found to delete")
	}
	return nil
}
