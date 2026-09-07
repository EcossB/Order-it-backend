package httphandler

import (
	"net/http"
	"order-it-backend/internal/domain"
	"order-it-backend/internal/service"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type OrderHandler struct {
	service *service.OrderService
}

func NewOrderHandler(service *service.OrderService) *OrderHandler {
	return &OrderHandler{service: service}
}

// Request Payload
type CreateOrderRequest struct {
	TenantId uuid.UUID                `json:"tenant_id" binding:"required"`
	TableId  *uuid.UUID               `json:"table_id"`
	WaiterId *uuid.UUID               `json:"waiter_id"`
	Items    []CreateOrderItemRequest `json:"items" binding:"required,gt=0"`
}

type CreateOrderItemRequest struct {
	MenuItemId uuid.UUID `json:"menu_item_id" binding:"required"`
	Quantity   int       `json:"quantity" binding:"required,gt=0"`
	Notes      string    `json:"notes"`
}

type UpdateOrderStatusRequest struct {
	TenantId uuid.UUID `json:"tenant_id" binding:"required"`
	Status   string    `json:"status" binding:"required"`
}

type UpdateOrderItemStatusRequest struct {
	TenantId uuid.UUID `json:"tenant_id" binding:"required"`
	Status   string    `json:"status" binding:"required"`
}

// Response Payload
type OrderResponse struct {
	Id        uuid.UUID           `json:"id"`
	TenantId  uuid.UUID           `json:"tenant_id"`
	TableId   *uuid.UUID          `json:"table_id"`
	WaiterId  *uuid.UUID          `json:"waiter_id"`
	Status    string              `json:"status"`
	CreatedAt time.Time           `json:"created_at"`
	Items     []OrderItemResponse `json:"items,omitempty"`
}

type OrderItemResponse struct {
	Id         uuid.UUID `json:"id"`
	MenuItemId uuid.UUID `json:"menu_item_id"`
	KitchenId  uuid.UUID `json:"kitchen_id"`
	Quantity   int       `json:"quantity"`
	Status     string    `json:"status"`
	Notes      string    `json:"notes"`
}

func (h *OrderHandler) Create(c *gin.Context) {
	var req CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	order := &domain.Order{
		TenantId: req.TenantId,
		TableId:  req.TableId,
		WaiterId: req.WaiterId,
	}

	var items []*domain.OrderItem
	for _, reqItem := range req.Items {
		items = append(items, &domain.OrderItem{
			MenuItemId: reqItem.MenuItemId,
			Quantity:   reqItem.Quantity,
			Notes:      reqItem.Notes,
		})
	}

	createdOrder, createdItems, err := h.service.CreateOrderWithItems(c.Request.Context(), order, items)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	res := mapToOrderResponse(createdOrder, createdItems)
	c.JSON(http.StatusCreated, res)
}

func (h *OrderHandler) GetAll(c *gin.Context) {
	tenantId, err := uuid.Parse(c.Query("tenant_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "tenant_id parameter is required and must be a UUID"})
		return
	}

	orders, err := h.service.GetAllOrdersByTenantId(c.Request.Context(), tenantId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var response = make([]OrderResponse, 0)
	for _, order := range orders {
		response = append(response, mapToOrderResponse(order, nil))
	}

	c.JSON(http.StatusOK, response)
}

func (h *OrderHandler) GetById(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid order id"})
		return
	}

	order, items, err := h.service.GetOrderById(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, mapToOrderResponse(order, items))
}

func (h *OrderHandler) UpdateStatus(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid order id"})
		return
	}

	var req UpdateOrderStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.UpdateOrderStatus(c.Request.Context(), id, req.TenantId, req.Status); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "order status updated successfully"})
}

func (h *OrderHandler) UpdateItemStatus(c *gin.Context) {
	itemId, err := uuid.Parse(c.Param("item_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid order item id"})
		return
	}

	var req UpdateOrderItemStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.UpdateOrderItemStatus(c.Request.Context(), itemId, req.TenantId, req.Status); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "order item status updated successfully"})
}

func (h *OrderHandler) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid order id"})
		return
	}

	if err := h.service.DeleteOrder(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "order deleted successfully"})
}

func mapToOrderResponse(order *domain.Order, items []*domain.OrderItem) OrderResponse {
	res := OrderResponse{
		Id:        order.Id,
		TenantId:  order.TenantId,
		TableId:   order.TableId,
		WaiterId:  order.WaiterId,
		Status:    order.Status,
		CreatedAt: order.CreatedAt,
	}

	if items != nil {
		res.Items = make([]OrderItemResponse, 0, len(items))
		for _, item := range items {
			res.Items = append(res.Items, OrderItemResponse{
				Id:         item.Id,
				MenuItemId: item.MenuItemId,
				KitchenId:  item.KitchenId,
				Quantity:   item.Quantity,
				Status:     item.Status,
				Notes:      item.Notes,
			})
		}
	}

	return res
}
