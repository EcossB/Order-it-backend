package httphandler

import (
	"net/http"
	"order-it-backend/internal/domain"
	"order-it-backend/internal/service"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UserHandler struct {
	service *service.UserService
}

type UserRequest struct {
	TenantId  uuid.UUID  `json:"tenant_id" binding:"required"`
	Role      uuid.UUID  `json:"role" binding:"required"`
	KitchenId *uuid.UUID `json:"kitchen_id"`
	Name      string     `json:"name" binding:"required"`
	Email     *string    `json:"email"`
	Password  *string    `json:"password"`
	Pin       *string    `json:"pin"`
}

type UserResponse struct {
	Id        uuid.UUID  `json:"id"`
	TenantId  uuid.UUID  `json:"tenant_id"`
	Role      uuid.UUID  `json:"role"`
	KitchenId *uuid.UUID `json:"kitchen_id"`
	Name      string     `json:"name"`
	Email     *string    `json:"email"`
	CreatedAt time.Time  `json:"created_at"`
}

func NewUserHandler(service *service.UserService) *UserHandler {
	return &UserHandler{service: service}
}

func (h *UserHandler) Create(c *gin.Context) {
	var req UserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := &domain.User{
		TenantId:  req.TenantId,
		Role:      req.Role,
		KitchenId: req.KitchenId,
		Name:      req.Name,
		Email:     req.Email,
	}

	user, err := h.service.Create(c.Request.Context(), user, req.Password, req.Pin)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, mapToUserResponse(user))
}

func (h *UserHandler) GetAll(c *gin.Context) {
	tenantId, err := uuid.Parse(c.Query("tenant_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "tenant_id parameter is required and must be a UUID"})
		return
	}

	users, err := h.service.GetAllByTenantId(c.Request.Context(), tenantId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var response = make([]UserResponse, 0)
	for _, user := range users {
		response = append(response, mapToUserResponse(user))
	}

	c.JSON(http.StatusOK, response)
}

func (h *UserHandler) GetById(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	user, err := h.service.GetById(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, mapToUserResponse(user))
}

func (h *UserHandler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	var req UserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := &domain.User{
		Id:        id,
		TenantId:  req.TenantId,
		Role:      req.Role,
		KitchenId: req.KitchenId,
		Name:      req.Name,
		Email:     req.Email,
	}

	if err := h.service.Update(c.Request.Context(), user, req.Password, req.Pin); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "user updated successfully"})
}

func (h *UserHandler) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	if err := h.service.Delete(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "user deleted successfully"})
}

func mapToUserResponse(user *domain.User) UserResponse {
	return UserResponse{
		Id:        user.Id,
		TenantId:  user.TenantId,
		Role:      user.Role,
		KitchenId: user.KitchenId,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
	}
}
