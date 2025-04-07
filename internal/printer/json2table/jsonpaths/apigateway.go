package jsonpaths

var (
	ApiGatewayGateway = map[string]string{
		"Id":             "id",
		"Name":           "properties.name",
		"Logs":           "properties.logs",
		"Metrics":        "properties.metrics",
		"Status":         "metadata.status",
		"PublicEndpoint": "metadata.publicEndpoint",
	}

	ApiGatewayRoute = map[string]string{
		"Id":            "id",
		"Name":          "properties.name",
		"Type":          "properties.type",
		"Paths":         "properties.paths",
		"Methods":       "properties.methods",
		"WebSocket":     "properties.webSocket",
		"Scheme":        "properties.upstreams.0.scheme",
		"LoadBalancer":  "properties.upstreams.0.loadbalancers",
		"Host":          "properties.upstreams.0.host",
		"Port":          "properties.upstreams.0.port",
		"Weight":        "properties.upstreams.0.weight",
		"Status":        "metadata.status",
		"StatusMessage": "metadata.statusMessage",
	}
)
