package jsonpaths

var VPNWireguardGateway = map[string]string{
	"ID":             "id",
	"Name":           "properties.name",
	"Description":    "properties.description",
	"GatewayIP":      "properties.gatewayIP",
	"InterfaceIPv4":  "properties.interfaceIPv4CIDR",
	"InterfaceIPv6":  "properties.interfaceIPv6CIDR",
	"ListenPort":     "properties.listenPort",
	"DatacenterId":   "properties.connections[0].datacenterId",
	"LanId":          "properties.connections[0].lanId",
	"ConnectionIPv4": "properties.connections[0].ipv4CIDR",
	"ConnectionIPv6": "properties.connections[0].ipv6CIDR",
	"Status":         "metadata.status",
	"StatusMessage":  "metadata.statusMessage",
	"PublicKey":      "metadata.publicKey",
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
