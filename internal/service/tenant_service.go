package service

import (
	"context"
	"errors"
	"strings"

	"order-it-backend/internal/domain"
)

// TenantService contiene la lógica de negocio (casos de uso) para las organizaciones
type TenantService struct {
	repo domain.TenantRepository
}

// NewTenantService inicializa el servicio inyectando el repositorio
func NewTenantService(repo domain.TenantRepository) *TenantService {
	return &TenantService{
		repo: repo,
	}
}

// CreateTenant valida y procesa la creación de un nuevo Tenant
func (s *TenantService) CreateTenant(ctx context.Context, tenant *domain.Tenant) error {
	// 1. Validaciones de Negocio (Business Rules)
	tenant.Name = strings.TrimSpace(tenant.Name)
	if tenant.Name == "" {
		return errors.New("el nombre de la organización no puede estar vacío")
	}

	tenant.Type = strings.ToUpper(strings.TrimSpace(tenant.Type))
	if tenant.Type != "INDEPENDENT" && tenant.Type != "PARK" {
		return errors.New("el tipo de organización debe ser INDEPENDENT o PARK")
	}

	// 2. Guardar en la base de datos a través de la interfaz
	return s.repo.Create(ctx, tenant)
}

func (s *TenantService) GetTenantById(ctx context.Context, id int) (*domain.Tenant, error) {
	if id <= 0 {
		return nil, errors.New("id de organización inválido")
	}
	return s.repo.GetById(ctx, id)
}

func (s *TenantService) GetAllTenants(ctx context.Context) ([]*domain.Tenant, error) {
	return s.repo.GetAll(ctx)
}

func (s *TenantService) UpdateTenant(ctx context.Context, tenant *domain.Tenant) error {
	if tenant.Id <= 0 {
		return errors.New("id de organización inválido")
	}

	tenant.Name = strings.TrimSpace(tenant.Name)
	if tenant.Name == "" {
		return errors.New("el nombre de la organización no puede estar vacío")
	}

	tenant.Type = strings.ToUpper(strings.TrimSpace(tenant.Type))
	if tenant.Type != "INDEPENDENT" && tenant.Type != "PARK" {
		return errors.New("el tipo de organización debe ser INDEPENDENT o PARK")
	}

	return s.repo.Update(ctx, tenant)
}

func (s *TenantService) DeleteTenant(ctx context.Context, id int) error {
	if id <= 0 {
		return errors.New("id de organización inválido")
	}
	// Aquí podrías agregar más lógica de negocio, por ejemplo:
	// "No permitir borrar un parque si tiene cocinas activas" (validándolo antes)

	return s.repo.Delete(ctx, id)
}
