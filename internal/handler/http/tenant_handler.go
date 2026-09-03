package httphandler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"order-it-backend/internal/domain"
	"order-it-backend/internal/service"
)

// TenantHandler maneja las peticiones HTTP (las entradas y salidas JSON) para los Tenants
type TenantHandler struct {
	service *service.TenantService
}

// NewTenantHandler inicializa el handler inyectando el servicio
func NewTenantHandler(service *service.TenantService) *TenantHandler {
	return &TenantHandler{
		service: service,
	}
}

// CreateTenantRequest define exactamente qué campos esperamos recibir en el JSON.
// Esto nos protege de recibir basura en la petición.
type CreateTenantRequest struct {
	Name string `json:"name" binding:"required"`
	Type string `json:"type" binding:"required"`
}

func (h *TenantHandler) Create(c *gin.Context) {
	var req CreateTenantRequest

	// 1. "Parsear" y Validar el JSON: Gin se encarga de revisar que vengan los campos required.
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "JSON inválido o faltan campos obligatorios"})
		return
	}

	// 2. Mapear del "Request" a nuestro Struct de Dominio
	tenant := &domain.Tenant{
		Name: req.Name,
		Type: req.Type,
	}

	// 3. Pasar el control a la capa de Negocio (Servicio)
	// Extraemos el context nativo con c.Request.Context() para mandarlo a la base de datos
	if err := h.service.CreateTenant(c.Request.Context(), tenant); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 4. Devolver la respuesta exitosa (Gin automáticamente lo convierte a JSON)
	c.JSON(http.StatusCreated, tenant)
}

func (h *TenantHandler) GetAll(c *gin.Context) {
	tenants, err := h.service.GetAllTenants(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener las organizaciones"})
		return
	}

	c.JSON(http.StatusOK, tenants)
}

func (h *TenantHandler) GetById(c *gin.Context) {
	// Extraer el ID de la URL (ej. /api/tenants/5)
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "El ID debe ser un número entero válido"})
		return
	}

	tenant, err := h.service.GetTenantById(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, tenant)
}

func (h *TenantHandler) Update(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "El ID debe ser un número entero válido"})
		return
	}

	var req CreateTenantRequest
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
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "El ID debe ser un número entero válido"})
		return
	}

	if err := h.service.DeleteTenant(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Organización eliminada exitosamente"})
}
