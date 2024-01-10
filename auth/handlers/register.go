package handlers

import (
	"errors"
	"internal/helpers"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/m3rashid/central/auth/components"
	"github.com/m3rashid/central/auth/models"
	"github.com/m3rashid/central/auth/utils"
)

func RenderRegisterScreen(ctx *fiber.Ctx) error {
	ctx.Set("Content-Type", "text/html")
	client, flowQueries, err := getClient(ctx)
	if err != nil {
		return errorComponent(ctx, models.Client{}, errors.New("client not found"))
	}

	return components.LoginOrRegister(components.LoginProps{
		IsRegister:         true,
		LoginEndpoint:      setUrlWithFlowQueries("/login", flowQueries),
		RegisterEndpoint:   setUrlWithFlowQueries("/register", flowQueries),
		SelectUserEndpoint: setUrlWithFlowQueries("/select-user", flowQueries),
	}, client).Render(ctx.Context(), ctx.Response().BodyWriter())
}

func HandleRegister(ctx *fiber.Ctx) error {
	ctx.Set("Content-Type", "text/html")
	client, flowQueries, err := getClient(ctx)
	if err != nil {
		return errorComponent(ctx, models.Client{}, errors.New("client not found"))
	}

	type Register struct {
		Name     string `json:"name" validate:"required"`
		Email    string `json:"email" validate:"required"`
		Password string `json:"password" validate:"required"`
	}

	var register Register
	if err := ctx.BodyParser(&register); err != nil {
		return errorComponent(ctx, client, errors.New("bad request"))
	}

	if err = validator.New().Struct(register); err != nil {
		return ctx.Redirect(client.FailureRedirectUri)
	}

	db, err := utils.GetDb()
	if err != nil {
		return ctx.Redirect(client.FailureRedirectUri)
	}

	hashedPassword, err := helpers.HashPassword(register.Password)
	if err != nil {
		return ctx.Redirect(client.FailureRedirectUri)
	}

	// create a new OTP and send it to the user's email
	otp := helpers.GenerateOTP()
	user := models.User{
		OTP:      otp,
		Name:     register.Name,
		Email:    register.Email,
		Password: hashedPassword,
	}
	db.Create(&user)

	addUserIDToCookie(ctx, user.ID)
	return ctx.Redirect(setUrlWithFlowQueries("/select-user", flowQueries))
}
