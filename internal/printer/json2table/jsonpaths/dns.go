package jsonpaths

// DNS json paths
var (
	DnsReverseRecord = map[string]string{
		"Id":          "id",
		"Name":        "properties.name",
		"Description": "properties.description",
		"IP":          "properties.ip",
	}

	DnsRecord = map[string]string{
		"Id":      "id",
		"Name":    "properties.name",
		"Content": "properties.content",
		"Type":    "properties.type",
		"Enabled": "properties.enabled",
		"FQDN":    "metadata.fqdn",
		"ZoneId":  "metadata.zoneId",

		// State is only in Zone Records
		"State": "metadata.state",

		// RootName is only in Secondary Zone Records
		"RootName": "metadata.rootName",
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
		"Id":              "id",
		"Algorithm":       "properties.keyParameters.algorithm",
		"KskBits":         "properties.keyParameters.kskBits",
		"ZskBits":         "properties.keyParameters.zskBits",
		"NsecMode":        "properties.nsecParameters.nsecMode",
		"Nsec3Iterations": "properties.nsecParameters.nsec3Iterations",
		"Nsec3SaltBits":   "properties.nsecParameters.nsec3SaltBits",
		"Validity":        "properties.validity",
	}

	DnsSecondaryZone = map[string]string{
		"Id":          "id",
		"Name":        "properties.zoneName",
		"Description": "properties.description",
		"PrimaryIPs":  "properties.primaryIps",
		"State":       "metadata.state",
	}

	DnsSecondaryZoneTransfer = map[string]string{
		"PrimaryIP":    "primaryIp",
		"Status":       "status",
		"ErrorMessage": "errorMessage",
	}
)
