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

func RenderConsentScreen(ctx *fiber.Ctx) error {
	ctx.Set("Content-Type", "text/html")
	client, flowQueries, err := getClient(ctx)
	if err != nil {
		return errorComponent(ctx, models.Client{}, errors.New("client not found"))
	}

	consentScreenScopes := components.ConsentScreenScope{}
	for modelName, scopes := range client.Scopes {
		for _, item := range scopes.([]interface{}) {
			consentScreenScopes[modelName] = append(consentScreenScopes[modelName], item.(string))
		}
	}

	component := components.ConsentScreen(components.ConsentScreenProps{
		Client:          client,
		Scopes:          consentScreenScopes,
		AllowConsentUrl: setUrlWithFlowQueries("/handle-consent", flowQueries) + "&" + consentQueryKey + "=true",
		DenyConsentUrl:  setUrlWithFlowQueries("/handle-consent", flowQueries) + "&" + consentQueryKey + "=false",
	})
	return component.Render(ctx.Context(), ctx.Response().BodyWriter())
}

func HandleConsent(ctx *fiber.Ctx) error {
	ctx.Set("Content-Type", "text/html")
	client, flowQueries, err := getClient(ctx)
	if err != nil || flowQueries.SelectedUserID == "0" {
		return errorComponent(ctx, models.Client{}, errors.New("client or user not found"))
	}

	consent := ctx.Query(consentQueryKey, "false")
	if consent == "true" {
		db, err := utils.GetDb()
		if err != nil {
			return errorComponent(ctx, models.Client{}, errors.New("unexpected error occured"))
		}

		u64, err := strconv.ParseUint(flowQueries.SelectedUserID, 10, 32)
		if err != nil {
			return errorComponent(ctx, models.Client{}, errors.New("unexpected error occured"))
		}

		db.Table(models.USER_TABLE_NAME).Where("id = ?", uint(u64)).Association("ConnectedApps").Append(&client)

		return ctx.Redirect(
			client.SuccessRedirectUri +
				helpers.Ternary[string](flowQueries.State != "", "?"+stateQueryKey+"="+flowQueries.State, ""),
		)
	}
	return ctx.Redirect(client.FailureRedirectUri)
}
