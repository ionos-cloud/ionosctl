package user

import (
	"context"
	"encoding/csv"
	"fmt"
	"github.com/cjrd/allocate"
	"github.com/mmatczuk/anyflag"
	"strings"

	"github.com/ionos-cloud/ionosctl/commands/dbaas/mongo/completer"
	"github.com/ionos-cloud/ionosctl/pkg/constants"
	"github.com/ionos-cloud/ionosctl/pkg/core"
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
		Example:   "ionosctl dbaas mongo user create --cluster-id CLUSTER_ID --name USERNAME --password PASSWORD --database DATABASE",
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
			c.Printer.Verbose("Creating users for cluster %s", clusterId)

			userProperties.Roles = &roles
			u, _, err := c.DbaasMongoServices.Users().Create(clusterId, sdkgo.User{Properties: &userProperties})
			if err != nil {
				return err
			}
			return c.Printer.Print(getUserPrint(c, &[]sdkgo.User{u}))
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
			return nil, fmt.Errorf("invalid input format: %s, need db=role\n", val)
		}
		r := sdkgo.UserRoles{Database: &dbAndRole[0], Role: &dbAndRole[1]}
		rs = append(rs, r)
	}
	return rs, nil
}

/*
    User:
      description: MongoDB database user.
      properties:
        type:
          $ref: '#/components/schemas/ResourceType'
        metadata:
          $ref: '#/components/schemas/UserMetadata'
        properties:
          $ref: '#/components/schemas/UserProperties'
      type: object

   UserProperties:
     description: Mongodb user properties.
     required:
       - username
       - password
       - database
     properties:
       username:
         type: string
       password:
         type: string
         writeOnly: true
         minLength: 10
       roles:
         type: array
         items:
           $ref: '#/components/schemas/UserRoles'

    UserRoles:
      description: a list of mongodb user role.
      properties:
        role:
          type: string
          enum:
          - read
          - readWrite
          - readAnyDatabase
          example: read
        database:
          type: string
      type: object
*/
