// Package middleware
package middleware

import (
	"playcorner-be/internal/auth"
	"playcorner-be/internal/models"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func AuthMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(models.ErrorResponse{
				Code:   401,
				Status: "UNAUTHORIZED",
				Data:   models.ErrorData{ErrorMsg: "Missing authorization header"},
			})
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return c.Status(fiber.StatusUnauthorized).JSON(models.ErrorResponse{
				Code:   401,
				Status: "UNAUTHORIZED",
				Data:   models.ErrorData{ErrorMsg: "Invalid authorization header format"},
			})
		}

		tokenString := parts[1]
		claims, err := auth.ValidateToken(tokenString)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(models.ErrorResponse{
				Code:   401,
				Status: "UNAUTHORIZED",
				Data:   models.ErrorData{ErrorMsg: "Invalid or expired token"},
			})
		}

		// Menyimpan user ID dari token ke dalam context untuk digunakan oleh handler selanjutnya
		c.Locals("userID", claims.UserID)
		return c.Next()
	}
}
