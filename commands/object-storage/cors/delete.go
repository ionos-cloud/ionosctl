package cors

import (
	"context"
	"fmt"

	"github.com/spf13/viper"

	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/pkg/confirm"
)

func DeleteCmd() *core.Command {
	cmd := core.NewCommand(context.Background(), nil, core.CommandBuilder{
		Namespace: "object-storage",
		Resource:  "cors",
		Verb:      "delete",
		Aliases:   []string{"d"},
		ShortDesc: "Delete the CORS configuration for a bucket",
		Example:   "ionosctl object-storage cors delete --name my-bucket\nionosctl object-storage cors delete --name my-bucket -f",
		PreCmdRun: func(c *core.PreCommandConfig) error {
			return core.CheckRequiredFlags(c.Command, c.NS, constants.FlagName)
		},
		CmdRun: func(c *core.CommandConfig) error {
			name := viper.GetString(core.GetFlagName(c.NS, constants.FlagName))

			if !confirm.FAsk(c.Command.Command.InOrStdin(), fmt.Sprintf("delete CORS configuration for bucket %q", name), viper.GetBool(constants.ArgForce)) {
				return fmt.Errorf(confirm.UserDenied)
			}

			s3, err := resolveRegionalClient(context.Background(), name)
			if err != nil {
				return err
			}

			_, err = s3.CORSApi.DeleteBucketCors(context.Background(), name).Execute()
			if err != nil {
				return err
			}

			fmt.Fprintf(c.Command.Command.OutOrStdout(), "CORS configuration for %q deleted successfully\n", name)
			return nil
		},
		InitClient: false,
	})

	cmd.AddStringFlag(constants.FlagName, constants.FlagNameShort, "", "Name of the bucket", core.RequiredFlagOption())

	cmd.Command.SilenceUsage = true
	cmd.Command.Flags().SortFlags = false
	return cmd
}
