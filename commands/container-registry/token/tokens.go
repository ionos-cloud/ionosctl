package token

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/pkg/functional"

	scope "github.com/ionos-cloud/ionosctl/v6/commands/container-registry/token/scopes"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/services/container-registry/resources"
	"github.com/ionos-cloud/sdk-go-bundle/products/containerregistry/v2"
	"github.com/spf13/cobra"
)

var (
	postHeaders  = []string{"CredentialsPassword"}
	AllTokenCols = []string{"TokenId", "DisplayName", "ExpiryDate", "CredentialsUsername", "CredentialsPassword", "Status", "RegistryId"}
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
	svcToken := resources.NewTokenService(client.Must(), context.Background())
	var allTokens []containerregistry.TokenResponse

	if regId != "" {

		tokens, _, _ := svcToken.List(regId)

		allTokens = append(allTokens, tokens.GetItems()...)

		return functional.Map(
			allTokens, func(reg containerregistry.TokenResponse) string {
				return reg.GetId()
			},
		)
	}

	regs, _, _ := resources.NewRegistriesService(client.Must(), context.Background()).List("")
	regsIDs := regs.GetItems()

	for _, regID := range regsIDs {
		tokens, _, _ := svcToken.List(regID.GetId())

		allTokens = append(allTokens, tokens.GetItems()...)
	}

	return functional.Map(
		allTokens, func(reg containerregistry.TokenResponse) string {
			return reg.GetId()
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
