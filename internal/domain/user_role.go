package domain

import (
	"context"
	"time"
)

type UserRole struct {
	Id          int
	Description string
	Status      int
	CreatedAt   time.Time
}

type UserRoleRepository interface {
	Create(ctx context.Context, userRole *UserRole) error
	GetById(ctx context.Context, id int) (*UserRole, error)
	GetAll(ctx context.Context) ([]*UserRole, error)
	Update(ctx context.Context, userRole *UserRole) error
	Delete(ctx context.Context, id int) error
}
