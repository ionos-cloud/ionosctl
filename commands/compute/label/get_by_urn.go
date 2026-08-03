package label

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	cloudapiv6 "github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6"
)

func LabelGetByUrnCmd() *core.Command {
	cmd := core.NewCommand(context.TODO(), nil, core.CommandBuilder{
		Namespace: "label",
		Resource:  "label",
		Verb:      "get-by-urn",
		ShortDesc: "Get a Label using URN",
		LongDesc: `Use this command to fetch a single Label by its URN, without having to specify --resource-type and the resource id separately.

A label's URN is its globally unique identifier and encodes everything needed to locate it: urn:label:<resource_type>:<resource_id>:<key> (for example urn:label:server:<server-uuid>:env). You can find the URN in the "URN" column of ` + "`label list`" + ` / ` + "`label get`" + `.

Required values to run command:

* Label URN`,
		Example:    `ionosctl compute label get-by-urn --label-urn "urn:label:server:SERVER_ID:env"`,
		PreCmdRun:  PreRunLabelUrn,
		CmdRun:     RunLabelGetByUrn,
		InitClient: true,
	})
	cmd.AddStringFlag(cloudapiv6.ArgLabelUrn, "", "", "The label's globally unique URN, formatted urn:label:<resource_type>:<resource_id>:<key> (e.g. urn:label:server:SERVER_ID:env)", core.RequiredFlagOption())

	return cmd
}
