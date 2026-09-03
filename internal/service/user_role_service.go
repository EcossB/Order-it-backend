package service

import (
	"context"
	"errors"
	"order-it-backend/internal/domain"
	"strings"
)

type UserRoleService struct {
	repository       domain.UserRoleRepository
	tenantRepository domain.TenantRepository
}

// injecting dependecy injection
func NewUserRoleService(repository domain.UserRoleRepository, tenantRepo domain.TenantRepository) *UserRoleService {
	return &UserRoleService{
		repository:       repository,
		tenantRepository: tenantRepo,
	}
}

func (service *UserRoleService) CreateUserRole(ctx context.Context, userRole *domain.UserRole) (*domain.UserRole, error) {

	if userRole.TenantId <= 0 {
		return nil, errors.New("el ID del tenant es obligatorio")
	}

	// 1. Validar que la organización (Tenant) realmente exista
	_, err := service.tenantRepository.GetById(ctx, userRole.TenantId)
	if err != nil {
		return nil, errors.New("la organización (tenant) especificada no existe")
	}

	userRole.Description = strings.TrimSpace(userRole.Description)
	if userRole.Description == "" {
		return nil, errors.New("es necesario agregar la descripción")
	}

	if userRole.Description != "WAITER" && userRole.Description != "CHEF" && userRole.Description != "ADMIN" {
		return nil, errors.New("el tipo de usuario debe ser WAITER, CHEF o ADMIN")
	}

	// 2. Validar que no exista un rol duplicado para este tenant
	existingRole, err := service.repository.GetByTenantAndDescription(ctx, userRole.TenantId, userRole.Description)
	if err != nil {
		return nil, err
	}
	if existingRole != nil {
		return nil, errors.New("este rol ya existe en la organización")
	}

	err = service.repository.Create(ctx, userRole)
	if err != nil {
		return nil, err
	}

	return userRole, nil
}

func (service *UserRoleService) DeleteUserRole(ctx context.Context, id int) error {
	if id <= 0 {
		return errors.New("id de rol de usuario inválido")
	}
	return service.repository.Delete(ctx, id)
}

func (service *UserRoleService) GetUserRoleById(ctx context.Context, id int) (*domain.UserRole, error) {
	if id <= 0 {
		return nil, errors.New("id de rol de usuario inválido")
	}
	return service.repository.GetById(ctx, id)
}

func (service *UserRoleService) GetAllUserRolesByTenantId(ctx context.Context, tenantId int) ([]*domain.UserRole, error) {
	if tenantId <= 0 {
		return nil, errors.New("el ID del tenant es obligatorio")
	}
	return service.repository.GetAllByTenantId(ctx, tenantId)
}

func (service *UserRoleService) UpdateUserRole(ctx context.Context, userRole *domain.UserRole) error {
	if userRole.Id <= 0 || userRole.TenantId <= 0 {
		return errors.New("el ID del rol y del tenant son obligatorios")
	}

	// Validar que la organización (Tenant) realmente exista
	_, err := service.tenantRepository.GetById(ctx, userRole.TenantId)
	if err != nil {
		return errors.New("la organización (tenant) especificada no existe")
	}

	userRole.Description = strings.TrimSpace(userRole.Description)
	if userRole.Description == "" {
		return errors.New("es necesario agregar la descripción")
	}

	if userRole.Description != "WAITER" && userRole.Description != "CHEF" && userRole.Description != "ADMIN" {
		return errors.New("el tipo de usuario debe ser WAITER, CHEF o ADMIN")
	}

	// Validar duplicados en la actualización
	existingRole, err := service.repository.GetByTenantAndDescription(ctx, userRole.TenantId, userRole.Description)
	if err != nil {
		return err
	}
	// Si existe, y el ID es diferente al que estamos intentando actualizar, entonces es un duplicado
	if existingRole != nil && existingRole.Id != userRole.Id {
		return errors.New("ya existe otro rol con esta descripción en la organización")
	}

	return service.repository.Update(ctx, userRole)
}
