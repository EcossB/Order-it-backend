package domain

import (
	"context"
	"time"
)

type Table struct {
	Id          int
	TenantId    int
	TableNumber string
	Status      string
	CreatedAt   time.Time
}

type TableRepository interface {
	Create(ctx context.Context, table *Table) error
	GetById(ctx context.Context, id int) (*Table, error)
	GetAll(ctx context.Context) ([]*Table, error)
	Update(ctx context.Context, table *Table) error
	Delete(ctx context.Context, id int) error
}
