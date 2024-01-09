package handlers

import (
	"errors"
	"github.com/gofiber/fiber/v2"
)

func InitFlow(ctx *fiber.Ctx) error {
	flowQueries := getFlowQueries(ctx)

	if flowQueries.ClientID == "" || flowQueries.ResponseType == "" {
		return errors.New("client_id and/or response_type missing")
	}

	userCookie := ctx.Cookies("local-users", "")
	if userCookie == "" {
		return ctx.Redirect(SetFlowQueries("/login", flowQueries))
	}

	return nil
}
