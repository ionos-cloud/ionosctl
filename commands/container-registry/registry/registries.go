package registry

import (
	"context"
	"github.com/fatih/structs"
	"github.com/ionos-cloud/ionosctl/pkg/config"
	"github.com/ionos-cloud/ionosctl/pkg/constants"
	"github.com/ionos-cloud/ionosctl/pkg/core"
	"github.com/ionos-cloud/ionosctl/pkg/printer"
	"github.com/ionos-cloud/ionosctl/pkg/utils"
	"github.com/ionos-cloud/ionosctl/services/container-registry/resources"
	ionoscloud "github.com/ionos-cloud/sdk-go-container-registry"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func RegistryCmd() *core.Command {
	regCmd := &core.Command{
		Command: &cobra.Command{
			Use:   "registry",
			Short: "Registry Operations",
			Long: "Manage container registries for storage of docker images and OCI compliant artifacts. " +
				"This operation is restricted to contract owner, admin, and users with 'accessAndManageRegistries' and " +
				"Share/Edit access permissions for the data center hosting the registry.",
			TraverseChildren: true,
		},
	}

	regCmd.Command.PersistentFlags().Bool(constants.ArgNoHeaders, false, "When using text output, don't print headers")

	regCmd.AddCommand(RegListCmd())

	return regCmd
}

func getRegistryPrint(
	resp *ionoscloud.APIResponse, c *core.CommandConfig, response *[]ionoscloud.RegistryResponse,
) printer.Result {
	r := printer.Result{}
	cols, _ := c.Command.Command.Flags().GetStringSlice(constants.ArgCols)

	if c != nil {
		if resp != nil {
			r.Resource = c.Resource
			r.Verb = c.Verb
			r.WaitForState = viper.GetBool(
				core.GetFlagName(
					c.NS, constants.ArgWaitForRequest,
				),
			) // this boolean is duplicated everywhere just to do an append of `& wait` to a verbose message
		}
		if response != nil {
			r.OutputJSON = response
			r.KeyValue = getRegRows(response)                       // map header -> rows
			r.Columns = printer.GetHeadersAllDefault(allCols, cols) // headers
		}
	}
	return r
}

type RegPrint struct {
	RegistryId  string `json:"RegistryId,omitempty"`
	DisplayName string `json:"DisplayName,omitempty"`
	Location    string `json:"Location,omitempty"`
	Hostname    string `json:"DisplayName,omitempty"`
}

func getRegRows(regs *[]ionoscloud.RegistryResponse) []map[string]interface{} {
	out := make([]map[string]interface{}, 0, len(*regs))
	for _, cert := range *regs {
		var regPrint RegPrint
		if idOk, ok := cert.GetIdOk(); ok && idOk != nil {
			regPrint.RegistryId = *idOk
		}
		if propertiesOk, ok := cert.GetPropertiesOk(); ok && propertiesOk != nil {
			if displayNameOk, ok := propertiesOk.GetNameOk(); ok && displayNameOk != nil {
				regPrint.DisplayName = *displayNameOk
			}
			if locationOk, ok := propertiesOk.GetLocationOk(); ok && locationOk != nil {
				regPrint.Location = *locationOk
			}
			if hostnameOk, ok := propertiesOk.GetHostnameOk(); ok && hostnameOk != nil {
				regPrint.Hostname = *hostnameOk
			}
		}
		o := structs.Map(regPrint)
		out = append(out, o)
	}
	return out
}

var allCols = structs.Names(RegPrint{})

func RegsIds() []string {
	client, _ := config.GetClient()
	svc := resources.NewRegistriesService(client, context.Background())
	certs, _, _ := svc.List("")
	return utils.Map(
		*certs.GetItems(), func(i int, reg ionoscloud.RegistryResponse) string {
			return *reg.GetId()
		},
	)
}
