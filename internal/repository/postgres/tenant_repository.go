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

type TenantRepository struct {
	pool *pgxpool.Pool
}

func NewTenantRepository(pool *pgxpool.Pool) domain.TenantRepository {
	return &TenantRepository{pool: pool}
}

// db obtiene la transacción del middleware (si existe), de lo contrario usa el pool directo
func (r *TenantRepository) db(ctx context.Context) domain.DBTX {
	if tx, ok := middleware.GetTx(ctx); ok {
		return tx
	}
	return r.pool
}

func (r *TenantRepository) Create(ctx context.Context, tenant *domain.Tenant) error {
	query := `INSERT INTO tenants (name, type) VALUES ($1, $2) RETURNING id, created_at`
	err := r.db(ctx).QueryRow(ctx, query, tenant.Name, tenant.Type).Scan(&tenant.Id, &tenant.CreatedAt)
	if err != nil {
		return err
	}
	return nil
}

func (r *TenantRepository) GetAll(ctx context.Context) ([]*domain.Tenant, error) {
	query := `SELECT id, name, type, created_at FROM tenants`
	rows, err := r.db(ctx).Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tenants []*domain.Tenant
	for rows.Next() {
		var tenant domain.Tenant
		err := rows.Scan(&tenant.Id, &tenant.Name, &tenant.Type, &tenant.CreatedAt)
		if err != nil {
			return nil, err
		}
		tenants = append(tenants, &tenant)
	}

	return tenants, nil
}

func (r *TenantRepository) GetById(ctx context.Context, id uuid.UUID) (*domain.Tenant, error) {
	query := `SELECT id, name, type, created_at FROM tenants WHERE id = $1`
	var tenant domain.Tenant

	err := r.db(ctx).QueryRow(ctx, query, id).Scan(&tenant.Id, &tenant.Name, &tenant.Type, &tenant.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("tenant no encontrado")
		}
		return nil, err
	}
	return &tenant, nil
}

func (r *TenantRepository) Update(ctx context.Context, tenant *domain.Tenant) error {
	query := `UPDATE tenants SET name = $1, type = $2 WHERE id = $3`
	commandTag, err := r.db(ctx).Exec(ctx, query, tenant.Name, tenant.Type, tenant.Id)
	if err != nil {
		return err
	}
	if commandTag.RowsAffected() == 0 {
		return errors.New("tenant no encontrado para actualizar")
	}
	return nil
}

func (r *TenantRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM tenants WHERE id = $1`
	commandTag, err := r.db(ctx).Exec(ctx, query, id)
	if err != nil {
		return err
	}
	if commandTag.RowsAffected() == 0 {
		return errors.New("tenant no encontrado para eliminar")
	}
	return nil
}
