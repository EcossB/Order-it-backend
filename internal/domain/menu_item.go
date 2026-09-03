package domain

import (
	"context"
	"time"
)

type MenuItem struct {
	Id          int
	KitchenId   int
	Name        string
	Description string
	Price       float64
	CreatedAt   time.Time
}

type MenuItemRepository interface {
	Create(ctx context.Context, menuItem *MenuItem) error
	GetById(ctx context.Context, id int) (*MenuItem, error)
	GetAll(ctx context.Context) ([]*MenuItem, error)
	Update(ctx context.Context, menuItem *MenuItem) error
	Delete(ctx context.Context, id int) error
}
