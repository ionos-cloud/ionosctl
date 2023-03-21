package token

import (
	"context"
	"fmt"
	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"strconv"
	"strings"
	"time"

	"github.com/ionos-cloud/ionosctl/v6/internal/functional"

	"github.com/fatih/structs"
	scope "github.com/ionos-cloud/ionosctl/v6/commands/container-registry/token/scopes"
	"github.com/ionos-cloud/ionosctl/v6/pkg/constants"
	"github.com/ionos-cloud/ionosctl/v6/pkg/core"
	"github.com/ionos-cloud/ionosctl/v6/pkg/printer"
	"github.com/ionos-cloud/ionosctl/v6/services/container-registry/resources"
	ionoscloud "github.com/ionos-cloud/sdk-go-container-registry"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func TokenCmd() *core.Command {
	tokenCmd := &core.Command{
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

	tokenCmd.Command.PersistentFlags().Bool(
		constants.ArgNoHeaders, false, "When using text output, don't print headers",
	)

	tokenCmd.AddCommand(TokenListCmd())
	tokenCmd.AddCommand(TokenPostCmd())
	tokenCmd.AddCommand(TokenGetCmd())
	tokenCmd.AddCommand(TokenDeleteCmd())
	tokenCmd.AddCommand(TokenUpdateCmd())
	tokenCmd.AddCommand(TokenReplaceCmd())
	tokenCmd.AddCommand(scope.TokenScopesCmd())

	return tokenCmd
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

func TokensIds(regId string) []string {
	svcToken := resources.NewTokenService(client.Must(), context.Background())
	var allTokens []ionoscloud.TokenResponse

	if regId != "" {

		tokens, _, _ := svcToken.List(regId)

		allTokens = append(allTokens, *tokens.GetItems()...)

		return functional.Map(
			allTokens, func(reg ionoscloud.TokenResponse) string {
				return *reg.GetId()
			},
		)
	}

	regs, _, _ := resources.NewRegistriesService(client.Must(), context.Background()).List("")
	regsIDs := *regs.GetItems()

	for _, regID := range regsIDs {
		tokens, _, _ := svcToken.List(*regID.GetId())

		allTokens = append(allTokens, *tokens.GetItems()...)
	}

	return functional.Map(
		allTokens, func(reg ionoscloud.TokenResponse) string {
			return *reg.GetId()
		},
	)
}

func ParseExpiryTime(expiryTime string) (time.Duration, error) {
	years := 0
	months := 0
	days := 0
	hours := 0

	if !strings.ContainsAny(expiryTime, "0123456789") {
		return 0, fmt.Errorf("invalid expiry time format")
	}

	number := ""

	for i := 0; i < len(expiryTime); i++ {
		if string(expiryTime[i]) != "y" && string(expiryTime[i]) != "m" && string(expiryTime[i]) != "d" && string(expiryTime[i]) != "h" {
			number += string(expiryTime[i])
		} else if expiryTime[i] == 'y' {
			years, _ = strconv.Atoi(number)
			number = ""
		} else if expiryTime[i] == 'm' {
			months, _ = strconv.Atoi(number)
			number = ""
		} else if expiryTime[i] == 'd' {
			days, _ = strconv.Atoi(number)
			number = ""
		} else if expiryTime[i] == 'h' {
			hours, _ = strconv.Atoi(number)
			number = ""
		}
	}

	return time.Duration(years)*time.Hour*24*365 + time.Duration(months)*time.Hour*24*30 + time.Duration(days)*time.Hour*24 + time.Duration(hours)*time.Hour, nil
}
