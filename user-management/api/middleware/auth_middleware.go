package middleware

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/wesleymassine/swordhealth/user-management/api/security"
)

const bearerPrefix = "Bearer "

// JWTMiddleware checks for a valid JWT token in the Authorization header.
func JWTMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Authorization header is missing",
			})
		}

		fmt.Println("Authorization", authHeader)

		var token string
		if len(authHeader) > len(bearerPrefix) && authHeader[:len(bearerPrefix)] == bearerPrefix {
			token = authHeader[len(bearerPrefix):]

			if token == "" {
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
					"error": "Token is empty",
				})
			}

		} else {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid authorization header",
			})
		}

		if len(token) == len("null") {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Please authenticate",
			})
		}

		fmt.Println("token", token)

		claims, err := security.ParseJWT(token)

		fmt.Println("ParseJWT ERROR", claims)

		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid or expired token",
			})
		}

		println("SET STORE========:", claims["user_id"], claims["role"])

		// Store the user ID and role in the context
		c.Locals("user_id", claims["user_id"])
		c.Locals("role", claims["role"])
		return c.Next()
	}
}

// RoleMiddleware ensures the user has the appropriate role to access the route.
func RoleMiddleware(role string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userRole := c.Locals("role")
		if userRole != role {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "Insufficient permissions",
			})
		}
		return c.Next()
	}
}
