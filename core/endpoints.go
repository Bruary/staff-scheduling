package core

import "github.com/Bruary/staff-scheduling/core/models"

type Endpoint struct {
	Path   string
	Method string
}

type EndpointConfig struct {
	RequireJWT  bool
	AccessLevel models.PermissionLevel
}

var Endpoints = map[Endpoint]EndpointConfig{
	{Path: "/api/v1/signup", Method: "POST"}:         {RequireJWT: false, AccessLevel: models.Basic},
	{Path: "/api/v1/login", Method: "POST"}:          {RequireJWT: false, AccessLevel: models.Basic},
	{Path: "/api/v1/user/permission", Method: "PUT"}: {RequireJWT: true, AccessLevel: models.Admin},
	{Path: "/api/v1/user", Method: "DELETE"}:         {RequireJWT: true, AccessLevel: models.Admin},
}
