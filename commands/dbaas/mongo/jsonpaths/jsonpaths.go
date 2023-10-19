package jsonpaths

var Cluster = map[string]string{
	"ClusterId":    "id",
	"Name":         "properties.displayName",
	"Edition":      "properties.edition",
	"Type":         "properties.type",
	"URL":          "properties.connectionString",
	"Instances":    "properties.instances",
	"Shards":       "properties.shards",
	"Health":       "metadata.health",
	"State":        "metadata.state",
	"MongoVersion": "properties.mongoDBVersion",
	"Location":     "properties.location",
	"TemplateId":   "properties.templateID",
	"Cores":        "properties.cores",
	"StorageType":  "properties.storageType",
}

var ApiVersion = map[string]string{
	"Version": "name",
	"Href":    "swaggerUrl",
}

var LogsMessage = map[string]string{
	"Message": "message",
	"Time":    "time",
}

var Snapshot = map[string]string{
	"SnapshotId":   "id",
	"CreationTime": "properties.creationTime",
	"Size":         "properties.size",
	"Version":      "properties.version",
}

var Templates = map[string]string{
	"TemplateId": "id",
	"Name":       "properties.name",
	"Edition":    "properties.edition",
	"Cores":      "properties.cores",
}

var User = map[string]string{
	"Username":  "properties.username",
	"CreatedBy": "metadata.createdBy",
}
