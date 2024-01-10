package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/m3rashid/central/campaigns/components"
)

func ConsentSuccessScreen(ctx *fiber.Ctx) error {
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Authorization successful",
	})
}

func ConsentFailureScreen(ctx *fiber.Ctx) error {
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": false,
		"message": "Authorization failed",
	})
}

func RenderLoginScreen(ctx *fiber.Ctx) error {
	ctx.Set("Content-Type", "text/html")
	return components.LoginScreen().Render(ctx.Context(), ctx.Response().BodyWriter())
}
