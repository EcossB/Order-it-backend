package domain

import (
	"context"
	"time"
)

type OrderItem struct {
	Id         int
	OrderId    int
	MenuItemId int
	KitchenId  int
	Quantity   int
	Status     string
	Notes      string
	CreatedAt  time.Time
}

type OrderItemRepository interface {
	Create(ctx context.Context, orderItem *OrderItem) error
	GetById(ctx context.Context, id int) (*OrderItem, error)
	GetAll(ctx context.Context) ([]*OrderItem, error)
	Update(ctx context.Context, orderItem *OrderItem) error
}
