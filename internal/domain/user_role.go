package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type UserRole struct {
	Id          uuid.UUID
	TenantId    uuid.UUID
	Description string
	Status      int
	CreatedAt   time.Time
}

type UserRoleRepository interface {
	Create(ctx context.Context, userRole *UserRole) error
	GetById(ctx context.Context, id uuid.UUID) (*UserRole, error)
	GetAllByTenantId(ctx context.Context, tenantId uuid.UUID) ([]*UserRole, error)
	GetByTenantAndDescription(ctx context.Context, tenantId uuid.UUID, description string) (*UserRole, error)
	Update(ctx context.Context, userRole *UserRole) error
	Delete(ctx context.Context, id uuid.UUID) error
}
