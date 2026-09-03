package main

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5"
)

func main() {
	// Conectar a la base de datos por defecto 'postgres' para crear la nueva BD
	defaultConnStr := "postgres://postgres:coss2003@localhost:5432/postgres"
	conn, err := pgx.Connect(context.Background(), defaultConnStr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error conectando a postgres: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())

	// Intentar crear la base de datos
	fmt.Println("Creando base de datos 'order_it'...")
	_, err = conn.Exec(context.Background(), "CREATE DATABASE order_it")
	if err != nil {
		fmt.Printf("Aviso: No se pudo crear (quizás ya existe): %v\n", err)
	} else {
		fmt.Println("Base de datos creada exitosamente.")
	}
	conn.Close(context.Background())

	// Conectar a la nueva base de datos para correr el esquema
	fmt.Println("Conectando a 'order_it'...")
	appConnStr := "postgres://postgres:coss2003@localhost:5432/order_it"
	appConn, err := pgx.Connect(context.Background(), appConnStr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error conectando a order_it: %v\n", err)
		os.Exit(1)
	}
	defer appConn.Close(context.Background())

	// Leer el archivo SQL
	schemaPath := "internal/models/schema.sql"
	schemaBytes, err := os.ReadFile(schemaPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error leyendo schema.sql: %v\n", err)
		os.Exit(1)
	}

	// Ejecutar el esquema
	fmt.Println("Ejecutando schema.sql...")
	_, err = appConn.Exec(context.Background(), string(schemaBytes))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error ejecutando el esquema: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("¡Esquema de base de datos inicializado correctamente!")
}
