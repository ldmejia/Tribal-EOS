package main

import "github.com/gofiber/fiber/v2"

func setupRoutes(app *fiber.App) {

	api := app.Group("/api")
	api.Post("/users", createUser)
	api.Post("/auth/login", login)

	protected := api.Group("/", authMiddleware)
	protected.Post("/teams", createTeam)
	protected.Get("/teams", getTeams)
	protected.Get("/users", getUsers)
	protected.Get("/users/:id", getUser)
	protected.Delete("/users/:id", deleteUser)
}