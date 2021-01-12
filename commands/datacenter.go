package commands

import (
	"context"
	"errors"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/ionos-cloud/ionosctl/pkg/config"
	"github.com/ionos-cloud/ionosctl/pkg/helpers"
	"github.com/ionos-cloud/ionosctl/pkg/resources"
	ionoscloud "github.com/ionos-cloud/sdk-go/v5"
	"github.com/spf13/cobra"
)

func listDataCenter() *cobra.Command {
	listDataCenterCmd := &cobra.Command{
		Use:     "datacenter",
		Aliases: []string{"dc"},
		Short:   "List command for Data Center",
		RunE: func(cmd *cobra.Command, args []string) error {
			err := runListDataCenter()
			return err
		},
	}
	return listDataCenterCmd
}

func createDataCenter() *cobra.Command {
	var datacenter resources.DataCenter

	postCmd := &cobra.Command{
		Use:     "datacenter",
		Aliases: []string{"dc"},
		Short:   "Create command for Data Center",
		RunE: func(cmd *cobra.Command, args []string) error {
			err := runCreateDataCenter(
				datacenter.GetName(),
				datacenter.GetDescription(),
				datacenter.GetRegion())
			return err
		},
	}

	flags := postCmd.Flags()
	flags.StringVarP(&datacenter.Name, "name", "n", "", "Name of the Data Center")
	flags.StringVarP(&datacenter.Description, "description", "d", "", "Description of the Data Center")
	flags.StringVarP(&datacenter.Region, "region", "l", "de/txl", "Location of the Data Center")

	return postCmd
}

func updateDataCenter() *cobra.Command {
	var datacenter resources.DataCenter

	updateCmd := &cobra.Command{
		Use:     "datacenter",
		Aliases: []string{"dc"},
		Short:   "Update command for Data Center",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if datacenter.GetId() == "" {
				return errors.New("no data center id provided")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			err := runUpdateDataCenter(
				datacenter.GetId(),
				datacenter.GetName(),
				datacenter.GetDescription())
			return err
		},
	}

	flags := updateCmd.Flags()
	flags.StringVarP(&datacenter.Id, "id", "i", "", "The unique ID of the Data Center")
	flags.StringVarP(&datacenter.Description, "description", "d", "", "Description of the Data Center")
	flags.StringVarP(&datacenter.Name, "name", "n", "", "Name of the Data Center")

	return updateCmd
}

func deleteDataCenter() *cobra.Command {
	var (
		datacenter resources.DataCenter
		deleteAll  bool
	)
	deleteCmd := &cobra.Command{
		Use:     "datacenter",
		Aliases: []string{"dc"},
		Short:   "Delete command for Data Center",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if datacenter.GetId() == "" && !deleteAll {
				return errors.New("no data center id provided")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			err := runDeleteDataCenter(datacenter.GetId(), deleteAll)
			return err
		},
	}

	flags := deleteCmd.Flags()
	flags.StringVarP(&datacenter.Id, "id", "i", "", "The unique ID of the Data Center")
	flags.BoolVar(&deleteAll, "all", false, "If set, it deletes all Data Centers")

	return deleteCmd
}

func runListDataCenter() error {
	apiClient, err := config.GetAPIClient()
	if err != nil {
		return err
	}

	req := apiClient.DataCenterApi.DatacentersGet(context.Background())
	datacenters, _, err := apiClient.DataCenterApi.DatacentersGetExecute(req)
	if err != nil {
		return err
	}

	if len(*datacenters.Items) != 0 {
		w := tabwriter.NewWriter(os.Stdout, 5, 0, 3, ' ', tabwriter.Debug)
		_, err = fmt.Fprintln(w, "Id\t Name\t Region\t")
		if err != nil {
			return err
		}

		for _, item := range *datacenters.Items {
			properties := item.GetProperties()

			datacenterName := *properties.GetName()
			datacenterRegion := *properties.GetLocation()
			datacenterId := *item.GetId()

			_, err = fmt.Fprintln(w,
				datacenterId, "\t",
				datacenterName, "\t",
				datacenterRegion, "\t")
			if err != nil {
				return err
			}
		}
		err = w.Flush()
		if err != nil {
			return err
		}
	} else {
		fmt.Println("no data centers found")
	}
	return nil
}

func runCreateDataCenter(name, description, region string) error {
	apiClient, err := config.GetAPIClient()
	if err != nil {
		return err
	}

	req := apiClient.DataCenterApi.DatacentersPost(context.Background())
	datacenter := ionoscloud.Datacenter{
		Properties: &ionoscloud.DatacenterProperties{
			Name:        &name,
			Description: &description,
			Location:    &region,
		},
	}
	req = req.Datacenter(datacenter)
	_, _, err = apiClient.DataCenterApi.DatacentersPostExecute(req)
	return err
}

func runUpdateDataCenter(datacenterId, name, description string) error {
	apiClient, err := config.GetAPIClient()
	if err != nil {
		return err
	}

	properties := ionoscloud.DatacenterProperties{}
	if name != "" {
		properties.SetName(name)
	}
	if description != "" {
		properties.SetDescription(description)
	}

	req := apiClient.DataCenterApi.DatacentersPatch(context.Background(), datacenterId)
	req = req.Datacenter(properties)
	_, _, err = apiClient.DataCenterApi.DatacentersPatchExecute(req)
	return err
}

func runDeleteDataCenter(datacenterId string, deleteAll bool) error {
	apiClient, err := config.GetAPIClient()
	if err != nil {
		return err
	}
	if datacenterId != "" {
		err = helpers.AskForConfirm("delete data center")
		if err != nil {
			return err
		}
		err = deleteDataCenterById(apiClient, datacenterId)
		if err != nil {
			return err
		}
	}

	if deleteAll {
		req := apiClient.DataCenterApi.DatacentersGet(context.Background())
		datacenters, _, err := apiClient.DataCenterApi.DatacentersGetExecute(req)
		if err != nil {
			return err
		}
		err = helpers.AskForConfirm("delete all data centers")
		if err != nil {
			return err
		}
		if len(*datacenters.Items) != 0 {
			for _, item := range *datacenters.Items {
				datacenterId := *item.GetId()
				err := deleteDataCenterById(apiClient, datacenterId)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func deleteDataCenterById(apiClient *ionoscloud.APIClient, datacenterId string) error {
	req := apiClient.DataCenterApi.DatacentersDelete(context.Background(), datacenterId)

	_, _, err := apiClient.DataCenterApi.DatacentersDeleteExecute(req)
	return err
}
