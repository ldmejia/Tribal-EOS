package main

import (
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func hashPassword(password string) (string, error) {

	bytes, err := bcrypt.GenerateFromPassword(
		[]byte(password),
		14,
	)

	return string(bytes), err
}

func checkPassword(password string, hash string) bool {

	err := bcrypt.CompareHashAndPassword(
		[]byte(hash),
		[]byte(password),
	)

	return err == nil
}

func generateToken(userID int) (string, error) {

	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * 72).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(os.Getenv("JWT_SECRET")))
}

func authMiddleware(c *fiber.Ctx) error {

	tokenString := c.Get("Authorization")

	if tokenString == "" {
		return c.Status(401).JSON("missing token")
	}

	tokenString = tokenString[len("Bearer "):]

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil || !token.Valid {
		return c.Status(401).JSON("invalid token")
	}

	claims := token.Claims.(jwt.MapClaims)

	c.Locals("user_id", int(claims["user_id"].(float64)))

	return c.Next()
}