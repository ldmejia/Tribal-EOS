package main

import "github.com/gofiber/fiber/v2"

func setupRoutes(app *fiber.App) {

	api := app.Group("/api")

	api.Post("/teams", createTeam)

	api.Get("/teams", getTeams)
}