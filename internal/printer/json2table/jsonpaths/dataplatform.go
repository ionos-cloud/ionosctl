package jsonpaths

// Dataplatform json paths
var (
	DataplatformCluster = map[string]string{
		"Id":           "id",
		"Name":         "properties.name",
		"Version":      "properties.dataPlatformVersion",
		"DatacenterId": "properties.datacenterId",
		"State":        "metadata.state",
	}

	DataplatformNodepool = map[string]string{
		"Id":               "id",
		"Name":             "properties.name",
		"Nodes":            "properties.nodeCount",
		"Cores":            "properties.coresCount",
		"CpuFamily":        "properties.cpuFamily",
		"State":            "metadata.state",
		"AvailabilityZone": "properties.availabilityZone",
		"Labels":           "properties.labels",
		"Annotations":      "properties.annotations",
		"ClusterId":        "href",
	}
)
