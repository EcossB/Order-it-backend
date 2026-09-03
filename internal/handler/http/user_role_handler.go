package httphandler

import (
	"net/http"
	"order-it-backend/internal/domain"
	"order-it-backend/internal/service"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UserRoleHandler struct {
	service *service.UserRoleService
}

func NewUserRoleHandler(service *service.UserRoleService) *UserRoleHandler {
	return &UserRoleHandler{service: service}
}

type UserRoleRequest struct {
	TenantId    uuid.UUID `json:"tenant_id" binding:"required"`
	Description string    `json:"description" binding:"required"`
}

type UserRoleResponse struct {
	Id          uuid.UUID `json:"id"`
	TenantId    uuid.UUID `json:"tenant_id"`
	Description string    `json:"description"`
	Status      int       `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
}

func (h *UserRoleHandler) CreateRole(c *gin.Context) {
	var req UserRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "JSON inválido o faltan campos por completar"})
		return
	}

	userRoleDomain := &domain.UserRole{
		TenantId:    req.TenantId,
		Description: req.Description,
	}

	userRoleDomain, err := h.service.CreateUserRole(c.Request.Context(), userRoleDomain)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	res := &UserRoleResponse{
		Id:          userRoleDomain.Id,
		TenantId:    userRoleDomain.TenantId,
		Description: userRoleDomain.Description,
		Status:      userRoleDomain.Status,
		CreatedAt:   userRoleDomain.CreatedAt,
	}

	c.JSON(http.StatusCreated, res)
}

func (h *UserRoleHandler) GetAllRoles(c *gin.Context) {
	tenantIdStr := c.Query("tenant_id")
	tenantId, err := uuid.Parse(tenantIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "el parámetro tenant_id es requerido y debe ser UUID"})
		return
	}

	userRoleDomains, err := h.service.GetAllUserRolesByTenantId(c.Request.Context(), tenantId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	userRolesResponse := make([]UserRoleResponse, 0)
	for _, domainRole := range userRoleDomains {
		userRolesResponse = append(userRolesResponse, UserRoleResponse{
			Id:          domainRole.Id,
			TenantId:    domainRole.TenantId,
			Description: domainRole.Description,
			Status:      domainRole.Status,
			CreatedAt:   domainRole.CreatedAt,
		})
	}

	c.JSON(http.StatusOK, userRolesResponse)
}

func (h *UserRoleHandler) GetRoleById(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id inválido"})
		return
	}

	domainRole, err := h.service.GetUserRoleById(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	res := &UserRoleResponse{
		Id:          domainRole.Id,
		TenantId:    domainRole.TenantId,
		Description: domainRole.Description,
		Status:      domainRole.Status,
		CreatedAt:   domainRole.CreatedAt,
	}

	c.JSON(http.StatusOK, res)
}

func (h *UserRoleHandler) UpdateRole(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "el id debe ser un UUID válido"})
		return
	}

	var req UserRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "JSON inválido o faltan campos por completar"})
		return
	}

	userRoleDomain := &domain.UserRole{
		Id:          id,
		TenantId:    req.TenantId,
		Description: req.Description,
	}

	err = h.service.UpdateUserRole(c.Request.Context(), userRoleDomain)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Rol actualizado exitosamente"})
}

func (h *UserRoleHandler) DeleteRole(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "el id debe ser un UUID válido"})
		return
	}

	err = h.service.DeleteUserRole(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Rol eliminado exitosamente"})
}
