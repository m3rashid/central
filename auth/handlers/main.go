package handlers

import (
	"errors"

	"github.com/gofiber/fiber/v2"
)

func InitFlow(ctx *fiber.Ctx) error {
	ctx.Set("Content-Type", "text/html")
	client, flowQueries, err := getClient(ctx)
	if err != nil {
		return errorComponent(ctx, client, errors.New("unexpected error occured"))
	}

	// check if cookie exists in the browser and if userIds are in it
	userIds := getLocalUserIDsFromCookie(ctx)
	if len(userIds) > 0 {
		// if users present, show the select account screen
		return ctx.Redirect(setUrlWithFlowQueries("/select-user", flowQueries))
	} else {
		// else, show the login screen
		return ctx.Redirect(setUrlWithFlowQueries("/login", flowQueries))
	}
}
