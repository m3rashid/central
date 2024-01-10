package discovery

import "time"

type EndpointsMap map[string]struct {
	Method      string `json:"method"`
	Description string `json:"description"`
}

type Scopes map[string][]string

type ResourceServerDetails struct {
	BaseUrl     string       `json:"url"`
	Endpoints   EndpointsMap `json:"endpoints"`
	LastUpdated time.Time    `json:"last_updated"`
	Scopes      Scopes       `json:"scopes"`
}

type ResourceServersMap = map[string]ResourceServerDetails
