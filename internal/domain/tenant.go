package domain

import (
	"context"
	"time"
)

type Tenant struct {
	Id        int
	Name      string
	Type      string
	CreatedAt time.Time
}

type TenantRepository interface {
	Create(ctx context.Context, tenant *Tenant) error
	GetById(ctx context.Context, id int) (*Tenant, error)
	GetAll(ctx context.Context) ([]*Tenant, error)
	Update(ctx context.Context, tenant *Tenant) error
	Delete(ctx context.Context, id int) error
}
