package jsonpaths

var VPNWireguardGateway = map[string]string{
	"ID":             "id",
	"Name":           "properties.name",
	"Description":    "properties.description",
	"GatewayIP":      "properties.gatewayIP",
	"DatacenterId":   "properties.connections[0].datacenterId",
	"LanId":          "properties.connections[0].lanId",
	"ConnectionIPv4": "properties.connections[0].ipv4CIDR",
	"ConnectionIPv6": "properties.connections[0].ipv6CIDR",
	"Status":         "metadata.status",
	"StatusMessage":  "metadata.statusMessage",
	"Version":        "properties.version",
}

var VPNWireguardPeer = map[string]string{
	"ID":           "id",
	"Name":         "properties.name",
	"Description":  "properties.description",
	"Host":         "properties.endpoint.host",
	"Port":         "properties.endpoint.port",
	"WhitelistIPs": "properties.allowedIPs",
	"PublicKey":    "properties.publicKey",
	"Status":       "metadata.status",
}

var VPNIPSecGateway = map[string]string{
	"ID":             "id",
	"Name":           "properties.name",
	"Description":    "properties.description",
	"GatewayIP":      "properties.gatewayIP",
	"Version":        "properties.version",
	"DatacenterId":   "properties.connections[0].datacenterId",
	"LanId":          "properties.connections[0].lanId",
	"ConnectionIPv4": "properties.connections[0].ipv4CIDR",
	"ConnectionIPv6": "properties.connections[0].ipv6CIDR",
	"Status":         "metadata.status",
	"StatusMessage":  "metadata.statusMessage",
}

var VPNIPSecTunnel = map[string]string{
	"ID":                     "id",
	"Name":                   "properties.name",
	"Description":            "properties.description",
	"RemoteHost":             "properties.remoteHost",
	"AuthMethod":             "properties.auth.method",
	"PSKKey":                 "properties.auth.psk.key",
	"IKEDiffieHellmanGroup":  "properties.ike.diffieHellmanGroup",
	"IKEEncryptionAlgorithm": "properties.ike.encryptionAlgorithm",
	"IKEIntegrityAlgorithm":  "properties.ike.integrityAlgorithm",
	"IKELifetime":            "properties.ike.lifetime",
	"ESPDiffieHellmanGroup":  "properties.esp.diffieHellmanGroup",
	"ESPEncryptionAlgorithm": "properties.esp.encryptionAlgorithm",
	"ESPIntegrityAlgorithm":  "properties.esp.integrityAlgorithm",
	"ESPLifetime":            "properties.esp.lifetime",
	"CloudNetworkCIDRs":      "properties.cloudNetworkCIDRs",
	"PeerNetworkCIDRs":       "properties.peerNetworkCIDRs",
	"Status":                 "metadata.status",
	"StatusMessage":          "metadata.statusMessage",
}
