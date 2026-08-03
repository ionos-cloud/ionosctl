package pcc

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	cloudapiv6 "github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6"
)

func PccCreateCmd() *core.Command {
	cmd := core.NewCommand(context.TODO(), nil, core.CommandBuilder{
		Namespace: "pcc",
		Resource:  "pcc",
		Verb:      "create",
		Aliases:   []string{"c"},
		ShortDesc: "Create a Cross-Connect",
		LongDesc: `Use this command to create a Cross-Connect (Private Cross-Connect).

This creates the empty Cross-Connect object only; no LAN is attached yet. After creation, attach LANs by running ` + "`" + `ionosctl compute lan update --datacenter-id <DC_ID> --lan-id <LAN_ID> --pcc-id <PCC_ID>` + "`" + ` for each private LAN you want to peer. All LANs you attach must be in the same contract and region and must have non-overlapping private IP ranges.`,
		Example: `# Create a Cross-Connect with a name and description
ionosctl compute pcc create --name "prod-interconnect" --description "Bridges DB and app VDCs"

# Create and wait for the request to complete
ionosctl compute pcc create --name "prod-interconnect" --wait`,
		PreCmdRun:  core.NoPreRun,
		CmdRun:     RunPccCreate,
		InitClient: true,
	})
	cmd.AddStringFlag(cloudapiv6.ArgName, cloudapiv6.ArgNameShort, "Unnamed PrivateCrossConnect", "The name for the Cross-Connect. Purely a display label; does not affect connectivity")
	cmd.AddStringFlag(cloudapiv6.ArgDescription, cloudapiv6.ArgDescriptionShort, "", "A human-readable description for the Cross-Connect")

	return cmd
}
