package discovery

import "time"

type EndpointsMap map[string]struct {
	Method      string `json:"method"`
	Description string `json:"description"`
}

type AllowedScopes map[string][]string

type ResourceServerDetails struct {
	BaseUrl       string        `json:"url"`
	Endpoints     EndpointsMap  `json:"endpoints"`
	LastUpdated   time.Time     `json:"last_updated"`
	AllowedScopes AllowedScopes `json:"scopes"`
}

type ResourceServersMap = map[string]ResourceServerDetails

// there will be no scopes in the request
// scopes are registered by the client app at the key creation time
// those all scopes are by default taken at the request handler
type IncomingScopeRequest = string

func (scopes *AllowedScopes) ToString() (string, error) {

	return "", nil
}
