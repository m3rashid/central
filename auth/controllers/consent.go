package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/m3rashid/central/auth/components"
)

func RenderConsentScreen(ctx *fiber.Ctx) error {
	ctx.Set("Content-Type", "text/html")
	component := components.Consent()
	return component.Render(ctx.Context(), ctx.Response().BodyWriter())
}
