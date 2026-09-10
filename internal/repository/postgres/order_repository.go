package postgres

import (
	"context"
	"errors"
	"order-it-backend/internal/domain"
	"order-it-backend/internal/middleware"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type OrderRepository struct {
	pool *pgxpool.Pool
}

func NewOrderRepository(pool *pgxpool.Pool) domain.OrderRepository {
	return &OrderRepository{pool: pool}
}

func (r *OrderRepository) db(ctx context.Context) domain.DBTX {
	if tx, ok := middleware.GetTx(ctx); ok {
		return tx
	}
	return r.pool
}

func (r *OrderRepository) Create(ctx context.Context, order *domain.Order) error {
	query := `INSERT INTO orders (tenant_id, table_id, waiter_id, status) 
			  VALUES ($1, $2, $3, $4) RETURNING id, created_at`

	err := r.db(ctx).QueryRow(ctx, query, order.TenantId, order.TableId, order.WaiterId, order.Status).Scan(&order.Id, &order.CreatedAt)
	if err != nil {
		return err
	}
	return nil
}

func (r *OrderRepository) GetById(ctx context.Context, id uuid.UUID) (*domain.Order, error) {
	query := `SELECT id, tenant_id, table_id, waiter_id, status, created_at FROM orders WHERE id = $1`

	var order domain.Order
	err := r.db(ctx).QueryRow(ctx, query, id).Scan(
		&order.Id, &order.TenantId, &order.TableId, &order.WaiterId, &order.Status, &order.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &order, nil
}

func (r *OrderRepository) GetAllByTenantId(ctx context.Context, tenantId uuid.UUID) ([]*domain.Order, error) {
	query := `SELECT id, tenant_id, table_id, waiter_id, status, created_at FROM orders WHERE tenant_id = $1 ORDER BY created_at DESC`

	rows, err := r.db(ctx).Query(ctx, query, tenantId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []*domain.Order
	for rows.Next() {
		var order domain.Order
		if err := rows.Scan(&order.Id, &order.TenantId, &order.TableId, &order.WaiterId, &order.Status, &order.CreatedAt); err != nil {
			return nil, err
		}
		orders = append(orders, &order)
	}
	return orders, nil
}

func (r *OrderRepository) Update(ctx context.Context, order *domain.Order) error {
	query := `UPDATE orders SET table_id = $1, waiter_id = $2, status = $3 WHERE id = $4 AND tenant_id = $5`

	cmd, err := r.db(ctx).Exec(ctx, query, order.TableId, order.WaiterId, order.Status, order.Id, order.TenantId)
	if err != nil {
		return err
	}
	if cmd.RowsAffected() == 0 {
		return errors.New("order not found to update")
	}
	return nil
}

func (r *OrderRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM orders WHERE id = $1`

	cmd, err := r.db(ctx).Exec(ctx, query, id)
	if err != nil {
		return err
	}
	if cmd.RowsAffected() == 0 {
		return errors.New("order not found to delete")
	}
	return nil
}
