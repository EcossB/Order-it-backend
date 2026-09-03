package httphandler

import (
	"github.com/gin-gonic/gin"
)

// Router agrupa todos los handlers y se encarga exclusivamente de definir las rutas
type Router struct {
	tenantHandler   *TenantHandler
	userRoleHandler *UserRoleHandler
	// Aquí inyectaremos los demás handlers en el futuro:

}

// NewRouter crea una nueva instancia del enrutador
func NewRouter(tenantHandler *TenantHandler, userRoleHandler *UserRoleHandler) *Router {
	return &Router{
		tenantHandler:   tenantHandler,
		userRoleHandler: userRoleHandler,
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

	// --- Futuras Rutas (ej. Users) ---
	// users := api.Group("/users")
	// {
	//     users.POST("", r.userHandler.Create)
	// }
}
