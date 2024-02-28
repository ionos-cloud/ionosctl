package jsonpaths

// DBaaS json paths
var (
	/*
			{
		      "id": "498ae72f-411f-11eb-9d07-046c59cc737e",
		      "properties": {
		        "clusterId": "498ae72f-411f-11eb-9d07-046c59cc737e",
		        "items": [
		          {
		            "created": "2020-12-10T13:37:50+01:00",
		            "size": 543
		          }
		        ]
		      }
			}
	*/
	DbaasMariadbBackup = map[string]string{
		"BackupId":  "id",
		"ClusterId": "properties.clusterId",
	}

	/*
		"properties":{"connections":[{"cidr":"192.168.1.209/24","datacenterId":"5ed7be83-9fc5-44ea-9c4d-f15f1fc83345","lanId":"1"}],"cores":1,"displayName":"asdfasdf","dnsName":"ma-redlvpvg74vj6qjc.mariadb.de-txl.ionos.com","instances":1,"maintenanceWindow":{"dayOfTheWeek":"Monday","time":"16:00:00"},"mariadbVersion":"10.11","ram":2048,"storageSize":10}}
	*/
	DbaasMariadbCluster = map[string]string{
		"ClusterId": "id",
		"Name":      "properties.displayName",
		"DNS":       "properties.dnsName",
		"Instances": "properties.instances",
		"Version":   "properties.mariadbVersion",
		"State":     "metadata.state",

		"Cores":       "properties.cores",
		"Ram":         "properties.ram",
		"StorageSize": "properties.storageSize",
		"StorageType": "properties.storageType",
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
		"Location":            "properties.location",
		"BackupLocation":      "properties.backupLocation",
		"State":               "metadata.state",
		"DisplayName":         "properties.displayName",
		"PostgresVersion":     "properties.postgresName",
		"Instances":           "properties.instances",
		"StorageType":         "properties.storageType",
		"SynchronizationMode": "properties.synchronizationMode",
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
)
