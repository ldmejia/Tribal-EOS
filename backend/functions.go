package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

var DB *pgxpool.Pool

func connectDB() {

    pool, err := pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))

    if err != nil {
        log.Fatal("database connection error:", err)
    }

    DB = pool

    fmt.Println("database connected")
}

func runSchema() {

	schema, err := os.ReadFile("database/schema.sql")

	if err != nil {
		log.Fatal("error reading schema:", err)
	}

	_, err = DB.Exec(context.Background(), string(schema))

	if err != nil {
		log.Fatal("error executing schema:", err)
	}

	fmt.Println("database schema loaded")
}