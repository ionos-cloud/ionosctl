package token

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/table"
	"github.com/ionos-cloud/ionosctl/v6/pkg/functional"

	scope "github.com/ionos-cloud/ionosctl/v6/commands/container-registry/token/scopes"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/sdk-go-bundle/products/containerregistry/v2"
	"github.com/spf13/cobra"
)

var allCols = []table.Column{
	{Name: "TokenId", JSONPath: "id", Default: true},
	{Name: "DisplayName", JSONPath: "properties.name", Default: true},
	{Name: "ExpiryDate", JSONPath: "properties.expiryDate", Default: true},
	{Name: "CredentialsUsername", JSONPath: "properties.credentials.username", Default: true},
	{Name: "CredentialsPassword", JSONPath: "properties.credentials.password", Default: true},
	{Name: "Status", JSONPath: "properties.status", Default: true},
	{Name: "RegistryId", JSONPath: "href", Default: true},
}

// postCols defines the columns shown after create/replace operations (includes CredentialsPassword).
var postCols = []table.Column{
	{Name: "CredentialsPassword", JSONPath: "properties.credentials.password", Default: true},
}

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

	tokenCmd.Command.PersistentFlags().StringSlice(constants.ArgCols, nil, table.ColsMessage(allCols))
	_ = tokenCmd.Command.RegisterFlagCompletionFunc(
		constants.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return table.AllCols(allCols), cobra.ShellCompDirectiveNoFileComp
		},
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

func TokensIds(regId string) []string {
	var allTokens []containerregistry.TokenResponse

	if regId != "" {
		tokens, _, _ := client.Must().RegistryClient.TokensApi.RegistriesTokensGet(context.Background(), regId).Execute()
		allTokens = append(allTokens, tokens.GetItems()...)
		return functional.Map(allTokens, func(t containerregistry.TokenResponse) string { return t.GetId() })
	}

	regs, _, _ := client.Must().RegistryClient.RegistriesApi.RegistriesGet(context.Background()).Execute()
	for _, reg := range regs.GetItems() {
		toks, _, _ := client.Must().RegistryClient.TokensApi.RegistriesTokensGet(context.Background(), reg.GetId()).Execute()
		allTokens = append(allTokens, toks.GetItems()...)
	}
	return functional.Map(allTokens, func(t containerregistry.TokenResponse) string { return t.GetId() })
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
