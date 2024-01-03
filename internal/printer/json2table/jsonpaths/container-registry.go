package jsonpaths

// Container Registry json paths
var (
	ContainerRegistryLocation = map[string]string{
		"LocationId": "id",
	}

	ContainerRegistryRegistry = map[string]string{
		"RegistryId":            "id",
		"DisplayName":           "properties.name",
		"Location":              "properties.location",
		"Hostname":              "properties.hostname",
		"GarbageCollectionDays": "properties.garbageCollectionSchedule.days",
		"GarbageCollectionTime": "properties.garbageCollectionSchedule.time",
		"VulnerabilityScanning": "properties.features.vulnerabilityScanning.enabled",
	}

	ContainerRegistryToken = map[string]string{
		"TokenId":             "id",
		"DisplayName":         "properties.name",
		"ExpiryDate":          "properties.expiryDate",
		"CredentialsUsername": "properties.credentials.username",
		"CredentialsPassword": "properties.credentials.password",
		"Status":              "properties.status",
	}

	ContainerRegistryArtifact = map[string]string{
		"Id":                     "id",
		"MediaType":              "properties.mediaType",
		"Repository":             "properties.repositoryName",
		"PushCount":              "metadata.pushCount",
		"PullCount":              "metadata.pullCount",
		"LastPushed":             "metadata.lastPushedAt",
		"TotalVulnerabilities":   "metadata.vulnTotalCount",
		"FixableVulnerabilities": "metadata.vulnFixableCount",
		"URN":                    "metadata.resourceURN",
	}
)
