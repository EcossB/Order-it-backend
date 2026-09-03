package domain

import (
	"context"
	"time"
)

type Kitchen struct {
	Id        int
	TenantId  int
	Name      string
	CreatedAt time.Time
}

type KitchenRepository interface {
	Create(ctx context.Context, kitchen *Kitchen) error
	GetById(ctx context.Context, id int) (*Kitchen, error)
	GetAll(ctx context.Context) ([]*Kitchen, error)
	Update(ctx context.Context, kitchen *Kitchen) error
	Delete(ctx context.Context, id int) error
}
