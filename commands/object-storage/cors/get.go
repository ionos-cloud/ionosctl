package cors

import (
	"context"
	"strconv"
	"strings"

	"github.com/spf13/viper"

	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/table"
)

func GetCmd() *core.Command {
	cmd := core.NewCommand(context.Background(), nil, core.CommandBuilder{
		Namespace: "object-storage",
		Resource:  "cors",
		Verb:      "get",
		Aliases:   []string{"g"},
		ShortDesc: "Get the CORS configuration for a bucket",
		Example:   "ionosctl object-storage cors get --name my-bucket",
		PreCmdRun: func(c *core.PreCommandConfig) error {
			return core.CheckRequiredFlags(c.Command, c.NS, constants.FlagName)
		},
		CmdRun: func(c *core.CommandConfig) error {
			name := viper.GetString(core.GetFlagName(c.NS, constants.FlagName))

			s3, err := resolveRegionalClient(context.Background(), name)
			if err != nil {
				return err
			}

			result, _, err := s3.CORSApi.GetBucketCors(context.Background(), name).Execute()
			if err != nil {
				return err
			}

			var rules []corsRuleInfo
			for _, r := range result.GetCORSRules() {
				info := corsRuleInfo{
					AllowedOrigins: strings.Join(r.GetAllowedOrigins(), ", "),
					AllowedMethods: strings.Join(r.GetAllowedMethods(), ", "),
					AllowedHeaders: strings.Join(r.GetAllowedHeaders(), ", "),
					ExposeHeaders:  strings.Join(r.GetExposeHeaders(), ", "),
				}

				if r.HasMaxAgeSeconds() {
					info.MaxAgeSeconds = strconv.FormatInt(int64(r.GetMaxAgeSeconds()), 10)
				}

				if r.HasID() {
					info.ID = strconv.FormatInt(int64(r.GetID()), 10)
				}

				rules = append(rules, info)
			}

			cols, _ := c.Command.Command.Flags().GetStringSlice(constants.ArgCols)
			return c.Out(table.Sprint(allCols, rules, cols))
		},
		InitClient: false,
	})

	cmd.AddStringFlag(constants.FlagName, constants.FlagNameShort, "", "Name of the bucket", core.RequiredFlagOption())

	cmd.Command.SilenceUsage = true
	cmd.Command.Flags().SortFlags = false
	return cmd
}
