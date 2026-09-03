package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type OrderItem struct {
	Id         uuid.UUID
	TenantId   uuid.UUID
	OrderId    uuid.UUID
	MenuItemId uuid.UUID
	KitchenId  uuid.UUID
	Quantity   int
	Status     string
	Notes      string
	CreatedAt  time.Time
}

type OrderItemRepository interface {
	Create(ctx context.Context, orderItem *OrderItem) error
	GetById(ctx context.Context, id uuid.UUID) (*OrderItem, error)
	GetAllByOrderId(ctx context.Context, orderId uuid.UUID) ([]*OrderItem, error)
	Update(ctx context.Context, orderItem *OrderItem) error
	Delete(ctx context.Context, id uuid.UUID) error
}
