package middleware

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Clave privada para evitar colisiones en el contexto
type txKey struct{}

// TransactionMiddleware crea una transacción para cada petición HTTP
func TransactionMiddleware(dbPool *pgxpool.Pool) gin.HandlerFunc {

	return func(c *gin.Context) {
		// Iniciar transacción
		tx, err := dbPool.Begin(c.Request.Context())
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "no se pudo iniciar la transacción de base de datos"})
			return
		}

		// Intentar obtener el tenant_id de los headers (ej. X-Tenant-ID)
		tenantID := c.GetHeader("X-Tenant-ID")

		// Si viene el tenant_id, configuramos el RLS para la transacción actual
		if tenantID != "" {
			_, err = tx.Exec(c.Request.Context(), "SELECT set_config('app.current_tenant_id', $1, true)", tenantID)
			if err != nil {
				tx.Rollback(c.Request.Context())
				c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "tenant_id inválido o acceso denegado por RLS"})
				return
			}
		}

		// Guardar la transacción en el contexto subyacente
		ctx := context.WithValue(c.Request.Context(), txKey{}, tx)
		c.Request = c.Request.WithContext(ctx)

		// Continuar con el resto de la ejecución (handlers, otros middlewares)
		c.Next()

		// Al finalizar, si hubo errores en Gin, hacemos rollback, si no, commit
		if len(c.Errors) > 0 || c.Writer.Status() >= 400 {
			tx.Rollback(c.Request.Context())
		} else {
			tx.Commit(c.Request.Context())
		}
	}
}

// GetTx extrae la transacción del contexto, retornando false si no existe
func GetTx(ctx context.Context) (pgx.Tx, bool) {
	tx, ok := ctx.Value(txKey{}).(pgx.Tx)
	return tx, ok
}
