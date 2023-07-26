package authv1

const (
	listTokenExample = `ionosctl token list`
	getTokenExample  = `ionosctl token get --token-id TOKEN_ID

ionosctl token get --token TOKEN`
	generateTokenExample = `ionosctl token generate`
	deleteTokenExample   = `ionosctl token delete --token-id TOKEN_ID

ionosctl token delete --token TOKEN

ionosctl token delete --expired

ionosctl token delete --current

ionosctl token delete --all`
	parseTokenExample = `ionosctl token parse --token TOKEN

ionosctl token parse --privileges --token TOKEN`
)
