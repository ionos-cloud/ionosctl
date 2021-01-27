package commands

import (
	"context"
	"io"
	"os"

	"github.com/ionos-cloud/ionosctl/pkg/builder"
	"github.com/ionos-cloud/ionosctl/pkg/config"
	"github.com/ionos-cloud/ionosctl/pkg/resources"
	"github.com/ionos-cloud/ionosctl/pkg/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func datacenter() *builder.Command {
	datacenterCmd := &builder.Command{
		Command: &cobra.Command{
			Use:     "datacenter",
			Aliases: []string{"dc"},
			Short:   "Data Center operations",
		},
	}

	/*
		List Command
	*/
	list := builder.NewCommand(context.TODO(), datacenterCmd, RunDataCenterList, "list", "List Data Centers", "", true)
	list.AddStringFlag(config.ArgDataCenterId, "", "", "The unique Data Center Id")
	list.Command.RegisterFlagCompletionFunc(config.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getDataCentersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})

	/*
		Create Command
	*/
	create := builder.NewCommand(context.TODO(), datacenterCmd, RunDataCenterCreate, "create", "Create a Data Center",
		"Create a Data Center. The name, description or region can be specified.", true)
	create.AddStringFlag(config.ArgDataCenterName, "", "", "Name of the Data Center")
	create.AddStringFlag(config.ArgDataCenterDescription, "", "", "Description of the Data Center")
	create.AddStringFlag(config.ArgDataCenterRegion, "", "de/txl", "Location for the Data Center")

	/*
		Update Command
	*/
	update := builder.NewCommand(context.TODO(), datacenterCmd, RunDataCenterUpdate, "update", "Update a Data Center",
		"Update a Data Center. Data Center Id is required", true)
	update.AddStringFlag(config.ArgDataCenterId, "", "", "The unique Data Center Id [Required flag]")
	update.Command.RegisterFlagCompletionFunc(config.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getDataCentersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	update.AddStringFlag(config.ArgDataCenterName, "", "", "Name of the Data Center")
	update.AddStringFlag(config.ArgDataCenterDescription, "", "", "Description of the Data Center")

	/*
		Delete Command
	*/
	delete := builder.NewCommand(context.TODO(), datacenterCmd, RunDataCenterDelete, "delete", "Delete a Data Center",
		"Delete a Data Center. Data Center Id is required.", true)
	delete.AddStringFlag(config.ArgDataCenterId, "", "", "The unique Data Center Id [Required flag]")
	delete.Command.RegisterFlagCompletionFunc(config.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getDataCentersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})

	return datacenterCmd
}

func RunDataCenterList(c *builder.CommandConfig) error {
	if viper.GetString(builder.GetFlagName(c.Name, config.ArgDataCenterId)) == "" {
		datacenters, _, err := c.DataCenters().List()
		if err != nil {
			return err
		}
		dcs := getDataCenters(datacenters)
		c.Printer.Result(&utils.SuccessResult{
			OutputJSON: datacenters,
			KeyValue:   getDataCentersKVMaps(dcs),
			Columns:    getDataCenterCols(),
		})
		return nil
	} else {
		datacenter, _, err := c.DataCenters().Get(viper.GetString(builder.GetFlagName(c.Name, config.ArgDataCenterId)))
		if err != nil {
			return err
		}

		c.Printer.Result(&utils.SuccessResult{
			KeyValue:   getDataCentersKVMaps([]resources.Datacenter{*datacenter}),
			Columns:    getDataCenterCols(),
			OutputJSON: datacenter,
		})

		return nil
	}
}

func RunDataCenterCreate(c *builder.CommandConfig) error {
	name := viper.GetString(builder.GetFlagName(c.Name, config.ArgDataCenterName))
	description := viper.GetString(builder.GetFlagName(c.Name, config.ArgDataCenterDescription))
	region := viper.GetString(builder.GetFlagName(c.Name, config.ArgDataCenterRegion))
	dc, resp, err := c.DataCenters().Create(name, description, region)
	if err != nil {
		return err
	}

	c.Printer.Result(&utils.SuccessResult{
		KeyValue:    getDataCentersKVMaps([]resources.Datacenter{*dc}),
		Columns:     getDataCenterCols(),
		OutputJSON:  dc,
		ApiResponse: resp,
		Resource:    "datacenter",
		Verb:        "create",
	})

	return nil
}

func RunDataCenterUpdate(c *builder.CommandConfig) error {
	if viper.GetString(builder.GetFlagName(c.Name, config.ArgDataCenterId)) == "" {
		return utils.NewRequiredFlagErr(config.ArgDataCenterId)
	}

	dc, resp, err := c.DataCenters().Update(
		viper.GetString(builder.GetFlagName(c.Name, config.ArgDataCenterId)),
		viper.GetString(builder.GetFlagName(c.Name, config.ArgDataCenterName)),
		viper.GetString(builder.GetFlagName(c.Name, config.ArgDataCenterDescription)),
	)
	if err != nil {
		return err
	}

	c.Printer.Result(&utils.SuccessResult{
		KeyValue:    getDataCentersKVMaps([]resources.Datacenter{*dc}),
		Columns:     getDataCenterCols(),
		OutputJSON:  dc,
		ApiResponse: resp,
		Resource:    "datacenter",
		Verb:        "update",
	})

	return nil
}

func RunDataCenterDelete(c *builder.CommandConfig) error {
	if viper.GetString(builder.GetFlagName(c.Name, config.ArgDataCenterId)) == "" {
		return utils.NewRequiredFlagErr(config.ArgDataCenterId)
	}
	err := utils.AskForConfirm(c.Printer.Stdin, c.Printer.Stdout, "delete data center")
	if err != nil {
		return err
	}
	resp, err := c.DataCenters().Delete(viper.GetString(builder.GetFlagName(c.Name, config.ArgDataCenterId)))
	if err != nil {
		return err
	}

	c.Printer.Result(&utils.SuccessResult{
		ApiResponse: resp,
		Resource:    "datacenter",
		Verb:        "delete",
	})

	return nil
}

func getDataCenterCols() []string {
	return []string{"ID", "Name", "Location", "Description"}
}

func getDataCenters(datacenters resources.Datacenters) []resources.Datacenter {
	dc := make([]resources.Datacenter, 0)
	for _, d := range *datacenters.Items {
		dc = append(dc, resources.Datacenter{d})
	}
	return dc
}

func getDataCentersKVMaps(dcs []resources.Datacenter) []map[string]interface{} {
	out := make([]map[string]interface{}, 0, len(dcs))
	for _, dc := range dcs {
		properties := dc.GetProperties()
		o := map[string]interface{}{
			"ID":       *dc.GetId(),
			"Name":     *properties.GetName(),
			"Location": *properties.GetLocation(),
		}
		if description, ok := properties.GetDescriptionOk(); ok {
			o["Description"] = *description
		}
		out = append(out, o)
	}
	return out
}

func getDataCentersIds(outErr io.Writer) []string {
	err := config.LoadFile()
	utils.CheckError(err, outErr)

	clientSvc, err := resources.NewClientService(
		viper.GetString(config.Username),
		viper.GetString(config.Password),
		viper.GetString(config.ArgServerUrl),
	)
	utils.CheckError(err, outErr)

	datacenterSvc := resources.NewDataCenterService(clientSvc.Get(), context.TODO())
	datacenters, _, err := datacenterSvc.List()
	utils.CheckError(err, outErr)

	dcIds := make([]string, 0)
	if datacenters.Datacenters.Items != nil {
		for _, d := range *datacenters.Datacenters.Items {
			dcIds = append(dcIds, *d.GetId())
		}
	}
	return dcIds
}
