package domain

import (
	"context"
	"time"
)

type User struct {
	Id        int
	TenantId  int
	Role      int  // Referencia a UserRole
	KitchenId *int // Puntero a int porque puede ser nulo para los meseros
	Name      string
	PinHash   string
	CreatedAt time.Time
}

type UserRepository interface {
	Create(ctx context.Context, user *User) error
	GetByName(ctx context.Context, name string) (*User, error)
	GetAll(ctx context.Context) ([]*User, error)
	Update(ctx context.Context, user *User) error
	Delete(ctx context.Context, id int) error
}
