package domain

import (
	"context"
	"time"
)

type UserRole struct {
	Id          int
	TenantId    int // Agregado para soportar Multi-Tenant
	Description string
	Status      int
	CreatedAt   time.Time
}

type UserRoleRepository interface {
	Create(ctx context.Context, userRole *UserRole) error
	GetById(ctx context.Context, id int) (*UserRole, error)
	// Cambiado de GetAll a GetAllByTenantId para filtrar por Organización
	GetAllByTenantId(ctx context.Context, tenantId int) ([]*UserRole, error)

	// Método para validar duplicados
	GetByTenantAndDescription(ctx context.Context, tenantId int, description string) (*UserRole, error)

	Update(ctx context.Context, userRole *UserRole) error
	Delete(ctx context.Context, id int) error
}
