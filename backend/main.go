package main

import (
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {

	godotenv.Load()

	connectDB()

	runSchema()

	app := fiber.New()

	setupRoutes(app)

	port := os.Getenv("PORT")

	app.Listen(":" + port)
}