package jsonpaths

var CDNDistribution = map[string]string{
	"Id":            "id",
	"Domain":        "properties.domain",
	"CertificateId": "properties.certificateId",
	"State":         "metadata.state",
}

var CDNRoutingRule = map[string]string{
	"Scheme":                   "scheme",
	"Prefix":                   "prefix",
	"Host":                     "upstream.host",
	"Caching":                  "upstream.caching",
	"Waf":                      "upstream.waf",
	"SniMode":                  "upstream.sniMode",
	"GeoRestrictionsAllowList": "upstream.geoRestrictions.allowList",
	"GeoRestrictionsBlockList": "upstream.geoRestrictions.blockList",
	"RateLimitClass":           "upstream.rateLimitClass",
}
