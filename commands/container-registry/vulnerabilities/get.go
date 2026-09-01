package vulnerabilities

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/spf13/viper"
)

func VulnerabilitiesGetCmd() *core.Command {
	cmd := core.NewCommand(
		context.TODO(), nil, core.CommandBuilder{
			Namespace:  "container-registry",
			Resource:   "vulnerabilities",
			Verb:       "get",
			ShortDesc:  "Get a single vulnerability finding by ID",
			LongDesc:   "Get the full details of one vulnerability finding by its ID (e.g. a CVE identifier), including CVSS score, severity, affected packages/versions, remediation recommendations and references. Find IDs via 'container-registry vulnerabilities list'.",
			Example:    "ionosctl container-registry vulnerabilities get --vulnerability-id VULNERABILITY_ID",
			PreCmdRun:  PreCmdGet,
			CmdRun:     CmdGet,
			InitClient: true,
		},
	)

	cmd.AddStringFlag(constants.FlagVulnerabilityId, "", "", "The ID of the vulnerability finding to retrieve (e.g. a CVE identifier, as shown by 'vulnerabilities list')")

	return cmd
}

func PreCmdGet(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlags(c.Command, c.NS, constants.FlagVulnerabilityId)
}

func CmdGet(c *core.CommandConfig) error {
	vulnId := viper.GetString(core.GetFlagName(c.NS, constants.FlagVulnerabilityId))

	vulnerability, _, err := client.Must().RegistryClient.VulnerabilitiesApi.VulnerabilitiesFindByID(
		context.
			Background(), vulnId,
	).Execute()
	if err != nil {
		return err
	}

	return c.Printer(allCols).Print(vulnerability)
}
