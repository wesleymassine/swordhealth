package api

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/wesleymassine/swordhealth/user-management/api/security"
)

// Login authenticates the user and returns a JWT token.
func (h *UserHandler) Login(c *fiber.Ctx) error {
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

	user, err := h.userUsecase.GetUserByEmail(c.Context(), loginRequest.Email)

	if err != nil {
		if user == nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "User not registered",
			})
		}

		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid email or password",
		})
	}

	// Verify password
	err = security.CheckPassword(user.Password, loginRequest.Password)

	// TODO: PASSWORD
	fmt.Println("TODO ERROR:", err)
	fmt.Println("Stored hashed password:", user.Password)
	fmt.Println("Provided login password:", loginRequest.Password)
	hashPassword, err := security.HashPassword(loginRequest.Password)
	fmt.Println("hash generate:", hashPassword)

	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid email or password",
		})
	}

	// Generate JWT
	token, err := security.GenerateJWT(user.ID, user.Role)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to generate token",
		})
	}

	return c.JSON(fiber.Map{
		"token": token,
	})
}
