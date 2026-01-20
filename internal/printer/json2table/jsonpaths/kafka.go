package jsonpaths

var (
	KafkaCluster = map[string]string{
		"Id":              "id",
		"Name":            "properties.name",
		"Version":         "properties.version",
		"Size":            "properties.size",
		"DatacenterId":    "properties.connections.0.datacenterId",
		"LanId":           "properties.connections.0.lanId",
		"BrokerAddresses": "properties.connections.0.brokerAddresses",
		"State":           "metadata.state",
		"StateMessage":    "metadata.message",
	}

	KafkaTopic = map[string]string{
		"Id":                 "id",
		"Name":               "properties.name",
		"ReplicationFactor":  "properties.replicationFactor",
		"NumberOfPartitions": "properties.numberOfPartitions",
		"RetentionTime":      "properties.logRetention.retentionTime",
		"SegmentByes":        "properties.logRetention.segmentBytes",
		"State":              "metadata.state",
		"ClusterId":          "href",
	}
)
