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

func createUser(c *fiber.Ctx) error {

	var user User

	if err := c.BodyParser(&user); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "invalid request",
		})
	}

	hashedPassword, err := hashPassword(user.Password)

	if err != nil {
		return c.Status(500).JSON(err)
	}

	_, err = DB.Exec(context.Background(),
		"INSERT INTO users (name, email, password, role) VALUES ($1,$2,$3,$4)",
		user.Name,
		user.Email,
		hashedPassword,
		user.Role,
	)

	if err != nil {
		return c.Status(500).JSON(err)
	}

	return c.JSON(fiber.Map{
		"message": "user created",
	})
}

func getUsers(c *fiber.Ctx) error {

	rows, err := DB.Query(context.Background(),
		"SELECT id, name, email, role FROM users")

	if err != nil {
		return c.Status(500).JSON(err)
	}

	defer rows.Close()

	var users []User

	for rows.Next() {

		var user User

		rows.Scan(
			&user.ID,
			&user.Name,
			&user.Email,
			&user.Role,
		)

		users = append(users, user)
	}

	return c.JSON(users)
}

func getUser(c *fiber.Ctx) error {

	id := c.Params("id")

	row := DB.QueryRow(context.Background(),
		"SELECT id, name, email, role FROM users WHERE id=$1", id)

	var user User

	err := row.Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Role,
	)

	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": "user not found",
		})
	}

	return c.JSON(user)
}

func deleteUser(c *fiber.Ctx) error {

	id := c.Params("id")

	_, err := DB.Exec(context.Background(),
		"DELETE FROM users WHERE id=$1", id)

	if err != nil {
		return c.Status(500).JSON(err)
	}

	return c.JSON(fiber.Map{
		"message": "user deleted",
	})
}

func login(c *fiber.Ctx) error {

	type LoginRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var body LoginRequest

	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON("invalid request")
	}

	row := DB.QueryRow(context.Background(),
		"SELECT id, password FROM users WHERE email=$1",
		body.Email,
	)

	var id int
	var hash string

	err := row.Scan(&id, &hash)

	if err != nil {
		return c.Status(401).JSON("invalid credentials")
	}

	if !checkPassword(body.Password, hash) {
		return c.Status(401).JSON("invalid credentials")
	}

	token, err := generateToken(id)

	if err != nil {
		return c.Status(500).JSON(err)
	}

	return c.JSON(fiber.Map{
		"token": token,
	})
}