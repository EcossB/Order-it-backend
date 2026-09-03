package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Order struct {
	Id        uuid.UUID
	TenantId  uuid.UUID
	TableId   *uuid.UUID
	WaiterId  *uuid.UUID
	Status    string
	CreatedAt time.Time
}

type OrderRepository interface {
	Create(ctx context.Context, order *Order) error
	GetById(ctx context.Context, id uuid.UUID) (*Order, error)
	GetAllByTenantId(ctx context.Context, tenantId uuid.UUID) ([]*Order, error)
	Update(ctx context.Context, order *Order) error
	Delete(ctx context.Context, id uuid.UUID) error
}
