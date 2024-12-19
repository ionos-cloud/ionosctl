package user

import (
	"context"
	"encoding/csv"
	"fmt"
	"strings"

	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/json2table/resource2table"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/jsontabwriter"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/tabheaders"
	"github.com/ionos-cloud/sdk-go-bundle/products/dbaas/mongo/v2"

	"github.com/ionos-cloud/ionosctl/v6/commands/dbaas/mongo/completer"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	sdkgo "github.com/ionos-cloud/sdk-go-bundle/products/dbaas/mongo/v2"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func UserCreateCmd() *core.Command {
	cmd := core.NewCommand(context.TODO(), nil, core.CommandBuilder{
		Namespace: "dbaas-mongo",
		Resource:  "user",
		Verb:      "create",
		Aliases:   []string{"c"},
		ShortDesc: "Create MongoDB users.",
		Example:   "ionosctl dbaas mongo user create --cluster-id CLUSTER_ID --name USERNAME --password PASSWORD --roles DATABASE=ROLE",
		PreCmdRun: func(c *core.PreCommandConfig) (err error) {
			err = c.Command.Command.MarkFlagRequired(constants.FlagClusterId)
			if err != nil {
				return err
			}
			err = c.Command.Command.MarkFlagRequired(constants.FlagName)
			if err != nil {
				return err
			}
			err = c.Command.Command.MarkFlagRequired(constants.ArgPassword)
			if err != nil {
				return err
			}
			err = c.Command.Command.MarkFlagRequired(FlagRoles)
			if err != nil {
				return fmt.Errorf("%w: DB1=Role1,DB2=Role2. Roles: read, readWrite, readAnyDatabase, readWriteAnyDatabase, dbAdmin, dbAdminAnyDatabase, clusterMonitor", err)
			}
			return nil
		},
		CmdRun: func(c *core.CommandConfig) error {
			clusterId := viper.GetString(core.GetFlagName(c.NS, constants.FlagClusterId))
			fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("Creating user for cluster %s", clusterId))

			input := mongo.UserProperties{}
			if fn := core.GetFlagName(c.NS, FlagRoles); viper.IsSet(fn) {
				roles, err := parseRoles(viper.GetString(fn))
				if err != nil {
					return err
				}

				input.Roles = roles
			}

			if fn := core.GetFlagName(c.NS, constants.FlagName); viper.IsSet(fn) {
				input.Username = viper.GetString(fn)
			}

			if fn := core.GetFlagName(c.NS, constants.ArgPassword); viper.IsSet(fn) {
				input.Password = viper.GetString(fn)
			}

			u, _, err := client.Must().MongoClient.UsersApi.ClustersUsersPost(context.Background(), clusterId).
				User(sdkgo.User{Properties: &input}).Execute()
			if err != nil {
				return err
			}

			cols, _ := c.Command.Command.Flags().GetStringSlice(constants.ArgCols)

			uConverted, err := resource2table.ConvertDbaasMongoUserToTable(u)
			if err != nil {
				return err
			}

			out, err := jsontabwriter.GenerateOutputPreconverted(u, uConverted, tabheaders.GetHeadersAllDefault(allCols, cols))
			if err != nil {
				return err
			}

			fmt.Fprintf(c.Command.Command.OutOrStdout(), out)
			return nil
		},
		InitClient: true,
	})

	cmd.AddStringFlag(constants.FlagClusterId, constants.FlagIdShort, "", "")
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.FlagClusterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.MongoClusterIds(), cobra.ShellCompDirectiveNoFileComp
	})

	// required Path flags
	cmd.AddStringFlag(constants.FlagName, constants.FlagNameShort, "", "The authentication username", core.RequiredFlagOption())
	cmd.AddStringFlag(constants.ArgPassword, constants.ArgPasswordShort, "", "The authentication password", core.RequiredFlagOption())

	cmd.AddStringFlag(FlagRoles, FlagRolesShort, "", "User's role for each db. DB1=Role1,DB2=Role2."+
		"Roles: read, readWrite, readAnyDatabase, readWriteAnyDatabase, dbAdmin, dbAdminAnyDatabase, clusterMonitor, enableSharding",
		core.RequiredFlagOption(),
		core.WithCompletionComplex(
			func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
				if toComplete[len(toComplete)-1] != '=' {
					toComplete += "="
				}
				return []string{
					toComplete + "read",
					toComplete + "readWrite",
					toComplete + "readAnyDatabase",
					toComplete + "readWriteAnyDatabase",
					toComplete + "dbAdmin",
					toComplete + "dbAdminAnyDatabase",
					toComplete + "clusterMonitor",
					toComplete + "enableSharding",
				}, cobra.ShellCompDirectiveNoFileComp
			}, "",
		),
	)

	cmd.Command.SilenceUsage = true

	return cmd
}

// parseRoles converts string "DB1=role,DB2=other-role" input into a slice of "UserRoles" objects
func parseRoles(val string) ([]sdkgo.UserRoles, error) {
	// step 1. Use identical slice conversion as pflag lib for compatibility reasons
	stringReader := strings.NewReader(val)
	csvReader := csv.NewReader(stringReader)
	tuples, err := csvReader.Read()
	if err != nil {
		return nil, fmt.Errorf("failed to read csv: %w", err)
	}

	// step 2. For each tuple, get its Database and Role
	var rs []sdkgo.UserRoles
	for _, t := range tuples {
		dbAndRole := strings.Split(t, "=")
		if len(dbAndRole) != 2 {
			return nil, fmt.Errorf("invalid input format: %s, need db=role1,db=role2", val)
		}
		r := sdkgo.UserRoles{Database: &dbAndRole[0], Role: &dbAndRole[1]}
		rs = append(rs, r)
	}
	return rs, nil
}
