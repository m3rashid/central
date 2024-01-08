package controllers

import "github.com/gofiber/fiber/v2"

func Discovery(ctx *fiber.Ctx) error {
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"/all":  "Get all Contacts",
		"/{id}": "Get Contact by ID",
	})
}
