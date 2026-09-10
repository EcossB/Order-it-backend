package httphandler

import (
	"order-it-backend/internal/websockets"

	"github.com/gin-gonic/gin"
)

// Router agrupa todos los handlers y se encarga exclusivamente de definir las rutas
type Router struct {
	tenantHandler   *TenantHandler
	userRoleHandler *UserRoleHandler
	kitchenHandler  *KitchenHandler
	tableHandler    *TableHandler
	menuItemHandler *MenuItemHandler
	userHandler     *UserHandler
	orderHandler    *OrderHandler
	hub             *websockets.Hub
}

// NewRouter crea una nueva instancia del enrutador
func NewRouter(tenantHandler *TenantHandler,
	userRoleHandler *UserRoleHandler,
	kitchenHandler *KitchenHandler,
	tableHandler *TableHandler,
	menuItemHandler *MenuItemHandler,
	userHandler *UserHandler,
	orderHandler *OrderHandler,
	hub *websockets.Hub) *Router {
	return &Router{
		tenantHandler:   tenantHandler,
		userRoleHandler: userRoleHandler,
		kitchenHandler:  kitchenHandler,
		tableHandler:    tableHandler,
		menuItemHandler: menuItemHandler,
		userHandler:     userHandler,
		orderHandler:    orderHandler,
		hub:             hub,
	}
}

// RegisterRoutes mapea las URLs a las funciones correspondientes de cada Handler
func (r *Router) RegisterRoutes(api *gin.RouterGroup) {

	// --- Rutas de Tenants (Organizaciones) ---
	tenants := api.Group("/tenants")
	{
		tenants.POST("", r.tenantHandler.Create)
		tenants.GET("", r.tenantHandler.GetAll)
		tenants.GET("/:id", r.tenantHandler.GetById)
		tenants.PUT("/:id", r.tenantHandler.Update)
		tenants.DELETE("/:id", r.tenantHandler.Delete)
	}

	roles := api.Group("/roles")
	{
		roles.POST("", r.userRoleHandler.CreateRole)
		roles.GET("", r.userRoleHandler.GetAllRoles)
		roles.GET("/:id", r.userRoleHandler.GetRoleById)
		roles.PUT("/:id", r.userRoleHandler.UpdateRole)
		roles.DELETE("/:id", r.userRoleHandler.DeleteRole)

	}

	kitchens := api.Group("/kitchens")
	{
		kitchens.POST("", r.kitchenHandler.Create)
		kitchens.GET("", r.kitchenHandler.GetAll)
		kitchens.GET("/:id", r.kitchenHandler.GetById)
		kitchens.PUT("/:id", r.kitchenHandler.Update)
		kitchens.DELETE("/:id", r.kitchenHandler.Delete)
	}

	tables := api.Group("/tables")
	{
		tables.POST("", r.tableHandler.Create)
		tables.GET("", r.tableHandler.GetAll)
		tables.GET("/:id", r.tableHandler.GetById)
		tables.PUT("/:id", r.tableHandler.Update)
		tables.PATCH("/:id/status", r.tableHandler.UpdateStatus)
		tables.DELETE("/:id", r.tableHandler.Delete)
	}

	menuItems := api.Group("/menu-items")
	{
		menuItems.POST("", r.menuItemHandler.Create)
		menuItems.GET("", r.menuItemHandler.GetAll)
		menuItems.GET("/:id", r.menuItemHandler.GetById)
		menuItems.PUT("/:id", r.menuItemHandler.Update)
		menuItems.DELETE("/:id", r.menuItemHandler.Delete)
	}

	users := api.Group("/users")
	{
		users.POST("", r.userHandler.Create)
		users.GET("", r.userHandler.GetAll)
		users.GET("/:id", r.userHandler.GetById)
		users.PUT("/:id", r.userHandler.Update)
		users.DELETE("/:id", r.userHandler.Delete)
	}

	orders := api.Group("/orders")
	{
		orders.POST("", r.orderHandler.Create)
		orders.GET("", r.orderHandler.GetAll)
		orders.GET("/:id", r.orderHandler.GetById)
		orders.PATCH("/:id/status", r.orderHandler.UpdateStatus)
		orders.PATCH("/items/:item_id/status", r.orderHandler.UpdateItemStatus)
		orders.DELETE("/:id", r.orderHandler.Delete)
	}

	api.GET("/ws", websockets.ServeWS(r.hub))
}
