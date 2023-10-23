package jsonpaths

// Cert Manager json paths
var (
	CertManagerCertificate = map[string]string{
		"CertId":      "id",
		"DisplayName": "properties.name",
	}

	CertManagerAPIVersion = map[string]string{
		"Href":    "href",
		"Name":    "name",
		"Version": "version",
	}
)
