package core

import "github.com/Bruary/staff-scheduling/models"

type Endpoint struct {
	Path   string
	Method string
}

type EndpointConfig struct {
	RequireJWT  bool
	AccessLevel models.PermissionLevel
}

var Endpoints = map[Endpoint]EndpointConfig{
	{Path: "/api/v1/signup", Method: "POST"}: {RequireJWT: false, AccessLevel: models.Basic},
	{Path: "/api/v1/login", Method: "POST"}:  {RequireJWT: false, AccessLevel: models.Basic},
}
