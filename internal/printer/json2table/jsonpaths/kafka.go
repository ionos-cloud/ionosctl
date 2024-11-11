package jsonpaths

var KafkaCluster = map[string]string{
	"Id":              "id",
	"Name":            "properties.name",
	"Version":         "properties.version",
	"Size":            "properties.size",
	"DatacenterId":    "properties.connections.0.datacenterId",
	"LanId":           "properties.connections.0.lanId",
	"BrokerAddresses": "properties.connections.0.brokers",
	"State":           "metadata.state",
}
