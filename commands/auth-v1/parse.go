package authv1

import (
	"context"

	"github.com/fatih/structs"
	"github.com/ionos-cloud/ionosctl/v6/internal/jwt"
	"github.com/ionos-cloud/ionosctl/v6/pkg/constants"
	"github.com/ionos-cloud/ionosctl/v6/pkg/core"
	"github.com/ionos-cloud/ionosctl/v6/pkg/printer"
	authv1 "github.com/ionos-cloud/ionosctl/v6/services/auth-v1"
	"github.com/spf13/viper"
)

func TokenParseCmd() *core.Command {
	cmd := core.NewCommand(context.Background(), nil, core.CommandBuilder{
		Namespace:  "token",
		Resource:   "token",
		Verb:       "parse",
		Aliases:    []string{"p"},
		ShortDesc:  "Parse the contents of a Token",
		LongDesc:   "", // TODO: write description
		Example:    "", // TODO: prep examples,
		PreCmdRun:  preRunTokenParse,
		CmdRun:     runTokenParse,
		InitClient: false,
	})
	cmd.AddStringFlag(authv1.ArgToken, authv1.ArgTokenShort, "", authv1.Token, core.RequiredFlagOption())

	return cmd
}

func preRunTokenParse(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlags(c.Command, c.NS, constants.ArgToken)
}

// TODO: implement parsing
func runTokenParse(c *core.CommandConfig) error {
	token := viper.GetString(core.GetFlagName(c.NS, authv1.ArgToken))

	headers, err := jwt.Headers(token)
	if err != nil {
		return err
	}

	claims, err := jwt.Claims(token)
	if err != nil {
		return err
	}

	var tokenInfo tokenInfoPrint

	tokenInfo.TokenId, err = jwt.Kid(headers)
	if err != nil {
		return err
	}

	tokenInfo.UserId, err = jwt.Uuid(claims)
	if err != nil {
		return err
	}

	tokenInfo.ContractNumber, err = jwt.ContractNumber(claims)
	if err != nil {
		return err
	}

	tokenInfo.Role, err = jwt.Role(claims)
	if err != nil {
		return err
	}

	return c.Printer.Print(getTokenInfoPrint(c, tokenInfo))
}

// Token content printing

type tokenInfoPrint struct {
	TokenId        string `json:"TokenId,omitempty"`
	UserId         string `json:"UserId,omitempty"`
	ContractNumber int64  `json:"ContractNumber,omitempty"`
	Role           string `json:"Role,omitempty"`
}

var allCols = structs.Names(tokenInfoPrint{})

func getTokenInfoPrint(c *core.CommandConfig, tokenInfo tokenInfoPrint) printer.Result {
	r := printer.Result{}
	cols, _ := c.Command.Command.Flags().GetStringSlice(constants.ArgCols)

	if c != nil {
		r.OutputJSON = tokenInfo
		r.Columns = printer.GetHeadersAllDefault(allCols, cols)
		r.KeyValue = makeTokenInfoPrintObject(tokenInfo)
	}

	return r
}

func makeTokenInfoPrintObject(tokenInfo tokenInfoPrint) []map[string]interface{} {
	var out = make([]map[string]interface{}, 0)

	out = append(out, structs.Map(tokenInfo))
	return out
}
