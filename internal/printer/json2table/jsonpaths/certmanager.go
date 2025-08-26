package jsonpaths

// Cert Manager json paths
var (
	CertManagerCertificate = map[string]string{
		"CertId":      "id",
		"DisplayName": "properties.name",
	}

	CertManagerProvider = map[string]string{
		"Id":        "id",
		"Name":      "properties.name",
		"Email":     "properties.email",
		"Server":    "properties.server",
		"KeyId":     "properties.externalAccountBinding.keyId",
		"KeySecret": "properties.externalAccountBinding.keySecret",
		"State":     "metadata.state",
	}

	CertManagerAutocertificate = map[string]string{
		"Id":               "id",
		"Name":             "properties.name",
		"CommonName":       "properties.commonName",
		"KeyAlgorithm":     "properties.keyAlgorithm",
		"Provider":         "properties.provider",
		"AlternativeNames": "properties.subjectAlternativeNames",
		"State":            "metadata.state",
	}
)
