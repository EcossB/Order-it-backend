package postgres

import (
	"context"
	"errors"
	"order-it-backend/internal/domain"
	"order-it-backend/internal/middleware"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TableRepository struct {
	pool *pgxpool.Pool
}

func NewTableRepository(pool *pgxpool.Pool) domain.TableRepository {
	return &TableRepository{pool: pool}
}

func (t *TableRepository) db(ctx context.Context) domain.DBTX {

	if tx, ok := middleware.GetTx(ctx); ok {
		return tx
	}

	return t.pool
}

func (t *TableRepository) Create(ctx context.Context, table *domain.Table) error {

	query := `INSERT INTO tables ( tenant_id, table_number) values ( $1, $2 ) returning id, status, created_at`

	err := t.db(ctx).QueryRow(ctx, query, table.TenantId, table.TableNumber).Scan(&table.Id, &table.Status, &table.CreatedAt)

	if err != nil {
		return err
	}

	return nil
}

func (t *TableRepository) GetById(ctx context.Context, id uuid.UUID) (*domain.Table, error) {

	query := `select id, tenant_id, table_number, status, created_at from tables WHERE id = $1`

	row := t.db(ctx).QueryRow(ctx, query, id)

	var table domain.Table

	err := row.Scan(&table.Id, &table.TenantId, &table.TableNumber, &table.Status, &table.CreatedAt)

	if err != nil {
		return nil, err
	}

	return &table, nil
}

func (t *TableRepository) GetAllByTenantId(ctx context.Context, tenantId uuid.UUID) ([]*domain.Table, error) {

	query := `select id, tenant_id, table_number, status, created_at from tables WHERE tenant_id = $1`

	row, err := t.db(ctx).Query(ctx, query, tenantId)

	if err != nil {
		return nil, err
	}

	var tables []*domain.Table

	for row.Next() {
		var table domain.Table
		err := row.Scan(&table.Id, &table.TenantId, &table.TableNumber, &table.Status, &table.CreatedAt)

		if err != nil {
			return nil, err
		}

		tables = append(tables, &table)
	}

	return tables, nil
}

func (t *TableRepository) GetByTenantIdAndName(ctx context.Context, tenantId uuid.UUID, name string) (*domain.Table, error) {

	query := `select id, tenant_id, table_number, status, created_at from tables where tenant_id = $1 and table_number = $2`
	row := t.db(ctx).QueryRow(ctx, query, tenantId, name)
	var table domain.Table
	err := row.Scan(&table.Id, &table.TenantId, &table.TableNumber, &table.Status, &table.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &table, nil
}

func (t *TableRepository) Update(ctx context.Context, table *domain.Table) error {

	query := `update tables set table_number = $1 where id = $2 and tenant_id = $3`

	cmd, err := t.db(ctx).Exec(ctx, query, table.TableNumber, table.Id, table.TenantId)
	if err != nil {
		return err
	}
	if cmd.RowsAffected() == 0 {
		return errors.New("table not found to update")
	}
	return nil
}

func (t *TableRepository) UpdateStatus(ctx context.Context, table *domain.Table) error {

	query := `update tables set status = $1 where id = $2 and tenant_id = $3`

	cmd, err := t.db(ctx).Exec(ctx, query, table.Status, table.Id, table.TenantId)

	if err != nil {
		return err
	}

	if cmd.RowsAffected() == 0 {
		return errors.New("table not found to update")
	}
	return nil
}

func (t *TableRepository) Delete(ctx context.Context, id uuid.UUID) error {

	query := `delete from tables where id = $1`

	cmd, err := t.db(ctx).Exec(ctx, query, id)

	if err != nil {
		return err
	}

	if cmd.RowsAffected() == 0 {
		return errors.New("kitchen not found to delete")
	}
	return nil

}
