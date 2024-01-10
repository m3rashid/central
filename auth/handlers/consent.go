package handlers

import (
	"errors"
	"internal/helpers"

	"github.com/gofiber/fiber/v2"
	"github.com/m3rashid/central/auth/components"
	"github.com/m3rashid/central/auth/models"
)

func RenderConsentScreen(ctx *fiber.Ctx) error {
	ctx.Set("Content-Type", "text/html")
	client, flowQueries, err := getClient(ctx)
	if err != nil {
		return errorComponent(ctx, models.Client{}, errors.New("client not found"))
	}

	consentScreenScopes := []components.ConsentScreenScope{}
	for _, scope := range client.Scopes {
		newScope := components.ConsentScreenScope{Name: scope.Name, Permissions: []string{}}
		// TODO: get scopes and permissions from discovery
		for _, permission := range scope.Permissions {
			newScope.Permissions = append(newScope.Permissions, permission.Name)
		}
		consentScreenScopes = append(consentScreenScopes, newScope)
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
	if err != nil {
		return errorComponent(ctx, models.Client{}, errors.New("client not found"))
	}

	consent := ctx.Query(consentQueryKey, "false")
	if consent == "true" {
		// TODO: add app to connected apps
		return ctx.Redirect(client.SuccessRedirectUri + helpers.Ternary[string](flowQueries.State != "", "?"+stateQueryKey+"="+flowQueries.State, ""))
	}
	return ctx.Redirect(client.FailureRedirectUri)
}
