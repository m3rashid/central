package handlers

import (
	"github.com/gofiber/fiber/v2"
	"internal/helpers"
	"strconv"
	"strings"
)

const StateKey = "state"
const ClientIDKey = "client_id"
const ResponseTypeKey = "ResponseType"

const LocalUsersCookieName = "local-users"

type FlowQueries struct {
	ClientID     string
	ResponseType string
	State        string
}

func getFlowQueries(ctx *fiber.Ctx) FlowQueries {
	return FlowQueries{
		ClientID:     ctx.Query(ClientIDKey, ""),
		ResponseType: ctx.Query(ResponseTypeKey, ""),
		State:        ctx.Query(StateKey, ""),
	}
}

func SetFlowQueries(baseUrl string, flowQueries FlowQueries) string {
	if flowQueries.ClientID == "" || flowQueries.ResponseType == "" {
		return baseUrl
	}

	return baseUrl + "?" +
		ClientIDKey + "=" + flowQueries.ClientID + "&" +
		ResponseTypeKey + "=" + flowQueries.ResponseType +
		helpers.Ternary[string](flowQueries.State != "", "&"+StateKey+"="+flowQueries.State, "")
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
		Name:     LocalUsersCookieName,
		Value:    userIDsString,
		Domain:   "localhost", // TODO: handle this for deployments
	})
}

func getLocalUserIDsFromCookie(ctx *fiber.Ctx) []uint {
	var users []uint
	userIDsString := ctx.Cookies(LocalUsersCookieName, "")
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
