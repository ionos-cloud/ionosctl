package jsonpaths

// DNS json paths
var (
	DnsRecord = map[string]string{
		"Id":      "id",
		"Name":    "properties.name",
		"Content": "properties.content",
		"Type":    "properties.type",
		"Enabled": "properties.name",
		"FQDN":    "metadata.fqdn",
		"State":   "metadata.state",
		"ZoneId":  "zoneId",
	}

	DnsZone = map[string]string{
		"Id":          "id",
		"Name":        "properties.zoneName",
		"Description": "properties.description",
		"NameServers": "metadata.nameServers",
		"Enabled":     "properties.enabled",
		"State":       "metadata.state",
	}

	DnsQuota = map[string]string{
		"ZonesUsed":           "quotaUsage.zones",
		"ZonesLimit":          "quotaLimits.zones",
		"SecondaryZonesUsed":  "quotaUsage.secondaryZones",
		"SecondaryZonesLimit": "quotaLimits.secondaryZones",
		"RecordsUsed":         "quotaUsage.records",
		"RecordsLimit":        "quotaLimits.records",
		"ReverseRecordsUsed":  "quotaUsage.reverseRecords",
		"ReverseRecordsLimit": "quotaLimits.reverseRecords",
	}

	DnsSecKey = map[string]string{
		"ID":        "id",
		"Algorithm": "properties.keyParameters.algorithm",
		"NsecMode":  "properties.nsecParameters.nsecMode",
	}
)
