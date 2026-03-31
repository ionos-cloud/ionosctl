package bucket

import (
	"context"

	"github.com/spf13/viper"

	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/table"
)

var versioningCols = []table.Column{
	{Name: "Name", JSONPath: "Name", Default: true},
	{Name: "Versioning", JSONPath: "Versioning", Default: true},
}

type versioningResult struct {
	Name       string `json:"Name"`
	Versioning string `json:"Versioning"`
}

func GetBucketVersioningCmd() *core.Command {
	cmd := core.NewCommand(context.Background(), nil, core.CommandBuilder{
		Namespace: "object-storage",
		Resource:  "bucket",
		Verb:      "get-versioning",
		Aliases:   []string{"gv"},
		ShortDesc: "Get the versioning state of a bucket",
		Example:   "ionosctl object-storage bucket get-versioning --name my-bucket",
		PreCmdRun: func(c *core.PreCommandConfig) error {
			return core.CheckRequiredFlags(c.Command, c.NS, constants.FlagName)
		},
		CmdRun: func(c *core.CommandConfig) error {
			name := viper.GetString(core.GetFlagName(c.NS, constants.FlagName))

			s3, err := client.GetObjectStorageClient("")
			if err != nil {
				return err
			}

			// Resolve the bucket's actual region to avoid redirect loops.
			loc, _, err := s3.BucketsApi.GetBucketLocation(context.Background(), name).Execute()
			if err != nil {
				return err
			}

			region := ""
			if loc != nil {
				region = loc.GetLocationConstraint()
			}

			s3Regional, err := client.GetObjectStorageClient(region)
			if err != nil {
				return err
			}

			apiResult, _, err := s3Regional.VersioningApi.GetBucketVersioning(context.Background(), name).Execute()
			if err != nil {
				return err
			}

			status := "Disabled"
			if apiResult.HasStatus() {
				status = string(apiResult.GetStatus())
			}

			result := versioningResult{
				Name:       name,
				Versioning: status,
			}

			cols, _ := c.Command.Command.Flags().GetStringSlice(constants.ArgCols)
			return c.Out(table.Sprint(versioningCols, result, cols))
		},
		InitClient: false,
	})

	cmd.AddStringFlag(constants.FlagName, constants.FlagNameShort, "", "Name of the bucket", core.RequiredFlagOption())

	cmd.Command.SilenceUsage = true
	cmd.Command.Flags().SortFlags = false
	return cmd
}
