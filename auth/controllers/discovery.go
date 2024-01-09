package controllers

import (
	"encoding/json"
	"fmt"
	"internal/discovery"
	"io"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
)

var resourceDiscoveryUrls = map[string]string{
	"contacts": "http://localhost:5001/discovery",
	// ...
}

var resourceServerDetailsMap = discovery.ResourceServersMap{}

func DiscoverResourceServers(ctx *fiber.Ctx) error {
	scope := ctx.Query("scope", "all")
	if scope == "all" {
		return ctx.Status(fiber.StatusOK).JSON(resourceServerDetailsMap)
	}

	details, resourceServerExists := resourceServerDetailsMap[scope]
	if !resourceServerExists {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "No such resource server exists",
		})
	}
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		scope: details,
	})
}

/**
 * ping the resource servers in the background
 * in specified time intervals and cache the results
 */
func PingResourceServers() {
	timeTick := time.NewTicker(time.Second * 5)

	for tick := range timeTick.C {
		fmt.Println(tick)

		for name, url := range resourceDiscoveryUrls {
			response, err := http.Get(url)
			if err != nil {
				return
			}

			byteArr, err := io.ReadAll(response.Body)
			if err != nil {
				return
			}

			var resourceServerDetails discovery.ResourceServerDetails
			if err := json.Unmarshal(byteArr, &resourceServerDetails); err != nil {
				return
			}

			currentResourceServerDetails, resourceServerExists := resourceServerDetailsMap[name]
			if resourceServerExists {
				currentResourceServerDetails = resourceServerDetails
			} else {
				resourceServerDetailsMap[name] = resourceServerDetails
			}
			currentResourceServerDetails.LastUpdated = time.Now()
			resourceDiscoveryUrls[name] = resourceServerDetails.BaseUrl + "/discovery"
		}
	}
}
