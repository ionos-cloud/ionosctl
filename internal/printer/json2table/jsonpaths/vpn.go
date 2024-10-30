package jsonpaths

var VPNWireguardGateway = map[string]string{
	"ID":             "id",
	"Name":           "properties.name",
	"Description":    "properties.description",
	"IP":             "properties.gateway_ip",
	"InterfaceIPv4":  "properties.interface_ipv4_cidr",
	"InterfaceIPv6":  "properties.interface_ipv6_cidr",
	"DatacenterId":   "properties.connections[0].datacenter_id",
	"LanId":          "properties.connections[0].lan_id",
	"ConnectionIPv4": "properties.connections[0].ipv4_cidr",
	"ConnectionIPv6": "properties.connections[0].ipv6_cidr",
	"InterfaceIP":    "properties.interface_ipv4_cidr",
	"ListenPort":     "properties.listen_port",
	"Status":         "metadata.status",
}
