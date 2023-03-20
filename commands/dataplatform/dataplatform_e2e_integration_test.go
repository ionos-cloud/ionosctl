// build +integration
package dataplatform

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/ionos-cloud/ionosctl/v6/commands/dataplatform/nodepool"
	client2 "github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/functional"
	ionoscloud "github.com/ionos-cloud/sdk-go-dataplatform"

	"github.com/cilium/fake"
	"github.com/ionos-cloud/ionosctl/v6/commands/dataplatform/cluster"
	"github.com/ionos-cloud/ionosctl/v6/pkg/constants"
	sdkcompute "github.com/ionos-cloud/sdk-go/v6"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

var (
	uniqueResourceName = "ionosctl-dataplatform-cluster-test-" + fake.AlphaNum(8)
	client             *client2.Client
	createdClusterId   string
	createdDcId        string
)

// If your test is failing because your credentials env var seem empty, try running with `godotenv -f <config-file> go test <test>`
func TestDataplatformCmd(t *testing.T) {
	var err error
	client, err = client2.Get()
	assert.NoError(t, err)
	go testClusterIdentifyRequiredNotSet(t)
	if setup() != nil {
		t.Fatalf("Failed setting up Dataplatform required resources: %s", err)
	}
	t.Cleanup(teardown)
	testClusterOk(t)
}

func testClusterOk(t *testing.T) {
	viper.Set(constants.ArgOutput, "text")
	viper.Set(constants.ArgCols, "Name")
	viper.Set(constants.ArgNoHeaders, true)
	fmt.Printf(viper.GetString(constants.ArgCols))

	c := cluster.ClusterCreateCmd()
	c.Command.Flags().Set(constants.FlagDatacenterId, createdDcId)
	c.Command.Flags().Set(constants.FlagName, uniqueResourceName)
	c.Command.Flags().Set(constants.FlagMaintenanceDay, "Friday")
	c.Command.Flags().Set(constants.FlagMaintenanceTime, "10:00:00")

	err := c.Command.Execute()
	assert.NoError(t, err)

	ls, resp, err := client.DataplatformClient.DataPlatformClusterApi.GetClusters(context.Background()).Name(uniqueResourceName).Execute()
	assert.NoError(t, err)
	assert.False(t, resp.HttpNotFound())
	items := *ls.Items
	assert.Len(t, items, 1)
	createdClusterId = *(items)[0].GetId()
	assert.Equal(t, uniqueResourceName, *(*ls.Items)[0].Properties.Name)
}

func testNodepoolOk(t *testing.T) {
	viper.Set(constants.ArgOutput, "text")
	viper.Set(constants.ArgCols, "Name")
	viper.Set(constants.ArgNoHeaders, true)
	fmt.Printf(viper.GetString(constants.ArgCols))

	c := nodepool.NodepoolCreateCmd()
	c.Command.Flags().Set(constants.FlagClusterId, createdClusterId)
	c.Command.Flags().Set(constants.FlagName, uniqueResourceName)
	c.Command.Flags().Set(constants.FlagNodeCount, "2")

	err := c.Command.Execute()
	assert.NoError(t, err)

	ls, resp, err := client.DataplatformClient.DataPlatformNodePoolApi.GetClusterNodepools(context.Background(), createdClusterId).Execute()
	assert.NoError(t, err)
	assert.False(t, resp.HttpNotFound())
	var foundNodepool ionoscloud.NodePoolResponseData
	assert.True(t, functional.Fold(*ls.GetItems(), func(found bool, x ionoscloud.NodePoolResponseData) bool {
		if *x.Properties.Name == uniqueResourceName {
			foundNodepool = x
			return true
		}
		return found
	}, false))
	assert.Equal(t, 2, foundNodepool.Properties.NodeCount)
}

func testClusterIdentifyRequiredNotSet(t *testing.T) {
	viper.Set(constants.ArgOutput, "text")
	viper.Set(constants.ArgCols, "Name")
	viper.Set(constants.ArgNoHeaders, true)
	fmt.Printf(viper.GetString(constants.ArgCols))

	c := cluster.ClusterCreateCmd()
	c.Command.Flags().Set(constants.FlagName, uniqueResourceName)
	c.Command.Flags().Set(constants.FlagMaintenanceDay, "Friday")
	c.Command.Flags().Set(constants.FlagMaintenanceTime, "10:00:00")

	err := c.Command.Execute()
	assert.ErrorContains(t, err, constants.FlagDatacenterId)
}

func setup() error {
	// make sure datacenter exists
	dcs, resp, err := client.CloudClient.DataCentersApi.DatacentersGet(context.Background()).Filter("name", uniqueResourceName).Depth(1).Execute()
	if resp.HttpNotFound() || len(*dcs.Items) < 1 {
		dc, _, err := client.CloudClient.DataCentersApi.DatacentersPost(context.Background()).Datacenter(sdkcompute.Datacenter{Properties: &sdkcompute.DatacenterProperties{Name: sdkcompute.PtrString(uniqueResourceName), Location: sdkcompute.PtrString("de/fra")}}).Execute()
		if err != nil {
			return fmt.Errorf("failed creating dc %w", err)
		}
		createdDcId = *dc.Id
		time.Sleep(10 * time.Second)
	} else if err != nil {
		return fmt.Errorf("failed getting dc %w", err)
	} else {
		createdDcId = *(*dcs.GetItems())[0].GetId()
	}
	fmt.Printf("dcId: %s\n", createdDcId)
	return nil
}

func teardown() {
	_, _, err := client.DataplatformClient.DataPlatformClusterApi.DeleteCluster(context.Background(), createdClusterId).Execute()
	if err != nil {
		fmt.Printf("failed deleting cluster: %v\n", err)
	}

	time.Sleep(5 * time.Second)

	_, err = client.CloudClient.DataCentersApi.DatacentersDelete(context.Background(), createdDcId).Execute()
	if err != nil {
		fmt.Printf("failed deleting dc: %v\n", err)
	}
}
