package completer

import (
	"context"
	"fmt"
	"strings"

	"github.com/cilium/fake"
	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func GetCidrCompletionFunc(cmd *core.Command) func(c *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	return func(c *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		databaseIp := "192.168.1.128" // fallback in case of no servers / errs
		ip, err := getNicIp(cmd)
		if err != nil || ip == "" {
			ip = databaseIp
		}
		cidrs := generateCidrs(cmd, ip)
		return []string{strings.Join(cidrs, ",")}, cobra.ShellCompDirectiveNoFileComp
	}
}

func getNicIp(cmd *core.Command) (string, error) {
	ls, _, err := client.Must().CloudClient.ServersApi.DatacentersServersGet(context.Background(),
		viper.GetString(core.GetFlagName(cmd.NS, constants.FlagDatacenterId))).Execute()
	if err != nil || ls.Items == nil || len(ls.Items) == 0 {
		return "", fmt.Errorf("failed getting servers %w", err)
	}

	for _, server := range ls.Items {
		if server.Id == nil {
			return "", fmt.Errorf("failed getting ID")
		}

		nics, _, err := client.Must().CloudClient.NetworkInterfacesApi.DatacentersServersNicsGet(context.Background(),
			viper.GetString(core.GetFlagName(cmd.NS, constants.FlagDatacenterId)), *server.Id).Execute()
		if err != nil || nics.Items == nil || len(nics.Items) == 0 {
			return "", fmt.Errorf("failed getting nics %w", err)
		}
		// Find the first nic with IPs not empty and return it
		for _, nic := range nics.Items {
			if nic.Properties.Ips != nil && len(nic.Properties.Ips) > 0 {
				return (nic.Properties.Ips)[0], nil
			}
		}
	}
	return "", fmt.Errorf("no NIC with IP")
}

func generateCidrs(cmd *core.Command, ip string) []string {
	var cidrs []string

	instances := viper.GetInt(core.GetFlagName(cmd.NS, constants.FlagInstances))
	if instances == 0 {
		instances = 1 // flag might not be set
	}

	for i := 0; i < instances; i++ {
		cidrs = append(cidrs, fake.IP(fake.WithIPv4(), fake.WithIPCIDR(ip+"/24"))+"/24")
	}

	return cidrs
}
