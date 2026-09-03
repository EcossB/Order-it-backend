package httphandler

import (
	"github.com/gin-gonic/gin"
)

// Router agrupa todos los handlers y se encarga exclusivamente de definir las rutas
type Router struct {
	tenantHandler   *TenantHandler
	userRoleHandler *UserRoleHandler
	kitchenHandler  *KitchenHandler
}

// NewRouter crea una nueva instancia del enrutador
func NewRouter(tenantHandler *TenantHandler, userRoleHandler *UserRoleHandler, kitchenHandler *KitchenHandler) *Router {
	return &Router{
		tenantHandler:   tenantHandler,
		userRoleHandler: userRoleHandler,
		kitchenHandler:  kitchenHandler,
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
}
