package registry

import (
	"context"
	"github.com/fatih/structs"
	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/functional"
	"github.com/ionos-cloud/ionosctl/v6/pkg/constants"
	"github.com/ionos-cloud/ionosctl/v6/pkg/core"
	"github.com/ionos-cloud/ionosctl/v6/pkg/printer"
	"github.com/ionos-cloud/ionosctl/v6/services/container-registry/resources"
	ionoscloud "github.com/ionos-cloud/sdk-go-container-registry"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func RegistryCmd() *core.Command {
	regCmd := &core.Command{
		Command: &cobra.Command{
			Use:     "registry",
			Aliases: []string{"reg", "registries", "r"},
			Short:   "Registry Operations",
			Long: "Manage container registries for storage of docker images and OCI compliant artifacts. " +
				"This operation is restricted to contract owner, admin, and users with 'accessAndManageRegistries' and " +
				"Share/Edit access permissions for the data center hosting the registry.",
			TraverseChildren: true,
		},
	}

	regCmd.Command.PersistentFlags().Bool(constants.ArgNoHeaders, false, "When using text output, don't print headers")

	regCmd.AddCommand(RegListCmd())
	regCmd.AddCommand(RegPostCmd())
	regCmd.AddCommand(RegGetCmd())
	regCmd.AddCommand(RegDeleteCmd())
	regCmd.AddCommand(RegUpdateCmd())
	regCmd.AddCommand(RegReplaceCmd())

	return regCmd
}

func getRegistryPrint(
	resp *ionoscloud.APIResponse, c *core.CommandConfig, response *[]ionoscloud.RegistryResponse,
	post bool,
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
			if !post {
				r.OutputJSON = response
				r.KeyValue = getRegRows(response)                       // map header -> rows
				r.Columns = printer.GetHeadersAllDefault(allCols, cols) // headers
			} else {
				r.OutputJSON = response
				r.KeyValue = getRegRows(response)
				postHeaders := []string{
					"DisplayName", "Location",
				} // map header -> rows
				r.Columns = printer.GetHeaders(allCols, postHeaders, cols) // headers
			}
		}
	}
	return r
}

type RegPrint struct {
	RegistryId            string `json:"RegistryId,omitempty"`
	DisplayName           string `json:"DisplayName,omitempty"`
	Location              string `json:"Location,omitempty"`
	Hostname              string `json:"Hostname,omitempty"`
	GarbageCollectionDays string `json:"GarbageCollectionDays,omitempty"`
	GarbageCollectionTime string `json:"GarbageCollectionTime,omitempty"`
}

func getRegRows(regs *[]ionoscloud.RegistryResponse) []map[string]interface{} {
	out := make([]map[string]interface{}, 0, len(*regs))
	for _, reg := range *regs {
		var regPrint RegPrint
		if idOk, ok := reg.GetIdOk(); ok && idOk != nil {
			regPrint.RegistryId = *idOk
		}
		if propertiesOk, ok := reg.GetPropertiesOk(); ok && propertiesOk != nil {
			if displayNameOk, ok := propertiesOk.GetNameOk(); ok && displayNameOk != nil {
				regPrint.DisplayName = *displayNameOk
			}
			if locationOk, ok := propertiesOk.GetLocationOk(); ok && locationOk != nil {
				regPrint.Location = *locationOk
			}
			if hostnameOk, ok := propertiesOk.GetHostnameOk(); ok && hostnameOk != nil {
				regPrint.Hostname = *hostnameOk
			}
			if gcOk, ok := propertiesOk.GetGarbageCollectionScheduleOk(); ok && gcOk != nil {
				if gcDaysOk, ok := gcOk.GetDaysOk(); ok && gcDaysOk != nil {
					for _, day := range *gcDaysOk {
						regPrint.GarbageCollectionDays += string(day) + ", "
					}
					regPrint.GarbageCollectionDays = regPrint.GarbageCollectionDays[:len(regPrint.GarbageCollectionDays)-2]
				}
				if gcTimeOk, ok := gcOk.GetTimeOk(); ok && gcTimeOk != nil {
					regPrint.GarbageCollectionTime = *gcTimeOk
				}
			}
		}
		o := structs.Map(regPrint)
		out = append(out, o)
	}
	return out
}

var allCols = structs.Names(RegPrint{})

func RegsIds() []string {
	//client, _ := config.GetClient()
	svc := resources.NewRegistriesService(client.Must(), context.Background())
	regs, _, _ := svc.List("")
	return functional.Map(
		*regs.GetItems(), func(reg ionoscloud.RegistryResponse) string {
			return *reg.GetId()
		},
	)
}
