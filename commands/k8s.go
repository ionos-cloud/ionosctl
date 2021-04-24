package commands

import (
	"context"
	"errors"
	"io"
	"os"

	"github.com/fatih/structs"
	"github.com/ionos-cloud/ionosctl/pkg/builder"
	"github.com/ionos-cloud/ionosctl/pkg/config"
	"github.com/ionos-cloud/ionosctl/pkg/resources"
	"github.com/ionos-cloud/ionosctl/pkg/utils"
	"github.com/ionos-cloud/ionosctl/pkg/utils/clierror"
	"github.com/ionos-cloud/ionosctl/pkg/utils/printer"
	ionoscloud "github.com/ionos-cloud/sdk-go/v5"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func k8s() *builder.Command {
	ctx := context.TODO()
	k8sCmd := &builder.Command{
		Command: &cobra.Command{
			Use:              "k8s",
			Short:            "User Operations",
			Long:             `The sub-command of ` + "`" + `ionosctl k8s` + "`" + ` allows you to list, get, create, update, delete Users under your account. To add Users to a Group, check the ` + "`" + `ionosctl group` + "`" + ` commands. To add S3Keys to a User, check the ` + "`" + `ionosctl s3key` + "`" + ` commands.`,
			TraverseChildren: true,
		},
	}
	globalFlags := k8sCmd.Command.PersistentFlags()
	globalFlags.StringSlice(config.ArgCols, defaultUserCols, "Columns to be printed in the standard output")
	_ = viper.BindPFlag(builder.GetGlobalFlagName(k8sCmd.Command.Name(), config.ArgCols), globalFlags.Lookup(config.ArgCols))

	/*
		List Command
	*/
	builder.NewCommand(ctx, k8sCmd, noPreRun, RunK8sClustersList, "list", "List Users",
		"Use this command to get a list of existing Users available on your account.", listUserExample, true)

	/*
		Get Command
	*/
	get := builder.NewCommand(ctx, k8sCmd, PreRunK8sIdValidate, RunK8sClusterGet, "get", "Get a User",
		"Use this command to retrieve details about a specific User.\n\nRequired values to run command:\n\n* User Id", getK8sExample, true)
	get.AddStringFlag(config.ArgK8sId, "", "", config.RequiredFlagUserId)
	_ = get.Command.RegisterFlagCompletionFunc(config.ArgK8sId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getK8sClustersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})

	/*
		Create Command
	*/
	create := builder.NewCommand(ctx, k8sCmd, PreRunK8sNameEmailPwdValidate, RunK8sClusterCreate, "create", "Create a User under a particular contract",
		`Use this command to create a User under a particular contract. You need to specify the firstname, lastname, email and password for the new User.

Note: The password set here cannot be updated through the API currently. It is recommended that a new User log into the DCD and change their password.

Required values to run a command:

* User First Name
* User Last Name
* User Email
* User Password`, createUserExample, true)
	create.AddStringFlag(config.ArgUserFirstName, "", "", "The firstname for the User "+config.RequiredFlag)
	create.AddStringFlag(config.ArgUserLastName, "", "", "The lastname for the User "+config.RequiredFlag)
	create.AddStringFlag(config.ArgUserEmail, "", "", "The email for the User "+config.RequiredFlag)
	create.AddStringFlag(config.ArgUserPassword, "", "", "The password for the User (must be at least 5 characters long) "+config.RequiredFlag)
	create.AddBoolFlag(config.ArgUserAdministrator, "", false, "Assigns the User to have administrative rights")
	create.AddBoolFlag(config.ArgUserForceSecAuth, "", false, "Indicates if secure (two-factor) authentication should be forced for the User")

	/*
		Update Command
	*/
	update := builder.NewCommand(ctx, k8sCmd, noPreRun, RunK8sClusterUpdate, "update", "Update a User",
		`Use this command to update details about a specific User including their privileges.

Note: The password attribute is immutable. It is not allowed in update requests. It is recommended that the new User log into the DCD and change their password.

Required values to run command:

* User Id`, updateUserExample, true)
	update.AddStringFlag(config.ArgUserFirstName, "", "", "The firstname for the User")
	update.AddStringFlag(config.ArgUserLastName, "", "", "The lastname for the User")
	update.AddStringFlag(config.ArgUserEmail, "", "", "The email for the User")
	update.AddBoolFlag(config.ArgUserAdministrator, "", false, "Assigns the User to have administrative rights")
	update.AddBoolFlag(config.ArgUserForceSecAuth, "", false, "Indicates if secure (two-factor) authentication should be forced for the User")
	update.AddStringFlag(config.ArgK8sId, "", "", config.RequiredFlagUserId)
	_ = update.Command.RegisterFlagCompletionFunc(config.ArgK8sId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getK8sClustersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})

	/*
		Delete Command
	*/
	deleteCmd := builder.NewCommand(ctx, k8sCmd, PreRunK8sIdValidate, RunK8sClusterDelete, "delete", "Blacklists the User, disabling them",
		`This command blacklists the User, disabling them. The User is not completely purged, therefore if you anticipate needing to create a User with the same name in the future, we suggest renaming the User before you delete it.

Required values to run command:

* User Id`, deleteUserExample, true)
	deleteCmd.AddStringFlag(config.ArgK8sId, "", "", config.RequiredFlagUserId)
	_ = deleteCmd.Command.RegisterFlagCompletionFunc(config.ArgK8sId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getK8sClustersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})

	return k8sCmd
}

func PreRunK8sIdValidate(c *builder.PreCommandConfig) error {
	return builder.CheckRequiredFlags(c.ParentName, c.Name, config.ArgK8sId)
}

func PreRunK8sNameEmailPwdValidate(c *builder.PreCommandConfig) error {
	return builder.CheckRequiredFlags(c.ParentName, c.Name, config.ArgUserFirstName, config.ArgUserLastName, config.ArgUserEmail, config.ArgUserPassword)
}

func RunK8sClustersList(c *builder.CommandConfig) error {
	k8ss, _, err := c.K8s().ListClusters()
	if err != nil {
		return err
	}
	return c.Printer.Print(getK8sClusterPrint(nil, c, getK8sClusters(k8ss)))
}

func RunK8sClusterGet(c *builder.CommandConfig) error {
	u, _, err := c.K8s().GetCluster(viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgK8sId)))
	if err != nil {
		return err
	}
	return c.Printer.Print(getK8sClusterPrint(nil, c, getK8sCluster(u)))
}

func RunK8sClusterCreate(c *builder.CommandConfig) error {
	firstname := viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgUserFirstName))
	lastname := viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgUserLastName))
	email := viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgUserEmail))
	pwd := viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgUserPassword))
	secureAuth := viper.GetBool(builder.GetFlagName(c.ParentName, c.Name, config.ArgUserForceSecAuth))
	admin := viper.GetBool(builder.GetFlagName(c.ParentName, c.Name, config.ArgUserAdministrator))
	newUser := resources.K8sCluster{
		Cluster: ionoscloud.Cluster{
			Properties: &ionoscloud.UserProperties{
				Firstname:     &firstname,
				Lastname:      &lastname,
				Email:         &email,
				Administrator: &admin,
				ForceSecAuth:  &secureAuth,
				Password:      &pwd,
			},
		},
	}
	u, resp, err := c.K8s().CreateCluster(newUser)
	if err != nil {
		return err
	}
	return c.Printer.Print(getK8sClusterPrint(resp, c, getK8sCluster(u)))
}

func RunK8sClusterUpdate(c *builder.CommandConfig) error {
	oldUser, resp, err := c.K8s().GetCluster(viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgK8sId)))
	if err != nil {
		return err
	}
	newProperties := getK8sInfo(oldUser, c)
	newUser := resources.K8sCluster{
		User: ionoscloud.User{
			Properties: &newProperties.UserProperties,
		},
	}
	k8sUpd, resp, err := c.K8s().UpdateCluster(viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgK8sId)), newUser)
	if err != nil {
		return err
	}
	return c.Printer.Print(getK8sClusterPrint(resp, c, getK8sCluster(k8sUpd)))
}

func RunK8sClusterDelete(c *builder.CommandConfig) error {
	err := utils.AskForConfirm(c.Stdin, c.Printer, "delete k8s")
	if err != nil {
		return err
	}
	resp, err := c.K8s().DeleteCluster(viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgK8sId)))
	if err != nil {
		return err
	}
	return c.Printer.Print(getK8sClusterPrint(resp, c, nil))
}

func getK8sInfo(oldUser *resources.K8sCluster, c *builder.CommandConfig) *resources.K8sClusterProperties {
	var (
		firstName, lastName, email string
		forceSecureAuth, admin     bool
	)
	if properties, ok := oldUser.GetPropertiesOk(); ok && properties != nil {
		if viper.IsSet(builder.GetFlagName(c.ParentName, c.Name, config.ArgUserFirstName)) {
			firstName = viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgUserFirstName))
		} else {
			if name, ok := properties.GetFirstnameOk(); ok && name != nil {
				firstName = *name
			}
		}
		if viper.IsSet(builder.GetFlagName(c.ParentName, c.Name, config.ArgUserLastName)) {
			lastName = viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgUserLastName))
		} else {
			if name, ok := properties.GetLastnameOk(); ok && name != nil {
				lastName = *name
			}
		}
	}
	return &resources.K8sClusterProperties{
		UserProperties: ionoscloud.UserProperties{
			Firstname:     &firstName,
			Lastname:      &lastName,
			Email:         &email,
			Administrator: &admin,
			ForceSecAuth:  &forceSecureAuth,
		},
	}
}

// Output Printing

var defaultK8sClusterCols = []string{"ClusterId", "ClusterName", "K8sVersion", "AvailableUpgradeVersions", "ViableNodePoolVersions"}

type K8sClusterPrint struct {
	ClusterId                string   `json:"ClusterId,omitempty"`
	Name                     string   `json:"ClusterName,omitempty"`
	K8sVersion               string   `json:"K8sVersion,omitempty"`
	AvailableUpgradeVersions []string `json:"AvailableUpgradeVersions,omitempty"`
	ViableNodePoolVersions   []string `json:"ViableNodePoolVersions,omitempty"`
}

func getK8sClusterPrint(resp *resources.Response, c *builder.CommandConfig, k8ss []resources.K8sCluster) printer.Result {
	r := printer.Result{}
	if c != nil {
		if resp != nil {
			r.ApiResponse = resp
			r.Resource = c.ParentName
			r.Verb = c.Name
			r.WaitFlag = viper.GetBool(builder.GetFlagName(c.ParentName, c.Name, config.ArgWait))
		}
		if k8ss != nil {
			r.OutputJSON = k8ss
			r.KeyValue = getK8sClustersKVMaps(k8ss)
			r.Columns = getK8sClusterCols(builder.GetGlobalFlagName(c.ParentName, config.ArgCols), c.Printer.GetStderr())
		}
	}
	return r
}

func getK8sClusterCols(flagName string, outErr io.Writer) []string {
	if viper.IsSet(flagName) {
		var k8sCols []string
		columnsMap := map[string]string{
			"UserId":            "UserId",
			"Firstname":         "Firstname",
			"Lastname":          "Lastname",
			"Email":             "Email",
			"Administrator":     "Administrator",
			"ForceSecAuth":      "ForceSecAuth",
			"SecAuthActive":     "SecAuthActive",
			"S3CanonicalUserId": "S3CanonicalUserId",
			"Active":            "Active",
		}
		for _, k := range viper.GetStringSlice(flagName) {
			col := columnsMap[k]
			if col != "" {
				k8sCols = append(k8sCols, col)
			} else {
				clierror.CheckError(errors.New("unknown column "+k), outErr)
			}
		}
		return k8sCols
	} else {
		return defaultK8sClusterCols
	}
}

func getK8sClusters(k8ss resources.K8sClusters) []resources.K8sCluster {
	u := make([]resources.K8sCluster, 0)
	if items, ok := k8ss.GetItemsOk(); ok && items != nil {
		for _, item := range *items {
			u = append(u, resources.K8sCluster{KubernetesCluster: item})
		}
	}
	return u
}

func getK8sCluster(u *resources.K8sCluster) []resources.K8sCluster {
	k8ss := make([]resources.K8sCluster, 0)
	if u != nil {
		k8ss = append(k8ss, resources.K8sCluster{KubernetesCluster: u.KubernetesCluster})
	}
	return k8ss
}

func getK8sClustersKVMaps(us []resources.K8sCluster) []map[string]interface{} {
	out := make([]map[string]interface{}, 0, len(us))
	for _, u := range us {
		var uPrint K8sClusterPrint
		if id, ok := u.GetIdOk(); ok && id != nil {
			uPrint.ClusterId = *id
		}
		if properties, ok := u.GetPropertiesOk(); ok && properties != nil {
			if name, ok := properties.GetNameOk(); ok && name != nil {
				uPrint.Name = *name
			}
			if v, ok := properties.GetK8sVersionOk(); ok && v != nil {
				uPrint.K8sVersion = *v
			}
			if v, ok := properties.GetAvailableUpgradeVersionsOk(); ok && v != nil {
				uPrint.AvailableUpgradeVersions = *v
			}
			if v, ok := properties.GetViableNodePoolVersionsOk(); ok && v != nil {
				uPrint.ViableNodePoolVersions = *v
			}
		}
		o := structs.Map(uPrint)
		out = append(out, o)
	}
	return out
}

func getK8sClustersIds(outErr io.Writer) []string {
	err := config.Load()
	clierror.CheckError(err, outErr)
	clientSvc, err := resources.NewClientService(
		viper.GetString(config.Username),
		viper.GetString(config.Password),
		viper.GetString(config.Token),
		viper.GetString(config.ArgServerUrl),
	)
	clierror.CheckError(err, outErr)
	k8sSvc := resources.NewK8sService(clientSvc.Get(), context.TODO())
	k8ss, _, err := k8sSvc.ListClusters()
	clierror.CheckError(err, outErr)
	k8ssIds := make([]string, 0)
	if items, ok := k8ss.KubernetesClusters.GetItemsOk(); ok && items != nil {
		for _, item := range *items {
			if itemId, ok := item.GetIdOk(); ok && itemId != nil {
				k8ssIds = append(k8ssIds, *itemId)
			}
		}
	} else {
		return nil
	}
	return k8ssIds
}
