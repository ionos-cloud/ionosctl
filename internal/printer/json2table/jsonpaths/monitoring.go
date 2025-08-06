package jsonpaths

var (
	MonitoringPipeline = map[string]string{
		"Id":              "id",
		"Name":            "properties.name",
		"Status":          "metadata.status",
		"GrafanaEndpoint": "metadata.grafanaEndpoint",
		"HttpEndpoint":    "metadata.httpEndpoint",
	}

	MonitoringCentral = map[string]string{
		"Id":              "id",
		"Enabled":         "properties.enabled",
		"GrafanaEndpoint": "metadata.grafanaEndpoint",
		"Products":        "metadata.products.*",
	}
)
