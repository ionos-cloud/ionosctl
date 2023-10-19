package jsonpaths

var Cluster = map[string]string{
	"Id":           "id",
	"Name":         "properties.name",
	"Version":      "properties.dataPlatformVersion",
	"DatacenterId": "properties.datacenterId",
	"State":        "metadata.state",
}

var Nodepool = map[string]string{
	"Id":               "id",
	"Name":             "properties.name",
	"Nodes":            "properties.nodeCount",
	"Cores":            "properties.coresCount",
	"CpuFamily":        "properties.cpuFamily",
	"State":            "metadata.state",
	"AvailabilityZone": "properties.availabilityZone",
	"Labels":           "properties.labels",
	"Annotations":      "properties.annotations",
}
