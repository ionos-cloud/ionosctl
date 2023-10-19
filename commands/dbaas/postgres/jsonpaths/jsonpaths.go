package jsonpaths

var ApiVersion = map[string]string{
	"Version": "name",
}

var Backup = map[string]string{
	"BackupId":                   "id",
	"ClusterId":                  "properties.clusterId",
	"EarliestRecoveryTargetTime": "properties.earliestRecoveryTargetTime",
	"Version":                    "properties.version",
	"Active":                     "properties.active",
	"CreatedDate":                "metadata.createdDate",
	"State":                      "metadata.state",
}

var Cluster = map[string]string{
	"ClusterId":           "id",
	"Location":            "properties.location",
	"BackupLocation":      "properties.backupLocation",
	"State":               "metadata.state",
	"DisplayName":         "properties.displayName",
	"PostgresVersion":     "properties.postgresName",
	"Instances":           "properties.instances",
	"StorageType":         "properties.storageType",
	"SynchronizationMode": "properties.synchronizationMode",
}

var LogsMessage = map[string]string{
	"Message": "message",
	"Time":    "time",
}
