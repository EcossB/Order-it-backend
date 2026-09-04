package httphandler

import (
	"net/http"
	"order-it-backend/internal/domain"
	"order-it-backend/internal/service"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type MenuItemHandler struct {
	service *service.MenuItemService
}

type MenuItemRequest struct {
	TenantId    uuid.UUID `json:"tenant_id" binding:"required"`
	KitchenId   uuid.UUID `json:"kitchen_id" binding:"required"`
	Name        string    `json:"name" binding:"required"`
	Description string    `json:"description"`
	Price       float64   `json:"price" binding:"required,min=0"`
	IsAvailable *bool     `json:"is_available"` // Pointer to handle false correctly
}

type MenuItemResponse struct {
	Id          uuid.UUID `json:"id"`
	TenantId    uuid.UUID `json:"tenant_id"`
	KitchenId   uuid.UUID `json:"kitchen_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       float64   `json:"price"`
	IsAvailable bool      `json:"is_available"`
	CreatedAt   time.Time `json:"created_at"`
}

func NewMenuItemHandler(service *service.MenuItemService) *MenuItemHandler {
	return &MenuItemHandler{service: service}
}

func (h *MenuItemHandler) Create(c *gin.Context) {
	var req MenuItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	isAvailable := true
	if req.IsAvailable != nil {
		isAvailable = *req.IsAvailable
	}

	item := &domain.MenuItem{
		TenantId:    req.TenantId,
		KitchenId:   req.KitchenId,
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		IsAvailable: isAvailable,
	}

	item, err := h.service.Create(c.Request.Context(), item)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, mapToMenuItemResponse(item))
}

func (h *MenuItemHandler) GetAll(c *gin.Context) {
	tenantId, err := uuid.Parse(c.Query("tenant_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "tenant_id parameter is required and must be a UUID"})
		return
	}

	items, err := h.service.GetAllByTenantId(c.Request.Context(), tenantId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var response = make([]MenuItemResponse, 0)
	for _, item := range items {
		response = append(response, mapToMenuItemResponse(item))
	}

	c.JSON(http.StatusOK, response)
}

func (h *MenuItemHandler) GetById(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid menu item id"})
		return
	}

	item, err := h.service.GetById(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, mapToMenuItemResponse(item))
}

func (h *MenuItemHandler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid menu item id"})
		return
	}

	var req MenuItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	isAvailable := true
	if req.IsAvailable != nil {
		isAvailable = *req.IsAvailable
	}

	item := &domain.MenuItem{
		Id:          id,
		TenantId:    req.TenantId,
		KitchenId:   req.KitchenId,
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		IsAvailable: isAvailable,
	}

	if err := h.service.Update(c.Request.Context(), item); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "menu item updated successfully"})
}

func (h *MenuItemHandler) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid menu item id"})
		return
	}

	if err := h.service.Delete(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "menu item deleted successfully"})
}

func mapToMenuItemResponse(item *domain.MenuItem) MenuItemResponse {
	return MenuItemResponse{
		Id:          item.Id,
		TenantId:    item.TenantId,
		KitchenId:   item.KitchenId,
		Name:        item.Name,
		Description: item.Description,
		Price:       item.Price,
		IsAvailable: item.IsAvailable,
		CreatedAt:   item.CreatedAt,
	}
}
