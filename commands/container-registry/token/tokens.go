package token

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

func TokenCmd() *core.Command {
	regCmd := &core.Command{
		Command: &cobra.Command{
			Use:     "token",
			Aliases: []string{"t", "tokens"},
			Short:   "Registry Tokens Operations",
			Long: "Manage container registries for storage of docker images and OCI compliant artifacts. " +
				"This operation is restricted to contract owner, admin, and users with 'accessAndManageRegistries' and " +
				"Share/Edit access permissions for the data center hosting the registry.",
			TraverseChildren: true,
		},
	}

	regCmd.Command.PersistentFlags().Bool(constants.ArgNoHeaders, false, "When using text output, don't print headers")

	regCmd.AddCommand(TokenListCmd())
	regCmd.AddCommand(TokenPostCmd())
	//regCmd.AddCommand(RegGetCmd())
	regCmd.AddCommand(TokenDeleteCmd())
	//regCmd.AddCommand(RegUpdateCmd())
	//regCmd.AddCommand(RegReplaceCmd())

	return regCmd
}

func getTokenPrint(
	resp *ionoscloud.APIResponse, c *core.CommandConfig, response *[]ionoscloud.TokenResponse,
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
				defaultHeaders := []string{"TokenId", "DisplayName", "ExpiryDate", "Status"}
				r.OutputJSON = response
				r.KeyValue = getTokensRows(response) // map header -> rows
				r.Columns = printer.GetHeaders(
					allCols, defaultHeaders, cols,
				) // headers
			} else {
				r.OutputJSON = response
				r.KeyValue = getTokensRows(response)
				postHeaders := []string{"DisplayName", "ExpiryDate", "Status"} // map header -> rows
				r.Columns = printer.GetHeaders(allCols, postHeaders, cols)     // headers
			}
		}
	}
	return r
}

type TokenPrint struct {
	TokenId             string `json:"TokenId,omitempty"`
	DisplayName         string `json:"DisplayName,omitempty"`
	ExpiryDate          string `json:"ExpiryDate,omitempty"`
	CredentialsUsername string `json:"CredentialsUsername,omitempty"`
	CredentialsPassword string `json:"CredentialsPassword,omitempty"`
	Status              string `json:"Status,omitempty"`
}

func getTokensRows(tokens *[]ionoscloud.TokenResponse) []map[string]interface{} {
	out := make([]map[string]interface{}, 0, len(*tokens))
	for _, token := range *tokens {
		var tokenPrint TokenPrint
		if idOk, ok := token.GetIdOk(); ok && idOk != nil {
			tokenPrint.TokenId = *idOk
		}
		if propertiesOk, ok := token.GetPropertiesOk(); ok && propertiesOk != nil {
			if displayNameOk, ok := propertiesOk.GetNameOk(); ok && displayNameOk != nil {
				tokenPrint.DisplayName = *displayNameOk
			}
			if expiryDateOk, ok := propertiesOk.GetExpiryDateOk(); ok && expiryDateOk != nil {
				tokenPrint.ExpiryDate = expiryDateOk.String()
			}
			if credentialsOk, ok := propertiesOk.GetCredentialsOk(); ok && credentialsOk != nil {
				if usernameOk, ok := credentialsOk.GetUsernameOk(); ok && usernameOk != nil {
					tokenPrint.CredentialsUsername = *usernameOk
				}
				if passwordOk, ok := credentialsOk.GetPasswordOk(); ok && passwordOk != nil {
					tokenPrint.CredentialsPassword = *passwordOk
				}
			}
			if statusOk, ok := propertiesOk.GetStatusOk(); ok && statusOk != nil {
				tokenPrint.Status = *statusOk
			}
		}
		o := structs.Map(tokenPrint)
		out = append(out, o)
	}
	return out
}

var allCols = structs.Names(TokenPrint{})

// TODO: check how other resources are doing this
func TokensIds() []string {
	client, _ := config.GetClient()
	svc := resources.NewRegistriesService(client, context.Background())
	regs, _, _ := svc.List("")

	svcToken := resources.NewTokenService(client, context.Background())
	regsIDs := *regs.GetItems()

	var allTokens []ionoscloud.TokenResponse

	for _, regID := range regsIDs {
		tokens, _, _ := svcToken.List(*regID.GetId())
		allTokens = append(allTokens, *tokens.GetItems()...)
	}

	return utils.Map(
		allTokens, func(i int, reg ionoscloud.TokenResponse) string {
			return *reg.GetId()
		},
	)
}
