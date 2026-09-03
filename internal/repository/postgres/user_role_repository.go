package postgres

import (
	"context"
	"errors"
	"order-it-backend/internal/domain"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRoleRepository struct {
	db *pgxpool.Pool
}

func NewUserRoleRepository(db *pgxpool.Pool) domain.UserRoleRepository {
	return &UserRoleRepository{db: db}
}

func (u *UserRoleRepository) Create(ctx context.Context, userRole *domain.UserRole) error {
	query := `INSERT INTO users_role (tenant_id, description) VALUES ($1, $2) RETURNING id, status, created_at;`

	err := u.db.QueryRow(ctx, query, userRole.TenantId, userRole.Description).Scan(&userRole.Id, &userRole.Status, &userRole.CreatedAt)
	if err != nil {
		return err
	}
	return nil
}

func (u *UserRoleRepository) GetById(ctx context.Context, id int) (*domain.UserRole, error) {
	query := `SELECT id, tenant_id, description, status, created_at FROM users_role WHERE id = $1;`

	var userRole domain.UserRole
	err := u.db.QueryRow(ctx, query, id).Scan(&userRole.Id, &userRole.TenantId, &userRole.Description, &userRole.Status, &userRole.CreatedAt)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("rol de usuario no encontrado")
		}
		return nil, err
	}
	return &userRole, nil
}

func (u *UserRoleRepository) GetByTenantAndDescription(ctx context.Context, tenantId int, description string) (*domain.UserRole, error) {
	query := `SELECT id, tenant_id, description, status, created_at FROM users_role WHERE tenant_id = $1 AND description = $2;`

	var userRole domain.UserRole
	err := u.db.QueryRow(ctx, query, tenantId, description).Scan(&userRole.Id, &userRole.TenantId, &userRole.Description, &userRole.Status, &userRole.CreatedAt)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil // Retornamos nil, nil para indicar que NO existe (lo cual es bueno para crear)
		}
		return nil, err // Cualquier otro error de base de datos
	}
	return &userRole, nil
}

func (u *UserRoleRepository) GetAllByTenantId(ctx context.Context, tenantId int) ([]*domain.UserRole, error) {
	// Ahora filtramos por tenant_id
	query := `SELECT id, tenant_id, description, status, created_at FROM users_role WHERE status = 1 AND tenant_id = $1;`

	var userRoles []*domain.UserRole
	rows, err := u.db.Query(ctx, query, tenantId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var userRole domain.UserRole
		if err := rows.Scan(&userRole.Id, &userRole.TenantId, &userRole.Description, &userRole.Status, &userRole.CreatedAt); err != nil {
			return nil, err
		}
		userRoles = append(userRoles, &userRole)
	}

	return userRoles, nil
}

func (u *UserRoleRepository) Update(ctx context.Context, userRole *domain.UserRole) error {
	query := `UPDATE users_role SET description = $1 WHERE id = $2 AND tenant_id = $3;`

	commandTag, err := u.db.Exec(ctx, query, userRole.Description, userRole.Id, userRole.TenantId)
	if err != nil {
		return err
	}

	if commandTag.RowsAffected() == 0 {
		return errors.New("rol de usuario no encontrado para actualizar")
	}
	return nil
}

func (u *UserRoleRepository) Delete(ctx context.Context, id int) error {
	query := `UPDATE users_role SET status = 0 WHERE id = $1;`
	commandTag, err := u.db.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	if commandTag.RowsAffected() == 0 {
		return errors.New("rol de usuario no encontrado para eliminar")
	}
	return nil
}
