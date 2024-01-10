package handlers

import (
	"errors"
	"internal/helpers"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/m3rashid/central/auth/components"
	"github.com/m3rashid/central/auth/models"
	"github.com/m3rashid/central/auth/utils"
)

func RenderSelectUserScreen(ctx *fiber.Ctx) error {
	ctx.Set("Content-Type", "text/html")
	client, flowQueries, err := getClient(ctx)
	if err != nil {
		return errorComponent(ctx, models.Client{}, errors.New("client not found"))
	}

	db, err := utils.GetDb()
	if err != nil {
		return ctx.Redirect(client.FailureRedirectUri)
	}

	userIds := getLocalUserIDsFromCookie(ctx)
	var users []models.User
	err = db.Where("id in ?", userIds).Find(&users).Error
	if err != nil {
		return ctx.Redirect(client.FailureRedirectUri)
	}

	return components.SelectUser(components.SelectUserProps{
		Users:             users,
		Client:            client,
		LoginEndpoint:     setUrlWithFlowQueries("/login", flowQueries),
		RegisterEndpoint:  setUrlWithFlowQueries("/register", flowQueries),
		UserIDSelectedURL: setUrlWithFlowQueries("/handle-select-user", flowQueries) + "&" + selectedUserIDKey + "=",
	}).Render(ctx.Context(), ctx.Response().BodyWriter())
}

func HandleSelectUser(ctx *fiber.Ctx) error {
	ctx.Set("Content-Type", "text/html")
	client, flowQueries, err := getClient(ctx)
	if err != nil {
		return errorComponent(ctx, models.Client{}, errors.New("client not found"))
	}

	selectedUserID64, err := strconv.ParseUint(ctx.Query(selectedUserIDKey, "0"), 10, 32)
	if err != nil || selectedUserID64 == 0 {
		return ctx.Redirect(client.FailureRedirectUri)
	}

	db, err := utils.GetDb()
	if err != nil {
		return ctx.Redirect(client.FailureRedirectUri)
	}
	var user models.User
	if err = db.Preload("ConnectedApps").Where("id = ?", uint(selectedUserID64)).First(&user).Error; err != nil {
		return ctx.Redirect(client.FailureRedirectUri)
	}

	appConnected := false
	for _, app := range user.ConnectedApps {
		if client.ClientID == app.ClientID {
			appConnected = true
			break
		}
	}

	return helpers.Ternary[error](
		appConnected,
		ctx.Redirect(setUrlWithFlowQueries(client.SuccessRedirectUri, flowQueries)),
		ctx.Redirect(setUrlWithFlowQueries("/consent", flowQueries)),
	)
}
