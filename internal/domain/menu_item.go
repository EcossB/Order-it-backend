package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type MenuItem struct {
	Id          uuid.UUID
	TenantId    uuid.UUID
	KitchenId   *uuid.UUID // Asumiendo que puede ser nulo, o uuid.UUID si no
	Name        string
	Description string
	Price       float64
	CreatedAt   time.Time
}

type MenuItemRepository interface {
	Create(ctx context.Context, menuItem *MenuItem) error
	GetById(ctx context.Context, id uuid.UUID) (*MenuItem, error)
	GetAllByTenantId(ctx context.Context, tenantId uuid.UUID) ([]*MenuItem, error)
	Update(ctx context.Context, menuItem *MenuItem) error
	Delete(ctx context.Context, id uuid.UUID) error
}
