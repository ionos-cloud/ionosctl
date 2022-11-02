package cluster

import (
	"context"
	"github.com/ionos-cloud/ionosctl/pkg/constants"
	"github.com/ionos-cloud/ionosctl/pkg/core"
	"github.com/spf13/viper"
)

func ClusterDeleteAllCmd() *core.Command {
	cmd := core.NewCommand(context.TODO(), nil /* circular dependency 🤡*/, core.CommandBuilder{
		Namespace: "dbaas-mongo",
		Resource:  "cluster",
		Verb:      "delete-all",
		ShortDesc: "Delete all Mongo Clusters. Highly destructive!",
		Example:   "ionosctl dbaas mongo cluster delete-all --force",
		PreCmdRun: func(c *core.PreCommandConfig) error {
			return c.Command.Command.MarkFlagRequired(constants.ArgForce)
		},
		CmdRun: func(c *core.CommandConfig) error {
			c.Printer.Verbose("Deleting All Clusters!")
			_, err := c.DbaasMongoServices.Clusters().DeleteAll(viper.GetString(core.GetFlagName(c.NS, constants.FlagName)))
			return err
		},
		InitClient: true,
	})

	cmd.AddBoolFlag(constants.ArgForce, constants.ArgForceShort, false, "WARNING: HIGHLY DESTRUCTIVE! If this flag is set, you are ready to 'pull the trigger'! This flag confirms that you are sure you want to delete all clusters")
	cmd.AddStringFlag(constants.FlagName, constants.FlagNameShort, "", "Filter clusters to be deleted by name")

	cmd.Command.SilenceUsage = true

	return cmd
}
