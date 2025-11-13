package routes

import (
	"golang-skeleton/app/entities/dtos/responses"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	api := app.Group("/api")

	api.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(responses.NewSuccessResponse("API is running", nil, fiber.StatusOK))
	})
}
