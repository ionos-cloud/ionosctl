package user

import (
	"context"
	"encoding/csv"
	"fmt"
	"strings"

	"github.com/cjrd/allocate"
	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/pkg/jsontabwriter"
	"github.com/ionos-cloud/ionosctl/v6/pkg/printer"
	"github.com/mmatczuk/anyflag"

	"github.com/ionos-cloud/ionosctl/v6/commands/dbaas/mongo/completer"
	"github.com/ionos-cloud/ionosctl/v6/pkg/constants"
	"github.com/ionos-cloud/ionosctl/v6/pkg/core"
	sdkgo "github.com/ionos-cloud/sdk-go-dbaas-mongo"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func UserCreateCmd() *core.Command {
	var userProperties = sdkgo.UserProperties{}
	var roles []sdkgo.UserRoles

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
			fmt.Fprintf(c.Stderr, jsontabwriter.GenerateVerboseOutput("Creating users for cluster %s", clusterId))

			userProperties.Roles = &roles
			u, _, err := client.Must().MongoClient.UsersApi.ClustersUsersPost(context.Background(), clusterId).
				User(sdkgo.User{Properties: &userProperties}).Execute()
			if err != nil {
				return err
			}

			cols, _ := c.Command.Command.Flags().GetStringSlice(constants.ArgCols)

			uConverted, err := convertUserToTable(u)
			if err != nil {
				return err
			}

			out, err := jsontabwriter.GenerateOutputPreconverted(u, uConverted, printer.GetHeadersAllDefault(allCols, cols))
			if err != nil {
				return err
			}

			fmt.Fprintf(c.Stdout, out)
			return nil
		},
		InitClient: true,
	})

	cmd.AddStringFlag(constants.FlagClusterId, constants.FlagIdShort, "", "")
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.FlagClusterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.MongoClusterIds(), cobra.ShellCompDirectiveNoFileComp
	})

	// required Path flags
	_ = allocate.Zero(&userProperties)
	cmd.AddStringVarFlag(userProperties.Username, constants.FlagName, constants.FlagNameShort, "", "The authentication username", core.RequiredFlagOption())
	cmd.AddStringVarFlag(userProperties.Password, constants.ArgPassword, constants.ArgPasswordShort, "", "The authentication password", core.RequiredFlagOption())

	sliceOfRolesFlag := anyflag.NewValue(nil, &roles, rolesParser)
	cmd.Command.Flags().VarP(sliceOfRolesFlag, FlagRoles, FlagRolesShort, "User's role for each db. DB1=Role1,DB2=Role2. Roles: read, readWrite, readAnyDatabase, readWriteAnyDatabase, dbAdmin, dbAdminAnyDatabase, clusterMonitor")
	_ = viper.BindPFlag(core.GetFlagName(cmd.NS, FlagRoles), cmd.Command.Flags().Lookup(FlagRoles))
	_ = cmd.Command.RegisterFlagCompletionFunc(FlagRoles, func(c *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
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
	})

	cmd.Command.SilenceUsage = true

	return cmd
}

// rolesParser converts user input into a slice of "UserRoles" objects
func rolesParser(val string) ([]sdkgo.UserRoles, error) {
	// step 1. Use identical slice conversion as pflag lib for compatibility reasons
	stringReader := strings.NewReader(val)
	csvReader := csv.NewReader(stringReader)
	tuples, err := csvReader.Read()
	if err != nil {
		return nil, err
	}

	// step 2. For each tuple, get its Database and Role
	var rs []sdkgo.UserRoles
	for _, t := range tuples {
		dbAndRole := strings.Split(t, "=")
		if len(dbAndRole) != 2 {
			return nil, fmt.Errorf("invalid input format: %s, need db=role1,db=role2\n", val)
		}
		r := sdkgo.UserRoles{Database: &dbAndRole[0], Role: &dbAndRole[1]}
		rs = append(rs, r)
	}
	return rs, nil
}
