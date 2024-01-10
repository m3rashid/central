package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/m3rashid/central/auth/components"
	"github.com/m3rashid/central/auth/models"
)

func RenderLogoutScreen(ctx *fiber.Ctx) error {
	// client, _, err := getClient(ctx)
	ctx.Set("Content-Type", "text/html")
	return components.LogoutScreen(models.Client{}).Render(ctx.Context(), ctx.Response().BodyWriter())
}

func HanldeLogout(ctx *fiber.Ctx) error {

	return nil
}
