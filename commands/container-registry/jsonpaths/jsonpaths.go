package jsonpaths

var Location = map[string]string{
	"LocationId": "id",
}

var Registry = map[string]string{
	"RegistryId":            "id",
	"DisplayName":           "properties.name",
	"Location":              "properties.location",
	"Hostname":              "properties.hostname",
	"GarbageCollectionDays": "properties.garbageCollectionSchedule.days",
	"GarbageCollectionTime": "properties.garbageCollectionSchedule.time",
}

var Token = map[string]string{
	"TokenId":             "id",
	"DisplayName":         "properties.name",
	"ExpiryDate":          "properties.expiryDate",
	"CredentialsUsername": "properties.credentials.username",
	"CredentialsPassword": "properties.credentials.password",
	"Status":              "properties.status",
}
