package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Tenant struct {
	Id        uuid.UUID
	Name      string
	Type      string
	CreatedAt time.Time
}

type TenantRepository interface {
	Create(ctx context.Context, tenant *Tenant) error
	GetById(ctx context.Context, id uuid.UUID) (*Tenant, error)
	GetAll(ctx context.Context) ([]*Tenant, error)
	Update(ctx context.Context, tenant *Tenant) error
	Delete(ctx context.Context, id uuid.UUID) error
}
