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

func RenderLoginScreen(ctx *fiber.Ctx) error {
	ctx.Set("Content-Type", "text/html")
	client, flowQueries, err := getClient(ctx)
	if err != nil {
		return errorComponent(ctx, models.Client{}, errors.New("client not found"))
	}

	return components.LoginOrRegister(components.LoginProps{
		IsRegister:         false,
		LoginEndpoint:      setUrlWithFlowQueries("/login", flowQueries),
		RegisterEndpoint:   setUrlWithFlowQueries("/register", flowQueries),
		SelectUserEndpoint: setUrlWithFlowQueries("/select-user", flowQueries),
	}, client).Render(ctx.Context(), ctx.Response().BodyWriter())
}

func HandleLogin(ctx *fiber.Ctx) error {
	ctx.Set("Content-Type", "text/html")
	client, flowQueries, err := getClient(ctx)
	if err != nil {
		return errorComponent(ctx, models.Client{}, errors.New("client not found"))
	}

	type Login struct {
		Email    string `json:"email" validate:"email,required"`
		Password string `json:"password" validate:"required"`
	}
	var login Login
	if err := ctx.BodyParser(&login); err != nil {
		return ctx.Redirect(client.FailureRedirectUri)
	}

	if err = validator.New().Struct(login); err != nil {
		return ctx.Redirect(client.FailureRedirectUri)
	}

	db, err := utils.GetDb()
	if err != nil {
		return ctx.Redirect(client.FailureRedirectUri)
	}

	var user models.User
	if err = db.Where("email = ?", login.Email).First(&user).Error; err != nil {
		return ctx.Redirect(client.FailureRedirectUri)
	}

	if user.ID == 0 {
		return ctx.Redirect(setUrlWithFlowQueries("/register", flowQueries))
	}

	if err := helpers.ComparePassword(user.Password, login.Password); err != nil {
		return ctx.Redirect(client.FailureRedirectUri)
	}

	addUserIDToCookie(ctx, user.ID)
	return ctx.Redirect(setUrlWithFlowQueries("/select-user", flowQueries))
}
