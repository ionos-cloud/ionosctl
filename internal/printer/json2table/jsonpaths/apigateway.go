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
)
