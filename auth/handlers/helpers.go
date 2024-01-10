package handlers

import (
	"errors"
	"internal/helpers"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/m3rashid/central/auth/components"
	"github.com/m3rashid/central/auth/models"
	"github.com/m3rashid/central/auth/utils"
)

const stateQueryKey = "state"
const clientIDQueryKey = "client_id"
const responseTypeQueryKey = "response_type"
const scopesQueryKey = "scopes"

const localUsersCookieName = "local_users"
const selectedUserIDKey = "selected_user"

func errorComponent(ctx *fiber.Ctx, client models.Client, err error) error {
	component := components.FlowError([]string{err.Error()}, client)
	return component.Render(ctx.Context(), ctx.Response().BodyWriter())
}

type FlowQueries struct {
	ClientID     string
	ResponseType string
	State        string
	Scopes       []string
}

func getFlowQueries(ctx *fiber.Ctx) (FlowQueries, error) {
	var flowQueries FlowQueries
	scopes := ctx.Query(scopesQueryKey, "")
	flowQueries.ClientID = ctx.Query(clientIDQueryKey, "")
	flowQueries.ResponseType = ctx.Query(responseTypeQueryKey, "")

	if flowQueries.ClientID == "" || flowQueries.ResponseType == "" {
		return flowQueries, errors.New("client_id and/or response_type missing")
	}

	flowQueries.Scopes = append(flowQueries.Scopes, strings.Split(scopes, ",")...)
	return flowQueries, nil
}

func setUrlWithFlowQueries(baseUrl string, flowQueries FlowQueries) string {
	if flowQueries.ClientID == "" || flowQueries.ResponseType == "" {
		return baseUrl
	}

	return baseUrl + "?" +
		clientIDQueryKey + "=" + flowQueries.ClientID + "&" +
		responseTypeQueryKey + "=" + flowQueries.ResponseType +
		helpers.Ternary[string](flowQueries.State != "", "&"+stateQueryKey+"="+flowQueries.State, "")
}

func setLocalUsersCookie(ctx *fiber.Ctx, userIDs []uint) {
	var userIDStringArray []string
	for _, userId := range userIDs {
		userIDStringArray = append(userIDStringArray, strconv.FormatUint(uint64(userId), 10))
	}

	userIDsString := strings.Join(userIDStringArray, ",")
	// TODO: also hash this

	ctx.Cookie(&fiber.Cookie{
		HTTPOnly: true,
		Name:     localUsersCookieName,
		Value:    userIDsString,
		Domain:   "localhost", // TODO: handle this for deployments
	})
}

func getLocalUserIDsFromCookie(ctx *fiber.Ctx) []uint {
	var users []uint
	userIDsString := ctx.Cookies(localUsersCookieName, "")
	userIDsStringArray := strings.Split(userIDsString, ",")
	for _, userIDString := range userIDsStringArray {
		u64, err := strconv.ParseUint(userIDString, 10, 32)
		if err != nil {
			return users
		}
		users = append(users, uint(u64))
	}

	return users
}

func addUserIDToCookie(ctx *fiber.Ctx, userID uint) {
	existingUserIDs := getLocalUserIDsFromCookie(ctx)
	existingUserIDs = append(existingUserIDs, userID)
	setLocalUsersCookie(ctx, existingUserIDs)
}

func removeUserIDFromCookie(ctx *fiber.Ctx, userID uint) {
	var newUserIDs []uint
	existingUserIDs := getLocalUserIDsFromCookie(ctx)
	for _, existingUserID := range existingUserIDs {
		if existingUserID != userID {
			newUserIDs = append(newUserIDs, existingUserID)
		}
	}

	setLocalUsersCookie(ctx, newUserIDs)
}

func getClient(ctx *fiber.Ctx) (models.Client, FlowQueries, error) {
	var client models.Client
	flowQueries, err := getFlowQueries(ctx)
	if err != nil {
		return client, flowQueries, nil
	}

	db, err := utils.GetDb()
	if err != nil {
		return client, flowQueries, err
	}

	err = db.Preload("Scopes").Preload("Scopes.Permission").Where("client_id = ?", flowQueries.ClientID).First(&client).Error
	if err != nil || client.ID == 0 {
		return client, flowQueries, err
	}

	// check if the scopes asked for, is included in the registered scopes of the client
	for _, registeredScope := range client.Scopes {
		var scopeMatched = false
		// TODO: make it granular -- also check the permission levels
		for _, scope := range flowQueries.Scopes {
			if scope == registeredScope.Name {
				scopeMatched = true
				break
			}
		}

		if !scopeMatched {
			return client, flowQueries, errors.New("invalid scope")
		}
	}

	return client, flowQueries, nil
}
