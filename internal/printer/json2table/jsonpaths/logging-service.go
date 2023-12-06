package jsonpaths

var (
	LoggingServicePipeline = map[string]string{
		"Id":             "id",
		"Name":           "properties.name",
		"TCPAddress":     "properties.tcpAddress",
		"HTTPAddress":    "properties.httpAddress",
		"GrafanaAddress": "properties.grafanaAddress",
		"CreatedDate":    "metadata.createdDate",
		"State":          "metadata.state",
	}

	LoggingServiceLogs = map[string]string{
		"Protocol": "protocol",
		"Source":   "source",
		"Public":   "public",
		"Tag":      "tag",
		"Labels":   "labels",
	}
)
