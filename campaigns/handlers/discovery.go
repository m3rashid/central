package handlers

import (
	"internal/discovery"

	"github.com/gofiber/fiber/v2"
)

var endpoints = discovery.ResourceServerDetails{
	BaseUrl: "http://localhost:5001",
	Endpoints: discovery.EndpointsMap{
		"/all": {
			Method:      "GET",
			Description: "Get all contacts",
		},
		"/:id": {
			Method:      "GET",
			Description: "Get contact by ID",
		},
	},
	AllowedScopes: discovery.AllowedScopes{
		"contacts": {"read", "create", "update", "delete"},
	},
}

func Discovery(ctx *fiber.Ctx) error {
	return ctx.Status(fiber.StatusOK).JSON(endpoints)
}
