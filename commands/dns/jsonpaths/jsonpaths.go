package jsonpaths

var Record = map[string]string{
	"Id":      "id",
	"Name":    "properties.name",
	"Content": "properties.content",
	"Type":    "properties.type",
	"Enabled": "properties.name",
	"FQDN":    "metadata.fqdn",
	"State":   "metadata.state",
	"ZoneId":  "zoneId",
}

var Zone = map[string]string{
	"Id":          "id",
	"Name":        "properties.zoneName",
	"Description": "properties.description",
	"NameServers": "metadata.nameServers",
	"Enabled":     "properties.enabled",
	"State":       "metadata.state",
}
