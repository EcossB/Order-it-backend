package httphandler

import (
	"net/http"
	"time"

	"order-it-backend/internal/domain"
	"order-it-backend/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type TenantHandler struct {
	service *service.TenantService
}

func NewTenantHandler(service *service.TenantService) *TenantHandler {
	return &TenantHandler{service: service}
}

type TenantRequest struct {
	Name string `json:"name" binding:"required"`
	Type string `json:"type" binding:"required"`
}

type TenantResponse struct {
	Id        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Type      string    `json:"type"`
	CreatedAt time.Time `json:"created_at"`
}

func (h *TenantHandler) Create(c *gin.Context) {
	var req TenantRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "JSON inválido o faltan campos obligatorios"})
		return
	}

	tenant := &domain.Tenant{
		Name: req.Name,
		Type: req.Type,
	}

	if err := h.service.CreateTenant(c.Request.Context(), tenant); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, TenantResponse{
		Id:        tenant.Id,
		Name:      tenant.Name,
		Type:      tenant.Type,
		CreatedAt: tenant.CreatedAt,
	})
}

func (h *TenantHandler) GetAll(c *gin.Context) {
	tenants, err := h.service.GetAllTenants(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var response []TenantResponse
	for _, t := range tenants {
		response = append(response, TenantResponse{
			Id:        t.Id,
			Name:      t.Name,
			Type:      t.Type,
			CreatedAt: t.CreatedAt,
		})
	}

	c.JSON(http.StatusOK, response)
}

func (h *TenantHandler) GetById(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido, debe ser un UUID"})
		return
	}

	tenant, err := h.service.GetTenantById(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, TenantResponse{
		Id:        tenant.Id,
		Name:      tenant.Name,
		Type:      tenant.Type,
		CreatedAt: tenant.CreatedAt,
	})
}

func (h *TenantHandler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido, debe ser un UUID"})
		return
	}

	var req TenantRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "JSON inválido o faltan campos obligatorios"})
		return
	}

	tenant := &domain.Tenant{
		Id:   id,
		Name: req.Name,
		Type: req.Type,
	}

	if err := h.service.UpdateTenant(c.Request.Context(), tenant); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Organización actualizada exitosamente"})
}

func (h *TenantHandler) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido, debe ser un UUID"})
		return
	}

	if err := h.service.DeleteTenant(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Organización eliminada exitosamente"})
}
