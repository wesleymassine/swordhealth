package middleware

import (
	fiber "github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

const bearerPrefix = "Bearer "

var jwtSecret = "ace353e1c2dd28f9fa8c40f3687f943f7a4c0576dedc702fc049f7f98f06467a" // Change to a strong key in production TODO

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
		return jwtSecret, nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		return claims, nil
	}

	return nil, err
}

func ApiKeyMiddleware(c *fiber.Ctx) error {
	// If API key is not in the Authorization header, check the query parameter
	apiKey := c.Get("x-api-key")

	// Check if the API key is valid
	if apiKey != "swordhealth" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid or missing API key",
		})
	}

	// If API key is valid, proceed to the nex handler
	return c.Next()
}
