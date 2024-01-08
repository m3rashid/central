package controllers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
)

type Endpoints map[string]struct {
	Method      string `json:"method"`
	Description string `json:"description"`
}

type ResourceServer struct {
	Name        string    `json:"name"`
	BaseUrl     string    `json:"url"`
	Endpoints   Endpoints `json:"endpoints"`
	LastUpdated time.Time `json:"last_updated"`
}

var ResourceDiscoveryUrls = map[string]string{
	"contacts": "http://localhost:5001/discovery",
	// ...
}

var resourceServers = []ResourceServer{}

/**
 * the idea is to pign the resource servers in the background
 * in a specified time interval and cache the results
 */
func PingResourceServers() {
	for tick := range time.Tick(3 * time.Second) {
		fmt.Println(tick)

		for name, url := range ResourceDiscoveryUrls {
			response, err := http.Get(url)
			if err != nil {
				return
			}

			fmt.Println(name, response)
		}
	}
}

func GetResourceUrls(ctx *fiber.Ctx) error {
	type ResourceDiscoveryInput struct {
		ResourceServerName string `json:"app_name" validate:"required"`
	}
	var resourceDiscoveryInput ResourceDiscoveryInput

	if err := ctx.BodyParser(&resourceDiscoveryInput); err != nil {
		return err
	}

	return nil
}
