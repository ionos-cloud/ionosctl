package authv1

const (
	listTokenExample     = `ionosctl token list`
	getTokenExample      = `ionosctl token get --token-id TOKEN_ID`
	generateTokenExample = `ionosctl token generate`
	deleteTokenExample   = `ionosctl token delete --token-id TOKEN_ID

ionosctl token delete --expired`
)
