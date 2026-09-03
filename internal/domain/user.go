package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type User struct {
	Id           uuid.UUID
	TenantId     uuid.UUID
	Role         uuid.UUID  // Referencia a UserRole
	KitchenId    *uuid.UUID // Puntero porque puede ser nulo para los meseros
	Name         string
	Email        *string // Puede ser nulo para meseros y chefs
	PasswordHash *string // Puede ser nulo para meseros y chefs
	PinHash      *string // Puede ser nulo para admins
	CreatedAt    time.Time
}

type UserRepository interface {
	Create(ctx context.Context, user *User) error
	GetById(ctx context.Context, id uuid.UUID) (*User, error)
	GetAllByTenantId(ctx context.Context, tenantId uuid.UUID) ([]*User, error)
	GetByTenantIdAndName(ctx context.Context, tenantId uuid.UUID, name string) (*User, error)
	Update(ctx context.Context, user *User) error
	Delete(ctx context.Context, id uuid.UUID) error
}
