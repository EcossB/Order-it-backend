package httphandler

import (
	"net/http"
	"order-it-backend/internal/domain"
	"order-it-backend/internal/service"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type KitchenHandler struct {
	service *service.KitchenService
}

func NewKitchenHandler(service *service.KitchenService) *KitchenHandler {
	return &KitchenHandler{service: service}
}

type KitchenRequest struct {
	TenantId uuid.UUID `json:"tenant_id" binding:"required"`
	Name     string    `json:"name" binding:"required"`
}

type KitchenResponse struct {
	Id        uuid.UUID `json:"id"`
	TenantId  uuid.UUID `json:"tenant_id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

func (h *KitchenHandler) Create(c *gin.Context) {
	var req KitchenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "JSON inválido o faltan campos obligatorios"})
		return
	}

	kitchen := &domain.Kitchen{
		TenantId: req.TenantId,
		Name:     req.Name,
	}

	kitchen, err := h.service.CreateKitchen(c.Request.Context(), kitchen)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, KitchenResponse{
		Id:        kitchen.Id,
		TenantId:  kitchen.TenantId,
		Name:      kitchen.Name,
		CreatedAt: kitchen.CreatedAt,
	})
}

func (h *KitchenHandler) GetAll(c *gin.Context) {
	tenantIdStr := c.Query("tenant_id")
	tenantId, err := uuid.Parse(tenantIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "el parámetro tenant_id es requerido y debe ser UUID (ej: ?tenant_id=...)"})
		return
	}

	kitchens, err := h.service.GetAllKitchensByTenantId(c.Request.Context(), tenantId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var response = make([]KitchenResponse, 0)
	for _, k := range kitchens {
		response = append(response, KitchenResponse{
			Id:        k.Id,
			TenantId:  k.TenantId,
			Name:      k.Name,
			CreatedAt: k.CreatedAt,
		})
	}

	c.JSON(http.StatusOK, response)
}

func (h *KitchenHandler) GetById(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de cocina inválido"})
		return
	}

	kitchen, err := h.service.GetKitchenById(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, KitchenResponse{
		Id:        kitchen.Id,
		TenantId:  kitchen.TenantId,
		Name:      kitchen.Name,
		CreatedAt: kitchen.CreatedAt,
	})
}

func (h *KitchenHandler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de cocina inválido"})
		return
	}

	var req KitchenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "JSON inválido o faltan campos obligatorios"})
		return
	}

	kitchen := &domain.Kitchen{
		Id:       id,
		TenantId: req.TenantId,
		Name:     req.Name,
	}

	if err := h.service.UpdateKitchen(c.Request.Context(), kitchen); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Cocina actualizada exitosamente"})
}

func (h *KitchenHandler) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de cocina inválido"})
		return
	}

	if err := h.service.DeleteKitchen(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Cocina eliminada exitosamente"})
}
