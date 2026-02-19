package label

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	cloudapiv6 "github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6"
)

func LabelGetByUrnCmd() *core.Command {
	cmd := core.NewCommand(context.TODO(), nil, core.CommandBuilder{
		Namespace:  "label",
		Resource:   "label",
		Verb:       "get-by-urn",
		ShortDesc:  "Get a Label using URN",
		LongDesc:   "Use this command to get information about a specified Label using its URN. A URN is used for uniqueness of a Label and composed using `urn:label:<resource_type>:<resource_uuid>:<key>`.\n\nRequired values to run command:\n\n* Label URN",
		Example:    `ionosctl compute label get-by-urn --label-urn "urn:label:server:SERVER_ID:test"`,
		PreCmdRun:  PreRunLabelUrn,
		CmdRun:     RunLabelGetByUrn,
		InitClient: true,
	})
	cmd.AddStringFlag(cloudapiv6.ArgLabelUrn, "", "", "URN for the Label [urn:label:<resource_type>:<resource_uuid>:<key>]", core.RequiredFlagOption())

	return cmd
}
