package vpn

import (
	"github.com/ionos-cloud/ionosctl/v6/commands/vpn/ipsec"

	"github.com/ionos-cloud/ionosctl/v6/commands/vpn/wireguard"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/spf13/cobra"
)

func Root() *core.Command {
	cmd := &core.Command{
		Command: &cobra.Command{
			Use:              "vpn",
			Short:            "VPN Operations",
			TraverseChildren: true,
		},
	}

	cmd.AddCommand(wireguard.Root())
	cmd.AddCommand(ipsec.Root())

	return core.WithRegionalFlags(cmd, constants.DefaultVPNApiURL, LocationToURL)
}

var LocationToURL = map[string]string{
	"de/fra": "https://vpn.de-fra.ionos.com",
	"de/txl": "https://vpn.de-txl.ionos.com",
	"es/vit": "https://vpn.es-vit.ionos.com",
	"gb/bhx": "https://vpn.gb-bhx.ionos.com",
	"gb/lhr": "https://vpn.gb-lhr.ionos.com",
	"us/ewr": "https://vpn.us-ewr.ionos.com",
	"us/las": "https://vpn.us-las.ionos.com",
	"us/mci": "https://vpn.us-mci.ionos.com",
	"fr/par": "https://vpn.fr-par.ionos.com",
}
