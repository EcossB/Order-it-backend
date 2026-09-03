package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Kitchen struct {
	Id        uuid.UUID
	TenantId  uuid.UUID
	Name      string
	CreatedAt time.Time
}

type KitchenRepository interface {
	Create(ctx context.Context, kitchen *Kitchen) error
	GetById(ctx context.Context, id uuid.UUID) (*Kitchen, error)
	GetAllByTenantId(ctx context.Context, tenantId uuid.UUID) ([]*Kitchen, error)
	GetAll(ctx context.Context) ([]*Kitchen, error)
	GetByTenantIdAndName(ctx context.Context, tenantId uuid.UUID, name string) (*Kitchen, error)
	Update(ctx context.Context, kitchen *Kitchen) error
	Delete(ctx context.Context, id uuid.UUID) error
}
