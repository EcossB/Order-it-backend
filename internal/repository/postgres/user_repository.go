package postgres

import (
	"context"
	"errors"
	"order-it-backend/internal/domain"
	"order-it-backend/internal/middleware"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	pool *pgxpool.Pool
}

func NewUserRepository(pool *pgxpool.Pool) domain.UserRepository {
	return &UserRepository{pool: pool}
}

func (r *UserRepository) db(ctx context.Context) domain.DBTX {
	if tx, ok := middleware.GetTx(ctx); ok {
		return tx
	}
	return r.pool
}

func (r *UserRepository) Create(ctx context.Context, user *domain.User) error {
	query := `INSERT INTO users (tenant_id, role, kitchen_id, name, email, password_hash, pin_hash) 
			  VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id, created_at`

	err := r.db(ctx).QueryRow(ctx, query, user.TenantId, user.Role, user.KitchenId, user.Name, user.Email, user.PasswordHash, user.PinHash).Scan(&user.Id, &user.CreatedAt)
	if err != nil {
		return err
	}
	return nil
}

func (r *UserRepository) GetById(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	query := `SELECT id, tenant_id, role, kitchen_id, name, email, password_hash, pin_hash, created_at FROM users WHERE id = $1`

	var user domain.User
	err := r.db(ctx).QueryRow(ctx, query, id).Scan(
		&user.Id, &user.TenantId, &user.Role, &user.KitchenId, &user.Name, &user.Email, &user.PasswordHash, &user.PinHash, &user.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) GetAllByTenantId(ctx context.Context, tenantId uuid.UUID) ([]*domain.User, error) {
	query := `SELECT id, tenant_id, role, kitchen_id, name, email, password_hash, pin_hash, created_at FROM users WHERE tenant_id = $1`

	rows, err := r.db(ctx).Query(ctx, query, tenantId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*domain.User
	for rows.Next() {
		var user domain.User
		if err := rows.Scan(&user.Id, &user.TenantId, &user.Role, &user.KitchenId, &user.Name, &user.Email, &user.PasswordHash, &user.PinHash, &user.CreatedAt); err != nil {
			return nil, err
		}
		users = append(users, &user)
	}
	return users, nil
}

func (r *UserRepository) GetByTenantIdAndName(ctx context.Context, tenantId uuid.UUID, name string) (*domain.User, error) {
	query := `SELECT id, tenant_id, role, kitchen_id, name, email, password_hash, pin_hash, created_at FROM users WHERE tenant_id = $1 AND name = $2`

	var user domain.User
	err := r.db(ctx).QueryRow(ctx, query, tenantId, name).Scan(
		&user.Id, &user.TenantId, &user.Role, &user.KitchenId, &user.Name, &user.Email, &user.PasswordHash, &user.PinHash, &user.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// Opcional: Para buscar por email en el signup
func (r *UserRepository) GetByTenantIdAndEmail(ctx context.Context, tenantId uuid.UUID, email string) (*domain.User, error) {
	query := `SELECT id, tenant_id, role, kitchen_id, name, email, password_hash, pin_hash, created_at FROM users WHERE tenant_id = $1 AND email = $2`

	var user domain.User
	err := r.db(ctx).QueryRow(ctx, query, tenantId, email).Scan(
		&user.Id, &user.TenantId, &user.Role, &user.KitchenId, &user.Name, &user.Email, &user.PasswordHash, &user.PinHash, &user.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) Update(ctx context.Context, user *domain.User) error {
	query := `UPDATE users SET role = $1, kitchen_id = $2, name = $3, email = $4, password_hash = $5, pin_hash = $6 WHERE id = $7 AND tenant_id = $8`

	cmd, err := r.db(ctx).Exec(ctx, query, user.Role, user.KitchenId, user.Name, user.Email, user.PasswordHash, user.PinHash, user.Id, user.TenantId)
	if err != nil {
		return err
	}
	if cmd.RowsAffected() == 0 {
		return errors.New("user not found to update")
	}
	return nil
}

func (r *UserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM users WHERE id = $1`

	cmd, err := r.db(ctx).Exec(ctx, query, id)
	if err != nil {
		return err
	}
	if cmd.RowsAffected() == 0 {
		return errors.New("user not found to delete")
	}
	return nil
}
