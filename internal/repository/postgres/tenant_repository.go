package postgres

import (
	"context"
	"errors"

	"order-it-backend/internal/domain"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// tenantRepository es la implementación concreta (postgres) de domain.TenantRepository
type tenantRepository struct {
	db *pgxpool.Pool
}

// NewTenantRepository crea una nueva instancia del repositorio inyectando la base de datos
func NewTenantRepository(db *pgxpool.Pool) domain.TenantRepository {
	return &tenantRepository{
		db: db,
	}
}

func (r *tenantRepository) Create(ctx context.Context, tenant *domain.Tenant) error {
	query := `
		INSERT INTO tenants (name, type)
		VALUES ($1, $2)
		RETURNING id, created_at
	`

	// Usamos QueryRow porque queremos recuperar el ID y CreatedAt que generó PostgreSQL automáticamente
	err := r.db.QueryRow(ctx, query, tenant.Name, tenant.Type).Scan(&tenant.Id, &tenant.CreatedAt)
	if err != nil {
		return err
	}
	return nil
}

func (r *tenantRepository) GetById(ctx context.Context, id int) (*domain.Tenant, error) {
	query := `
		SELECT id, name, type, created_at
		FROM tenants
		WHERE id = $1
	`
	var tenant domain.Tenant
	err := r.db.QueryRow(ctx, query, id).Scan(
		&tenant.Id,
		&tenant.Name,
		&tenant.Type,
		&tenant.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("tenant not found") // Traducimos el error de PG a un error de negocio
		}
		return nil, err
	}

	return &tenant, nil
}

func (r *tenantRepository) GetAll(ctx context.Context) ([]*domain.Tenant, error) {
	query := `
		SELECT id, name, type, created_at
		FROM tenants
		ORDER BY created_at DESC
	`
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tenants []*domain.Tenant
	for rows.Next() {
		var t domain.Tenant
		if err := rows.Scan(
			&t.Id,
			&t.Name,
			&t.Type,
			&t.CreatedAt,
		); err != nil {
			return nil, err
		}
		tenants = append(tenants, &t)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return tenants, nil
}

func (r *tenantRepository) Update(ctx context.Context, tenant *domain.Tenant) error {
	query := `
		UPDATE tenants
		SET name = $1, type = $2
		WHERE id = $3
	`
	commandTag, err := r.db.Exec(ctx, query, tenant.Name, tenant.Type, tenant.Id)
	if err != nil {
		return err
	}

	if commandTag.RowsAffected() == 0 {
		return errors.New("tenant not found to update")
	}

	return nil
}

func (r *tenantRepository) Delete(ctx context.Context, id int) error {
	query := `
		DELETE FROM tenants
		WHERE id = $1
	`
	commandTag, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	if commandTag.RowsAffected() == 0 {
		return errors.New("tenant not found to delete")
	}

	return nil
}
