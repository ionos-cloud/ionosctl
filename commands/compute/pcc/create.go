package pcc

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	cloudapiv6 "github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6"
)

func PccCreateCmd() *core.Command {
	cmd := core.NewCommand(context.TODO(), nil, core.CommandBuilder{
		Namespace:  "pcc",
		Resource:   "pcc",
		Verb:       "create",
		Aliases:    []string{"c"},
		ShortDesc:  "Create a Cross-Connect",
		LongDesc:   "Use this command to create a Cross-Connect. You can specify the name and the description for the Cross-Connect.",
		Example:    `ionosctl pcc create --name NAME --description DESCRIPTION --wait-for-request`,
		PreCmdRun:  core.NoPreRun,
		CmdRun:     RunPccCreate,
		InitClient: true,
	})
	cmd.AddStringFlag(cloudapiv6.ArgName, cloudapiv6.ArgNameShort, "Unnamed PrivateCrossConnect", "The name for the Cross-Connect")
	cmd.AddStringFlag(cloudapiv6.ArgDescription, cloudapiv6.ArgDescriptionShort, "", "The description for the Cross-Connect")
	cmd.AddBoolFlag(constants.ArgWaitForRequest, constants.ArgWaitForRequestShort, constants.DefaultWait, "Wait for the Request for Cross-Connect creation to be executed")
	cmd.AddIntFlag(constants.ArgTimeout, constants.ArgTimeoutShort, constants.DefaultTimeoutSeconds, "Timeout option for Request for Cross-Connect creation [seconds]")

	return cmd
}
