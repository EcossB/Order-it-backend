package service

import (
	"context"
	"errors"
	"order-it-backend/internal/domain"
	"strings"

	"github.com/google/uuid"
)

type KitchenService struct {
	repository       domain.KitchenRepository
	tenantRepository domain.TenantRepository
}

func NewKitchenService(repository domain.KitchenRepository, tenantRepository domain.TenantRepository) *KitchenService {
	return &KitchenService{
		repository:       repository,
		tenantRepository: tenantRepository,
	}
}

func (k *KitchenService) CreateKitchen(ctx context.Context, kitchen *domain.Kitchen) (*domain.Kitchen, error) {
	if kitchen.TenantId == uuid.Nil {
		return nil, errors.New("el ID de la organización es obligatorio")
	}

	_, err := k.tenantRepository.GetById(ctx, kitchen.TenantId)
	if err != nil {
		return nil, errors.New("la organización (tenant) especificada no existe")
	}

	kitchen.Name = strings.TrimSpace(kitchen.Name)
	if kitchen.Name == "" {
		return nil, errors.New("el nombre de la cocina no puede estar vacío")
	}

	existingKitchen, err := k.repository.GetByTenantIdAndName(ctx, kitchen.TenantId, kitchen.Name)
	if err != nil {
		return nil, err
	}
	if existingKitchen != nil {
		return nil, errors.New("esta cocina ya esta definida en esta organización")
	}

	err = k.repository.Create(ctx, kitchen)
	if err != nil {
		return nil, err
	}
	return kitchen, nil
}

func (k *KitchenService) GetKitchenById(ctx context.Context, id uuid.UUID) (*domain.Kitchen, error) {
	if id == uuid.Nil {
		return nil, errors.New("id de cocina inválido")
	}
	return k.repository.GetById(ctx, id)
}

func (k *KitchenService) GetAllKitchensByTenantId(ctx context.Context, tenantId uuid.UUID) ([]*domain.Kitchen, error) {
	if tenantId == uuid.Nil {
		return nil, errors.New("el ID de la organización es obligatorio")
	}
	return k.repository.GetAllByTenantId(ctx, tenantId)
}

func (k *KitchenService) UpdateKitchen(ctx context.Context, kitchen *domain.Kitchen) error {
	if kitchen.Id == uuid.Nil || kitchen.TenantId == uuid.Nil {
		return errors.New("el ID de la cocina y la organización son obligatorios")
	}

	_, err := k.tenantRepository.GetById(ctx, kitchen.TenantId)
	if err != nil {
		return errors.New("la organización (tenant) especificada no existe")
	}

	kitchen.Name = strings.TrimSpace(kitchen.Name)
	if kitchen.Name == "" {
		return errors.New("el nombre de la cocina no puede estar vacío")
	}

	existingKitchen, err := k.repository.GetByTenantIdAndName(ctx, kitchen.TenantId, kitchen.Name)
	if err != nil {
		return err
	}
	if existingKitchen != nil && existingKitchen.Id != kitchen.Id {
		return errors.New("ya existe otra cocina con este nombre en la organización")
	}

	return k.repository.Update(ctx, kitchen)
}

func (k *KitchenService) DeleteKitchen(ctx context.Context, id uuid.UUID) error {
	if id == uuid.Nil {
		return errors.New("id de cocina inválido")
	}
	return k.repository.Delete(ctx, id)
}
