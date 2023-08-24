package token

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/pkg/json2table"

	"github.com/ionos-cloud/ionosctl/v6/internal/functional"

	scope "github.com/ionos-cloud/ionosctl/v6/commands/container-registry/token/scopes"
	"github.com/ionos-cloud/ionosctl/v6/pkg/constants"
	"github.com/ionos-cloud/ionosctl/v6/pkg/core"
	"github.com/ionos-cloud/ionosctl/v6/services/container-registry/resources"
	ionoscloud "github.com/ionos-cloud/sdk-go-container-registry"
	"github.com/spf13/cobra"
)

var (
	allJSONPaths = map[string]string{
		"TokenId":     "id",
		"DisplayName": "properties.name",
		//"ExpiryDate":          "properties.expiryDate",
		"CredentialsUsername": "properties.credentials.username",
		"CredentialsPassword": "properties.credential.password",
		"Status":              "properties.status",
	}

	postHeaders  = []string{"CredentialsPassword"}
	AllTokenCols = []string{"TokenId", "DisplayName", "ExpiryDate", "CredentialsUsername", "CredentialsPassword", "Status"}
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

func ConvertTokensToTable(tokens ionoscloud.TokensResponse) ([]map[string]interface{}, error) {
	items, ok := tokens.GetItemsOk()
	if !ok || items == nil {
		return nil, fmt.Errorf("could not retrieve Container Registry Token items")
	}

	var tokensConverted []map[string]interface{}
	for _, item := range *items {
		temp, err := ConvertTokenToTable(item)
		if err != nil {
			return nil, err
		}

		tokensConverted = append(tokensConverted, temp...)
	}

	return tokensConverted, nil
}

func ConvertTokenToTable(token ionoscloud.TokenResponse) ([]map[string]interface{}, error) {
	properties, ok := token.GetPropertiesOk()
	if !ok || properties == nil {
		return nil, fmt.Errorf("could not retrieve Container Registry Token properties")
	}

	expiryDate, ok := properties.GetExpiryDateOk()
	if !ok || expiryDate == nil {
		return nil, fmt.Errorf("could not retrieve Container Registry Token Expiry Date")
	}

	temp, err := json2table.ConvertJSONToTable("", allJSONPaths, token)
	if err != nil {
		return nil, fmt.Errorf("could not convert from JSON to Table format: %w", err)
	}

	temp[0]["ExpiryDate"] = expiryDate.String()
	return temp, nil
}
