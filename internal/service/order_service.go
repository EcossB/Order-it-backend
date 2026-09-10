package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"order-it-backend/internal/websockets"

	"order-it-backend/internal/domain"

	"github.com/google/uuid"
)

type OrderService struct {
	orderRepo    domain.OrderRepository
	itemRepo     domain.OrderItemRepository
	menuItemRepo domain.MenuItemRepository
	tableRepo    domain.TableRepository
	tenantRepo   domain.TenantRepository
	hub          *websockets.Hub
}

func NewOrderService(
	orderRepo domain.OrderRepository,
	itemRepo domain.OrderItemRepository,
	menuItemRepo domain.MenuItemRepository,
	tableRepo domain.TableRepository,
	tenantRepo domain.TenantRepository,
	hub *websockets.Hub,
) *OrderService {
	return &OrderService{
		orderRepo:    orderRepo,
		itemRepo:     itemRepo,
		menuItemRepo: menuItemRepo,
		tableRepo:    tableRepo,
		tenantRepo:   tenantRepo,
		hub:          hub,
	}
}

func (s *OrderService) CreateOrderWithItems(ctx context.Context, order *domain.Order, items []*domain.OrderItem) (*domain.Order, []*domain.OrderItem, error) {
	if order.TenantId == uuid.Nil {
		return nil, nil, errors.New("invalid tenant id")
	}
	if len(items) == 0 {
		return nil, nil, errors.New("order must have at least one item")
	}

	// 1. Validate Tenant
	if _, err := s.tenantRepo.GetById(ctx, order.TenantId); err != nil {
		return nil, nil, errors.New("tenant does not exist")
	}

	// 2. Validate Table if provided
	if order.TableId != nil {
		table, err := s.tableRepo.GetById(ctx, *order.TableId)
		if err != nil {
			return nil, nil, errors.New("table does not exist")
		}
		if table.TenantId != order.TenantId {
			return nil, nil, errors.New("table does not belong to this tenant")
		}
	}

	// 3. Set pending status and Create Order Master Record
	order.Status = "PENDING"
	if err := s.orderRepo.Create(ctx, order); err != nil {
		return nil, nil, err
	}

	// 4. Fetch all Menu Items in a single query (Fix N+1)
	var menuItemIds []uuid.UUID
	for _, item := range items {
		menuItemIds = append(menuItemIds, item.MenuItemId)
	}

	menuItems, err := s.menuItemRepo.GetByIds(ctx, menuItemIds)
	if err != nil {
		return nil, nil, err
	}

	menuItemMap := make(map[uuid.UUID]*domain.MenuItem)
	for _, mi := range menuItems {
		menuItemMap[mi.Id] = mi
	}

	var createdItems []*domain.OrderItem

	// 5. Validate and Prepare Order Items
	for _, item := range items {
		if item.Quantity <= 0 {
			return nil, nil, errors.New("item quantity must be greater than zero")
		}

		menuItem, exists := menuItemMap[item.MenuItemId]
		if !exists {
			return nil, nil, errors.New("menu item does not exist")
		}
		if menuItem.TenantId != order.TenantId {
			return nil, nil, errors.New("menu item does not belong to this tenant")
		}
		if !menuItem.IsAvailable {
			return nil, nil, errors.New("menu item is not available: " + menuItem.Name)
		}

		item.TenantId = order.TenantId
		item.OrderId = order.Id
		item.KitchenId = menuItem.KitchenId
		item.Status = "PENDING"
		createdItems = append(createdItems, item)
	}

	// 6. Bulk Insert Order Items (Fix N+1)
	if err := s.itemRepo.CreateBulk(ctx, createdItems); err != nil {
		return nil, nil, err
	}

	// 7. avisar a la cocina en tiempo real que se creo una orden.

	for _, item := range createdItems {
		topic := fmt.Sprintf("tenant:%s:kitchen:%s", item.TenantId, item.KitchenId)

		//armamos un json con lo que necesita ver el chef usando un mapa anonimo.
		eventData := map[string]interface{}{
			"type":         "NEW_ORDER_ITEM",
			"order_id":     item.OrderId,
			"item_id":      item.Id,
			"menu_item_id": item.MenuItemId,
			"quantity":     item.Quantity,
			"notes":        item.Notes,
		}

		//convertimos el mapa a bytes (JSON)
		jsonData, _ := json.Marshal(eventData)

		// Disparamos el mensaje al HUB
		s.hub.Broadcast(topic, jsonData)

	}

	return order, createdItems, nil
}

func (s *OrderService) GetOrderById(ctx context.Context, id uuid.UUID) (*domain.Order, []*domain.OrderItem, error) {
	if id == uuid.Nil {
		return nil, nil, errors.New("invalid order id")
	}

	order, err := s.orderRepo.GetById(ctx, id)
	if err != nil {
		return nil, nil, err
	}

	items, err := s.itemRepo.GetAllByOrderId(ctx, id)
	if err != nil {
		return nil, nil, err
	}

	return order, items, nil
}

func (s *OrderService) GetAllOrdersByTenantId(ctx context.Context, tenantId uuid.UUID) ([]*domain.Order, error) {
	if tenantId == uuid.Nil {
		return nil, errors.New("invalid tenant id")
	}
	return s.orderRepo.GetAllByTenantId(ctx, tenantId)
}

func (s *OrderService) UpdateOrderStatus(ctx context.Context, orderId uuid.UUID, tenantId uuid.UUID, status string) error {
	order, err := s.orderRepo.GetById(ctx, orderId)
	if err != nil {
		return err
	}
	if order.TenantId != tenantId {
		return errors.New("order does not belong to this tenant")
	}

	order.Status = status
	return s.orderRepo.Update(ctx, order)
}

func (s *OrderService) UpdateOrderItemStatus(ctx context.Context, itemId uuid.UUID, tenantId uuid.UUID, status string) error {
	item, err := s.itemRepo.GetById(ctx, itemId)
	if err != nil {
		return err
	}
	if item.TenantId != tenantId {
		return errors.New("order item does not belong to this tenant")
	}

	item.Status = status
	return s.itemRepo.Update(ctx, item)
}

func (s *OrderService) DeleteOrder(ctx context.Context, orderId uuid.UUID) error {
	return s.orderRepo.Delete(ctx, orderId)
}
