package certmanager

import (
	"context"
	"github.com/ionos-cloud/ionosctl/pkg/constants"
	"github.com/ionos-cloud/ionosctl/pkg/core"
)

func CertListCmd() *core.Command {
	cmd := core.NewCommand(context.TODO(), nil, core.CommandBuilder{
		Namespace:  "certmanager",
		Resource:   "certificates",
		Verb:       "list",
		Aliases:    []string{"l"},
		ShortDesc:  "List all Certificates",
		LongDesc:   "Use this command to retrieve all Certificates.",
		Example:    "ionosctl certificate-manager list",
		PreCmdRun:  core.NoPreRun,
		CmdRun:     CmdList,
		InitClient: true,
	})

	cmd.AddBoolFlag(constants.ArgNoHeaders, "", false, "Get response without headers")

	return cmd
}

func CmdList(c *core.CommandConfig) error {
	c.Printer.Verbose("Getting Certificates...")
	cert, _, err := c.CertificateManagerServices.Certs().List()
	if err != nil {
		return err
	}
	list := cert.GetItems()
	return c.Printer.Print(getCertPrint(nil, c, list))
}
