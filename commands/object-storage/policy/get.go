package policy

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/spf13/viper"

	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/table"
)

type statementInfo struct {
	Sid       string `json:"Sid"`
	Effect    string `json:"Effect"`
	Action    string `json:"Action"`
	Resource  string `json:"Resource"`
	Principal string `json:"Principal"`
	Condition string `json:"Condition"`
}

func GetCmd() *core.Command {
	cmd := core.NewCommand(context.Background(), nil, core.CommandBuilder{
		Namespace: "object-storage",
		Resource:  "policy",
		Verb:      "get",
		Aliases:   []string{"g"},
		ShortDesc: "Get the bucket policy",
		Example:   "ionosctl object-storage policy get --name my-bucket",
		PreCmdRun: func(c *core.PreCommandConfig) error {
			return core.CheckRequiredFlags(c.Command, c.NS, constants.FlagName)
		},
		CmdRun: func(c *core.CommandConfig) error {
			name := viper.GetString(core.GetFlagName(c.NS, constants.FlagName))

			s3, _, err := client.GetRegionalObjectStorageClient(context.Background(), name)
			if err != nil {
				return err
			}

			result, _, err := s3.PolicyApi.GetBucketPolicy(context.Background(), name).Execute()
			if err != nil {
				return err
			}

			var statements []statementInfo
			for _, s := range result.GetStatement() {
				si := statementInfo{
					Sid:    s.GetSid(),
					Effect: s.GetEffect(),
					Action: strings.Join(s.GetAction(), ", "),
				}

				si.Resource = strings.Join(s.GetResource(), ", ")

				if s.HasPrincipal() {
					p := s.GetPrincipal()
					si.Principal = strings.Join(p.GetAWS(), ", ")
				}

				if s.HasCondition() {
					cond := s.GetCondition()
					b, err := json.Marshal(cond)
					if err == nil {
						si.Condition = string(b)
					}
				}

				statements = append(statements, si)
			}

			cols, _ := c.Command.Command.Flags().GetStringSlice(constants.ArgCols)
			return c.Out(table.Sprint(allCols, statements, cols))
		},
		InitClient: false,
	})

	cmd.AddStringFlag(constants.FlagName, constants.FlagNameShort, "", "Name of the bucket", core.RequiredFlagOption())

	cmd.Command.SilenceUsage = true
	cmd.Command.Flags().SortFlags = false
	return cmd
}
