package handlers

import "github.com/gofiber/fiber/v2"

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
