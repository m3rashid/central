package discovery

import "time"

type EndpointsMap map[string]struct {
	Method      string `json:"method"`
	Description string `json:"description"`
}

type ResourceServerDetails struct {
	BaseUrl     string       `json:"url"`
	Endpoints   EndpointsMap `json:"endpoints"`
	LastUpdated time.Time    `json:"last_updated"`
}

type ResourceServersMap = map[string]ResourceServerDetails
