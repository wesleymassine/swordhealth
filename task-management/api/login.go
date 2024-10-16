package api

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)


func Login(c *fiber.Ctx) error {
	//TODO: REQUEST and filds
	var loginRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.BodyParser(&loginRequest); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid input",
		})
	}

	fmt.Println("email", loginRequest.Email)
	fmt.Println("password", loginRequest.Password)

	if loginRequest.Email == "manager@gmail.com" && loginRequest.Password == "manager" {
		// Create JWT token
		claims := jwt.MapClaims{
			"user_id": 1,
			"role":    "manager",
			"exp":     time.Now().Add(time.Hour * 24).Unix(),
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, err := token.SignedString([]byte(jwtSecret))
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not login"})
		}

		return c.JSON(fiber.Map{"token": tokenString})
	}

	return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid credentials"})
}
