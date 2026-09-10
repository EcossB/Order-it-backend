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

type OrderItemRepository struct {
	pool *pgxpool.Pool
}

func NewOrderItemRepository(pool *pgxpool.Pool) domain.OrderItemRepository {
	return &OrderItemRepository{pool: pool}
}

func (r *OrderItemRepository) db(ctx context.Context) domain.DBTX {
	if tx, ok := middleware.GetTx(ctx); ok {
		return tx
	}
	return r.pool
}

func (r *OrderItemRepository) Create(ctx context.Context, item *domain.OrderItem) error {
	query := `INSERT INTO order_items (tenant_id, order_id, menu_item_id, kitchen_id, quantity, status, notes) 
			  VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id, created_at`

	err := r.db(ctx).QueryRow(ctx, query, item.TenantId, item.OrderId, item.MenuItemId, item.KitchenId, item.Quantity, item.Status, item.Notes).Scan(&item.Id, &item.CreatedAt)
	if err != nil {
		return err
	}
	return nil
}

func (r *OrderItemRepository) CreateBulk(ctx context.Context, orderItems []*domain.OrderItem) error {
	if len(orderItems) == 0 {
		return nil
	}

	batch := &pgx.Batch{}
	query := `INSERT INTO order_items (tenant_id, order_id, menu_item_id, kitchen_id, quantity, status, notes) 
			  VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id, created_at`

	for _, item := range orderItems {
		batch.Queue(query, item.TenantId, item.OrderId, item.MenuItemId, item.KitchenId, item.Quantity, item.Status, item.Notes)
	}

	br := r.db(ctx).SendBatch(ctx, batch)
	defer br.Close()

	for _, item := range orderItems {
		err := br.QueryRow().Scan(&item.Id, &item.CreatedAt)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *OrderItemRepository) GetById(ctx context.Context, id uuid.UUID) (*domain.OrderItem, error) {
	query := `SELECT id, tenant_id, order_id, menu_item_id, kitchen_id, quantity, status, notes, created_at FROM order_items WHERE id = $1`

	var item domain.OrderItem
	err := r.db(ctx).QueryRow(ctx, query, id).Scan(
		&item.Id, &item.TenantId, &item.OrderId, &item.MenuItemId, &item.KitchenId, &item.Quantity, &item.Status, &item.Notes, &item.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *OrderItemRepository) GetAllByOrderId(ctx context.Context, orderId uuid.UUID) ([]*domain.OrderItem, error) {
	query := `SELECT id, tenant_id, order_id, menu_item_id, kitchen_id, quantity, status, notes, created_at FROM order_items WHERE order_id = $1 ORDER BY created_at ASC`

	rows, err := r.db(ctx).Query(ctx, query, orderId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []*domain.OrderItem
	for rows.Next() {
		var item domain.OrderItem
		if err := rows.Scan(&item.Id, &item.TenantId, &item.OrderId, &item.MenuItemId, &item.KitchenId, &item.Quantity, &item.Status, &item.Notes, &item.CreatedAt); err != nil {
			return nil, err
		}
		items = append(items, &item)
	}
	return items, nil
}

func (r *OrderItemRepository) Update(ctx context.Context, item *domain.OrderItem) error {
	query := `UPDATE order_items SET quantity = $1, status = $2, notes = $3 WHERE id = $4 AND tenant_id = $5`

	cmd, err := r.db(ctx).Exec(ctx, query, item.Quantity, item.Status, item.Notes, item.Id, item.TenantId)
	if err != nil {
		return err
	}
	if cmd.RowsAffected() == 0 {
		return errors.New("order item not found to update")
	}
	return nil
}

func (r *OrderItemRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM order_items WHERE id = $1`

	cmd, err := r.db(ctx).Exec(ctx, query, id)
	if err != nil {
		return err
	}
	if cmd.RowsAffected() == 0 {
		return errors.New("order item not found to delete")
	}
	return nil
}
