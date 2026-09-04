package service

import (
	"context"
	"errors"
	"strings"

	"order-it-backend/internal/domain"

	"github.com/google/uuid"
)

type MenuItemService struct {
	repo        domain.MenuItemRepository
	tenantRepo  domain.TenantRepository
	kitchenRepo domain.KitchenRepository
}

func NewMenuItemService(repo domain.MenuItemRepository, tenantRepo domain.TenantRepository, kitchenRepo domain.KitchenRepository) *MenuItemService {
	return &MenuItemService{repo: repo, tenantRepo: tenantRepo, kitchenRepo: kitchenRepo}
}

func (s *MenuItemService) Create(ctx context.Context, item *domain.MenuItem) (*domain.MenuItem, error) {
	if item.TenantId == uuid.Nil {
		return nil, errors.New("invalid tenant id")
	}
	if item.KitchenId == uuid.Nil {
		return nil, errors.New("invalid kitchen id")
	}

	item.Name = strings.TrimSpace(item.Name)
	if item.Name == "" {
		return nil, errors.New("name cannot be empty")
	}
	if item.Price < 0 {
		return nil, errors.New("price cannot be negative")
	}

	_, err := s.tenantRepo.GetById(ctx, item.TenantId)
	if err != nil {
		return nil, errors.New("tenant does not exist")
	}

	kitchen, err := s.kitchenRepo.GetById(ctx, item.KitchenId)
	if err != nil {
		return nil, errors.New("kitchen does not exist")
	}

	if kitchen.TenantId != item.TenantId {
		return nil, errors.New("kitchen does not belong to this tenant")
	}

	err = s.repo.Create(ctx, item)
	if err != nil {
		return nil, err
	}

	return item, nil
}

func (s *MenuItemService) GetById(ctx context.Context, id uuid.UUID) (*domain.MenuItem, error) {
	if id == uuid.Nil {
		return nil, errors.New("invalid menu item id")
	}
	return s.repo.GetById(ctx, id)
}

func (s *MenuItemService) GetAllByTenantId(ctx context.Context, tenantId uuid.UUID) ([]*domain.MenuItem, error) {
	if tenantId == uuid.Nil {
		return nil, errors.New("invalid tenant id")
	}
	return s.repo.GetAllByTenantId(ctx, tenantId)
}

func (s *MenuItemService) Update(ctx context.Context, item *domain.MenuItem) error {
	if item.Id == uuid.Nil {
		return errors.New("invalid menu item id")
	}
	if item.TenantId == uuid.Nil {
		return errors.New("invalid tenant id")
	}
	if item.KitchenId == uuid.Nil {
		return errors.New("invalid kitchen id")
	}

	item.Name = strings.TrimSpace(item.Name)
	if item.Name == "" {
		return errors.New("name cannot be empty")
	}
	if item.Price < 0 {
		return errors.New("price cannot be negative")
	}

	kitchen, err := s.kitchenRepo.GetById(ctx, item.KitchenId)
	if err != nil {
		return errors.New("kitchen does not exist")
	}
	if kitchen.TenantId != item.TenantId {
		return errors.New("kitchen does not belong to this tenant")
	}

	return s.repo.Update(ctx, item)
}

func (s *MenuItemService) Delete(ctx context.Context, id uuid.UUID) error {
	if id == uuid.Nil {
		return errors.New("invalid menu item id")
	}
	return s.repo.Delete(ctx, id)
}
