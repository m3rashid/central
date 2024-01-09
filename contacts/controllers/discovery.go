package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/m3rashid/central/internal/discovery"
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
}

func Discovery(ctx *fiber.Ctx) error {
	return ctx.Status(fiber.StatusOK).JSON(endpoints)
}
