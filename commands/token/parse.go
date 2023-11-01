package token

import (
	"context"
	"fmt"

	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/jwt"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/json2table/jsonpaths"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/jsontabwriter"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/tabheaders"
	"github.com/ionos-cloud/ionosctl/v6/pkg/core"
	authservice "github.com/ionos-cloud/ionosctl/v6/services/auth-v1"
	"github.com/spf13/viper"
)

func TokenParseCmd() *core.Command {
	cmd := core.NewCommand(context.Background(), nil, core.CommandBuilder{
		Namespace: "token",
		Resource:  "token",
		Verb:      "parse",
		Aliases:   []string{"p"},
		ShortDesc: "Parse the contents of a Token",
		LongDesc: `Use this command to parse a Token and find out Token ID, User ID, Contract Number, Role.
If you want to view the privileges associated with the token, you must set the --privileges flag. When this flag is set, the command will output a list of privileges instead of the default output.

Required values to run:

* Token`,
		Example:    parseTokenExample,
		PreCmdRun:  preRunTokenParse,
		CmdRun:     runTokenParse,
		InitClient: false,
	})
	cmd.AddStringFlag(authservice.ArgToken, authservice.ArgTokenShort, "", authservice.Token, core.RequiredFlagOption())
	cmd.AddBoolFlag(authservice.ArgPrivileges, authservice.ArgPrivilegesShort, false, authservice.Privileges)

	return cmd
}

func preRunTokenParse(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlags(c.Command, c.NS, constants.ArgToken)
}

func runTokenParse(c *core.CommandConfig) error {
	cols := viper.GetStringSlice(core.GetFlagName(c.Resource, constants.ArgCols))

	token := viper.GetString(core.GetFlagName(c.NS, authservice.ArgToken))

	if viper.IsSet(core.GetFlagName(c.NS, authservice.ArgPrivileges)) {
		claims, err := jwt.Claims(token)
		if err != nil {
			return err
		}

		privileges, err := jwt.Privileges(claims)
		if err != nil {
			return err
		}

		privilegesConverted := makeTokenPrivilegesPrintObject(privileges)

		out, err := jsontabwriter.GenerateOutputPreconverted(privileges, privilegesConverted,
			tabheaders.GetHeadersAllDefault([]string{"Privileges"}, cols))
		if err != nil {
			return err
		}

		fmt.Fprintf(c.Command.Command.OutOrStdout(), out)
		return nil
	}

	headers, err := jwt.Headers(token)
	if err != nil {
		return err
	}

	claims, err := jwt.Claims(token)
	if err != nil {
		return err
	}

	var info TokenInfo

	info.TokenId, err = jwt.Kid(headers)
	if err != nil {
		return err
	}

	info.UserId, err = jwt.Uuid(claims)
	if err != nil {
		return err
	}

	info.ContractNumber, err = jwt.ContractNumber(claims)
	if err != nil {
		return err
	}

	info.Role, err = jwt.Role(claims)
	if err != nil {
		return err
	}

	out, err := jsontabwriter.GenerateOutput("", jsonpaths.AuthTokenInfo, info, tabheaders.GetHeadersAllDefault(allCols, cols))
	if err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), out)
	return nil
}

type TokenInfo struct {
	TokenId        string `json:"tokenId,omitempty"`
	UserId         string `json:"userId,omitempty"`
	ContractNumber int64  `json:"contractNumber,omitempty"`
	Role           string `json:"role,omitempty"`
}

var (
	allCols = []string{"TokenId", "UserId", "ContractNumber", "Role"}
)

func makeTokenPrivilegesPrintObject(privileges []string) []map[string]interface{} {
	var out = make([]map[string]interface{}, 0)

	for _, priv := range privileges {
		temp := make(map[string]interface{}, 0)
		temp["Privileges"] = priv
		out = append(out, temp)
	}

	return out
}
