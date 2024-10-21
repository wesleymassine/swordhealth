package middleware

import (
	"os"

	fiber "github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

const bearerPrefix = "Bearer "

// Middleware to check if the user is authenticated
func AuthRequired(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Authorization header is missing",
		})
	}

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

	claims, err := ParseJWT(token)

	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid or expired token",
		})
	}

	// Attach user information to the context
	c.Locals("user_id", claims["user_id"])
	c.Locals("role", claims["role"])

	return c.Next()
}

// RoleMiddleware ensures the user has the appropriate role to access the route.
func RoleMiddleware(role string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userRole := c.Locals("role").(string)
		if userRole != role {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "Insufficient permissions",
			})
		}
		return c.Next()
	}
}

// ParseJWT parses and validates a JWT token, returning the claims if valid.
func ParseJWT(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrInvalidKey
		}
		return os.Getenv("JWT_SECRET"), nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		return claims, nil
	}

	return nil, err
}
