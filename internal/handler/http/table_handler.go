package httphandler

import (
	"net/http"
	"order-it-backend/internal/domain"
	"order-it-backend/internal/service"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type TableHandler struct {
	service *service.TableService
}

type TableRequest struct {
	TenantId    uuid.UUID `json:"tenant_id" binding:"required"`
	TableNumber string    `json:"table_number" binding:"required"`
}

type TableRequestUpdate struct {
	TenantId uuid.UUID `json:"tenant_id" binding:"required"`
	Status   string    `json:"status" binding:"required"`
}

type TableResponse struct {
	Id          uuid.UUID `json:"id" `
	TenantId    uuid.UUID `json:"tenant_id"`
	TableNumber string    `json:"table_number"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
}

func NewTableHandler(service *service.TableService) *TableHandler {
	return &TableHandler{service: service}
}

func (h *TableHandler) Create(c *gin.Context) {

	var tableRequest TableRequest
	if err := c.ShouldBindJSON(&tableRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "el json es invalido o faltan campos por completar."})
	}

	// mapeo de dto a domain

	tableDomain := &domain.Table{
		TenantId:    tableRequest.TenantId,
		TableNumber: tableRequest.TableNumber,
	}

	tableDomain, err := h.service.CreateTable(c.Request.Context(), tableDomain)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ocurrio un error intentando guardar la mesa."})
		return
	}

	//mapeo de domain a dto

	tableResponse := TableResponse{
		Id:          tableDomain.Id,
		TenantId:    tableDomain.TenantId,
		TableNumber: tableDomain.TableNumber,
		Status:      tableDomain.Status,
		CreatedAt:   tableDomain.CreatedAt,
	}

	c.JSON(http.StatusCreated, tableResponse)

}

func (h *TableHandler) GetAll(c *gin.Context) {

	tenantId, err := uuid.Parse(c.Query("tenant_id"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "el parámetro tenant_id es requerido y debe ser UUID (ej: ?tenant_id=...)"})
		return
	}

	tables, err := h.service.GetTableByTenantId(c.Request.Context(), tenantId)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var tableResponse = make([]TableResponse, 0)

	for _, table := range tables {
		tableResponse = append(tableResponse, TableResponse{
			Id:          table.Id,
			TenantId:    table.TenantId,
			TableNumber: table.TableNumber,
			Status:      table.Status,
			CreatedAt:   table.CreatedAt,
		})
	}

	c.JSON(http.StatusOK, tableResponse)

}

func (h *TableHandler) GetById(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de mesa inválido"})
		return
	}

	table, err := h.service.GetTableById(c.Request.Context(), id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	tableResponse := TableResponse{
		Id:          table.Id,
		TenantId:    table.TenantId,
		TableNumber: table.TableNumber,
		Status:      table.Status,
		CreatedAt:   table.CreatedAt,
	}

	c.JSON(http.StatusOK, tableResponse)

}

func (h *TableHandler) UpdateStatus(c *gin.Context) {

	id, err := uuid.Parse(c.Param("id"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de mesa inválido"})
		return
	}

	var tableRequestUpdate TableRequestUpdate

	if err := c.ShouldBindJSON(&tableRequestUpdate); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tableDomain := &domain.Table{
		Id:       id,
		Status:   tableRequestUpdate.Status,
		TenantId: tableRequestUpdate.TenantId,
	}

	err = h.service.UpdateTableStatus(c.Request.Context(), tableDomain)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "mesa actualizada exitosamente"})

}

func (h *TableHandler) Update(c *gin.Context) {

	id, err := uuid.Parse(c.Param("id"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de mesa inválido"})
		return
	}

	var tableRequest TableRequest

	if err := c.ShouldBindJSON(&tableRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tableDomain := &domain.Table{
		Id:          id,
		TenantId:    tableRequest.TenantId,
		TableNumber: tableRequest.TableNumber,
	}

	err = h.service.UpdateTable(c.Request.Context(), tableDomain)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Mesa actualizada exitosamente"})

}

func (h *TableHandler) Delete(c *gin.Context) {

	id, err := uuid.Parse(c.Param("id"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de mesa inválido"})
		return
	}

	if err := h.service.DeleteTableById(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Mesa eliminada exitosamente"})

}
