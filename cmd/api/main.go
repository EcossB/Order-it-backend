package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"

	httphandler "order-it-backend/internal/handler/http"
	"order-it-backend/internal/middleware"
	"order-it-backend/internal/repository/postgres"
	"order-it-backend/internal/service"
)

func main() {
	// 1. Configuración de Base de Datos
	// Usamos pgxpool para manejar un pool de conexiones concurrentes, ideal para WebSockets.
	dbURL := "postgres://postgres:coss2003@localhost:5432/order_it"

	ctx := context.Background()
	dbPool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error conectando a la base de datos: %v\n", err)
		os.Exit(1)
	}
	defer dbPool.Close()

	fmt.Println("Conectado a la base de datos PostgreSQL exitosamente.")

	// 2. Inyección de Dependencias (Dependency Injection)
	// Aquí conectamos todas las capas de la Clean Architecture.

	// Capa 1: Repositorios (Hablan con la BD)
	tenantRepo := postgres.NewTenantRepository(dbPool)
	userRoleRepo := postgres.NewUserRoleRepository(dbPool)
	kitchenRepo := postgres.NewKitchenRepository(dbPool)

	// Capa 2: Servicios (Lógica de Negocio)
	tenantService := service.NewTenantService(tenantRepo)
	userRoleService := service.NewUserRoleService(userRoleRepo, tenantRepo)
	kitchenService := service.NewKitchenService(kitchenRepo, tenantRepo)

	// Capa 3: Handlers (Reciben peticiones HTTP y hablan con el Servicio)
	tenantHandler := httphandler.NewTenantHandler(tenantService)
	userRoleHandler := httphandler.NewUserRoleHandler(userRoleService)
	kitchenHandler := httphandler.NewKitchenHandler(kitchenService)

	// Capa 4: Enrutador Centralizado
	appRouter := httphandler.NewRouter(tenantHandler, userRoleHandler, kitchenHandler)

	// 3. Configuración de Rutas (Router) con GIN
	router := gin.Default()

	// Ruta de ejemplo (Health Check)
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Order It! API is running",
		})
	})

	// Grupo de rutas principal de la API
	api := router.Group("/api")
	api.Use(middleware.TransactionMiddleware(dbPool))
	{
		// Registramos todas las rutas a través de nuestro archivo router.go
		appRouter.RegisterRoutes(api)
	}

	// 4. Iniciar Servidor HTTP
	port := ":8080"
	fmt.Printf("Servidor backend corriendo en http://localhost%s...\n", port)

	if err := router.Run(port); err != nil {
		log.Fatalf("Error iniciando el servidor: %v\n", err)
	}
}
