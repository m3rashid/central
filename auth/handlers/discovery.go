package handlers

import (
	"encoding/json"
	"fmt"
	"internal/discovery"
	"internal/helpers"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
)

var localResourceDiscoveryUrls = map[string]string{
	"contacts": "http://localhost:5001/discovery",
	// ...
}

var dockerResourceDiscoveryUrls = map[string]string{
	"contacts": "http://contacts-webserver:5001/discovery",
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
	serverEnvironment := os.Getenv("ENVIRONMENT")
	resourceDiscoveryUrls := helpers.Ternary[map[string]string](
		serverEnvironment == "docker",
		dockerResourceDiscoveryUrls,
		localResourceDiscoveryUrls,
	)

	for range timeTick.C {
		for name, url := range resourceDiscoveryUrls {
			log.Printf("Discovering %s server at %s", name, url)
			response, err := http.Get(url)
			if err != nil {
				fmt.Printf("Error in getting response from %s server\n", name)
				break
			}

			byteArr, err := io.ReadAll(response.Body)
			if err != nil {
				break
			}

			var resourceServerDetails discovery.ResourceServerDetails
			if err := json.Unmarshal(byteArr, &resourceServerDetails); err != nil {
				fmt.Printf("Error in parsing response from %s server\n", name)
				break
			}

			resourceServerDetails.LastUpdated = time.Now()
			resourceServerDetailsMap[name] = resourceServerDetails
		}
	}
}
