package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/wesleymassine/swordhealth/user-management/api/middleware"
	"github.com/wesleymassine/swordhealth/user-management/api/security"
)

func TestJWTMiddleware(t *testing.T) {
	app := fiber.New()
	app.Use(middleware.JWTMiddleware())

	app.Get("/protected", func(c *fiber.Ctx) error {
		return c.SendStatus(fiber.StatusOK)
	})

	// Create a valid JWT
	token, _ := security.GenerateJWT(1, "admin")

	// Make request with JWT
	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	resp, _ := app.Test(req)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestRoleMiddleware(t *testing.T) {
	app := fiber.New()
	app.Use(middleware.JWTMiddleware())
	app.Use(middleware.RoleMiddleware("admin"))

	app.Get("/admin-only", func(c *fiber.Ctx) error {
		return c.SendStatus(fiber.StatusOK)
	})

	// Create a valid JWT for non-admin user
	token, _ := security.GenerateJWT(1, "user")

	// Make request with JWT
	req := httptest.NewRequest(http.MethodGet, "/admin-only", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	resp, _ := app.Test(req)

	assert.Equal(t, http.StatusForbidden, resp.StatusCode)
}
