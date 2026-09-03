package domain

import (
	"context"
	"time"
)

type Order struct {
	Id        int
	TenantId  int
	TableId   int
	WaiterId  int
	Status    string
	CreatedAt time.Time
}

type OrderRepository interface {
	Create(ctx context.Context, order *Order) error
	GetById(ctx context.Context, id int) (*Order, error)
	GetAll(ctx context.Context) ([]*Order, error)
	Update(ctx context.Context, order *Order) error
}
