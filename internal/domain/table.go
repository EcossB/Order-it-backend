package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Table struct {
	Id          uuid.UUID
	TenantId    uuid.UUID
	TableNumber string
	Status      string
	CreatedAt   time.Time
}

type TableRepository interface {
	Create(ctx context.Context, table *Table) error
	GetById(ctx context.Context, id uuid.UUID) (*Table, error)
	GetAllByTenantId(ctx context.Context, tenantId uuid.UUID) ([]*Table, error)
	Update(ctx context.Context, table *Table) error
	Delete(ctx context.Context, id uuid.UUID) error
}
