package main

import (
	"context"

	"github.com/gofiber/fiber/v2"
)

func createTeam(c *fiber.Ctx) error {

	var team Team

	if err := c.BodyParser(&team); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	_, err := DB.Exec(context.Background(),
		"INSERT INTO teams (name, description) VALUES ($1,$2)",
		team.Name,
		team.Description,
	)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "team created",
	})
}

func getTeams(c *fiber.Ctx) error {

	rows, err := DB.Query(context.Background(),
		"SELECT id, name, description FROM teams")

	if err != nil {
		return c.Status(500).JSON(err)
	}

	defer rows.Close()

	var teams []Team

	for rows.Next() {

		var team Team

		rows.Scan(&team.ID, &team.Name, &team.Description)

		teams = append(teams, team)
	}

	return c.JSON(teams)
}