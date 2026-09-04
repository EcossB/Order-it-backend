package postgres

import (
	"context"
	"errors"
	"order-it-backend/internal/domain"
	"order-it-backend/internal/middleware"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type MenuItemRepository struct {
	pool *pgxpool.Pool
}

func NewMenuItemRepository(pool *pgxpool.Pool) domain.MenuItemRepository {
	return &MenuItemRepository{pool: pool}
}

func (r *MenuItemRepository) db(ctx context.Context) domain.DBTX {
	if tx, ok := middleware.GetTx(ctx); ok {
		return tx
	}
	return r.pool
}

func (r *MenuItemRepository) Create(ctx context.Context, item *domain.MenuItem) error {
	query := `INSERT INTO menu_items (tenant_id, kitchen_id, name, description, price, is_available) 
			  VALUES ($1, $2, $3, $4, $5, $6) RETURNING id, created_at`

	err := r.db(ctx).QueryRow(ctx, query, item.TenantId, item.KitchenId, item.Name, item.Description, item.Price, item.IsAvailable).Scan(&item.Id, &item.CreatedAt)
	if err != nil {
		return err
	}
	return nil
}

func (r *MenuItemRepository) GetById(ctx context.Context, id uuid.UUID) (*domain.MenuItem, error) {
	query := `SELECT id, tenant_id, kitchen_id, name, description, price, is_available, created_at FROM menu_items WHERE id = $1`

	var item domain.MenuItem
	err := r.db(ctx).QueryRow(ctx, query, id).Scan(
		&item.Id, &item.TenantId, &item.KitchenId, &item.Name, &item.Description, &item.Price, &item.IsAvailable, &item.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *MenuItemRepository) GetAllByTenantId(ctx context.Context, tenantId uuid.UUID) ([]*domain.MenuItem, error) {
	query := `SELECT id, tenant_id, kitchen_id, name, description, price, is_available, created_at FROM menu_items WHERE tenant_id = $1`

	rows, err := r.db(ctx).Query(ctx, query, tenantId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []*domain.MenuItem
	for rows.Next() {
		var item domain.MenuItem
		if err := rows.Scan(&item.Id, &item.TenantId, &item.KitchenId, &item.Name, &item.Description, &item.Price, &item.IsAvailable, &item.CreatedAt); err != nil {
			return nil, err
		}
		items = append(items, &item)
	}
	return items, nil
}

func (r *MenuItemRepository) Update(ctx context.Context, item *domain.MenuItem) error {
	query := `UPDATE menu_items SET kitchen_id = $1, name = $2, description = $3, price = $4, is_available = $5 WHERE id = $6 AND tenant_id = $7`

	cmd, err := r.db(ctx).Exec(ctx, query, item.KitchenId, item.Name, item.Description, item.Price, item.IsAvailable, item.Id, item.TenantId)
	if err != nil {
		return err
	}
	if cmd.RowsAffected() == 0 {
		return errors.New("menu item not found to update")
	}
	return nil
}

func (r *MenuItemRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM menu_items WHERE id = $1`

	cmd, err := r.db(ctx).Exec(ctx, query, id)
	if err != nil {
		return err
	}
	if cmd.RowsAffected() == 0 {
		return errors.New("menu item not found to delete")
	}
	return nil
}
