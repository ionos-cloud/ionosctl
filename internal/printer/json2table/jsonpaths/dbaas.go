package jsonpaths

// DBaaS json paths
var (
	DbaasMariadbBackup = map[string]string{
		"BackupId":  "id",
		"ClusterId": "properties.clusterId",
	}

	DbaasMariadbCluster = map[string]string{
		"ClusterId": "id",
		"Name":      "properties.displayName",
		"DNS":       "properties.dnsName",
		"Instances": "properties.instances",
		"Version":   "properties.mariadbVersion",
		"State":     "metadata.state",

		"Cores":       "properties.cores",
		"RAM":         "properties.ram",
		"StorageSize": "properties.storageSize",

		"MaintenanceDay":  "properties.maintenanceWindow.dayOfTheWeek",
		"MaintenanceTime": "properties.maintenanceWindow.time",
	}

	DbaasMongoCluster = map[string]string{
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

	DbaasMongoAPIVersion = map[string]string{
		"Version": "name",
		"Href":    "swaggerUrl",
	}

	DbaasLogsMessage = map[string]string{
		"Message": "message",
		"Time":    "time",
	}

	DbaasMongoSnapshot = map[string]string{
		"SnapshotId":   "id",
		"CreationTime": "properties.creationTime",
		"Size":         "properties.size",
		"Version":      "properties.version",
	}

	DbaasMongoTemplates = map[string]string{
		"TemplateId": "id",
		"Name":       "properties.name",
		"Edition":    "properties.edition",
		"Cores":      "properties.cores",
	}

	DbaasMongoUser = map[string]string{
		"Username":  "properties.username",
		"CreatedBy": "metadata.createdBy",
	}

	DbaasPostgresApiVersion = map[string]string{
		"Version": "name",
	}

	DbaasPostgresBackup = map[string]string{
		"BackupId":                   "id",
		"ClusterId":                  "properties.clusterId",
		"EarliestRecoveryTargetTime": "properties.earliestRecoveryTargetTime",
		"Version":                    "properties.version",
		"Active":                     "properties.active",
		"CreatedDate":                "metadata.createdDate",
		"State":                      "metadata.state",
	}

	DbaasPostgresCluster = map[string]string{
		"ClusterId":           "id",
		"State":               "metadata.state",
		"Location":            "properties.location",
		"BackupLocation":      "properties.backupLocation",
		"Cores":               "properties.cores",
		"DisplayName":         "properties.displayName",
		"PostgresVersion":     "properties.postgresVersion",
		"Instances":           "properties.instances",
		"StorageType":         "properties.storageType",
		"SynchronizationMode": "properties.synchronizationMode",
	}

	DbaasPostgresV2Cluster = map[string]string{
		"ClusterId":           "id",
		"State":               "metadata.state",
		"Location":            "properties.location",
		"BackupLocation":      "properties.backupLocation",
		"Cores":               "properties.cores",
		"DisplayName":         "properties.displayName",
		"PostgresVersion":     "properties.postgresVersion",
		"Instances":           "properties.instances",
		"StorageType":         "properties.storageType",
		"SynchronizationMode": "properties.synchronizationMode",
	}

	DbaasPostgresV2Backup = map[string]string{
		"BackupId":                   "id",
		"ClusterId":                  "properties.clusterId",
		"PostgresClusterVersion":     "properties.postgresClusterVersion",
		"IsActive":                   "properties.isActive",
		"EarliestRecoveryTargetTime": "properties.earliestRecoveryTargetTime",
		"LatestRecoveryTargetTime":   "properties.latestRecoveryTargetTime",
		"Location":                   "properties.location",
		"State":                      "metadata.state",
		"CreatedDate":                "metadata.createdDate",
	}

	DbaasPostgresV2BackupLocation = map[string]string{
		"LocationId": "id",
		"Location":   "properties.location",
	}

	DbaasPostgresVersion = map[string]string{
		"Version": "properties.version",
		"Id":      "id",
	}

	DbaasPostgresUser = map[string]string{
		"Id":       "id",
		"Username": "properties.username",
		"System":   "properties.system",
	}

	DbaasPostgresDatabase = map[string]string{
		"Id":    "id",
		"Name":  "properties.name",
		"Owner": "properties.owner",
	}

	DbaasInMemoryDBSnapshot = map[string]string{
		"Id":                 "id",
		"ReplicasetId":       "metadata.replicasetId",
		"DatacenterId":       "metadata.datacenterId",
		"Time":               "metadata.snapshotTime",
		"State":              "metadata.state",
		"RestoredSnapshotId": "metadata.restoredSnapshotId",
	}

	DbaasInMemoryDBReplicaSet = map[string]string{
		"Id":              "id",
		"Name":            "properties.displayName",
		"Version":         "properties.version",
		"DNSName":         "metadata.dnsName",
		"Replicas":        "properties.replicas",
		"Cores":           "properties.resources.cores",
		"RAM":             "properties.resources.ram",
		"StorageSize":     "properties.resources.storage",
		"State":           "metadata.state",
		"BackupLocation":  "properties.backupLocation",
		"PersistenceMode": "properties.persistenceMode",
		"EvictionPolicy":  "properties.evictionPolicy",
		"MaintenanceDay":  "properties.maintenanceWindow.dayOfTheWeek",
		"MaintenanceTime": "properties.maintenanceWindow.time",
		"DatacenterId":    "properties.connections[0].datacenterId",
		"LanId":           "properties.connections[0].lanId",
		"Username":        "properties.credentials.username",
	}
)
