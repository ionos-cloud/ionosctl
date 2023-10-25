package jsonpaths

// Auth json paths
var (
	AuthToken = map[string]string{
		"TokenId":        "id",
		"CreatedDate":    "createdDate",
		"ExpirationDate": "expirationDate",
		"Href":           "href",
	}
	AuthTokenInfo = map[string]string{
		"TokenId":        "tokenId",
		"UserId":         "userId",
		"ContractNumber": "contractNumber",
		"Role":           "role",
	}
)
