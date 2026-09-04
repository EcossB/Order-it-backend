package domain

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

// DBTX es una interfaz que agrupa los métodos comunes entre pgxpool.Pool y pgx.Tx
// Esto permite a los repositorios ejecutar queries usando la transacción inyectada por el middleware,
// o caer de regreso al Pool de conexiones directo si no hay transacción.
type DBTX interface {
	Exec(context.Context, string, ...any) (pgconn.CommandTag, error)
	Query(context.Context, string, ...any) (pgx.Rows, error)
	QueryRow(context.Context, string, ...any) pgx.Row
}
