//go:build integration
// +build integration

package dataplatform

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/ionos-cloud/ionosctl/v6/commands/dataplatform/nodepool"
	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/pkg/functional"
	"github.com/ionos-cloud/sdk-go-bundle/products/dataplatform/v2"

	"github.com/cilium/fake"
	"github.com/ionos-cloud/ionosctl/v6/commands/dataplatform/cluster"
	sdkcompute "github.com/ionos-cloud/sdk-go/v6"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

var (
	uniqueResourceName = "ionosctl-dataplatform-cluster-test-" + fake.AlphaNum(8)
	createdDcId        string
	createdClusterId   string
	createdNodepoolId  string
)

// If your test is failing because your credentials env var seem empty, try running with `godotenv -f <config-file> go test <test>`
func TestDataplatformCmd(t *testing.T) {
	testClusterIdentifyRequiredNotSet(t)
	if err := setup(); err != nil {
		t.Fatalf("Failed setting up Dataplatform required resources: %s", err)
	}
	t.Cleanup(teardown)
	testClusterOk(t)
	testNodepoolOk(t)
}

func testClusterOk(t *testing.T) {
	viper.Set(constants.FlagOutput, "text")
	viper.Set(constants.FlagCols, "Name")
	viper.Set(constants.FlagNoHeaders, true)

	c := cluster.ClusterCreateCmd()
	c.Command.Flags().Set(constants.FlagDatacenterId, createdDcId)
	c.Command.Flags().Set(constants.FlagName, uniqueResourceName)
	c.Command.Flags().Set(constants.FlagMaintenanceDay, "Friday")
	c.Command.Flags().Set(constants.FlagMaintenanceTime, "10:00:00")

	err := c.Command.Execute()
	if err != nil {
		if got403 := assert.ErrorContains(t, err, "403 Forbidden"); got403 {
			dc, _, errdc := client.Must().CloudClient.DataCentersApi.DatacentersFindById(context.Background(), createdDcId).Execute()
			if errdc != nil {
				t.Fatalf(fmt.Errorf("debugging 403 Forbidden went wrong: failed finding datacenter by ID %s: %w", createdDcId, errdc).Error())
			}
			dcAsJson, errJson := json.MarshalIndent(dc, "", "  ")
			assert.NoError(t, errJson)
			t.Fatalf(fmt.Sprintf("Received a '403 Forbidden' error:\n%s\n while attempting to create a dataplatform cluster in the datacenter:\n%s ", err.Error(), dcAsJson))
		}

		t.Fatalf(fmt.Errorf("failed creating a dataplatform cluster: %w", err).Error())
	}

	ls, resp, err := client.Must().DataplatformClient.DataPlatformClusterApi.ClustersGet(context.Background()).Name(uniqueResourceName).Execute()
	assert.NoError(t, err, fmt.Errorf("failed verifying created cluster via SDK: %w", err).Error())
	assert.False(t, resp.HttpNotFound())
	items := *ls.Items
	assert.Len(t, items, 1)
	createdClusterId = *(items)[0].GetId()
	assert.Equal(t, uniqueResourceName, *(*ls.Items)[0].Properties.Name)
}

func testNodepoolOk(t *testing.T) {
	viper.Set(constants.FlagOutput, "text")
	viper.Set(constants.FlagCols, "Name")
	viper.Set(constants.FlagNoHeaders, true)

	c := nodepool.NodepoolCreateCmd()
	c.Command.Flags().Set(constants.FlagClusterId, createdClusterId)
	c.Command.Flags().Set(constants.FlagName, uniqueResourceName)
	c.Command.Flags().Set(constants.FlagNodeCount, "1")

	err := c.Command.Execute()
	assert.NoError(t, err)

	ls, resp, err := client.Must().DataplatformClient.DataPlatformNodePoolApi.ClustersNodepoolsGet(context.Background(), createdClusterId).Execute()
	assert.NoError(t, err)
	assert.False(t, resp.HttpNotFound())
	var foundNodepool dataplatform.NodePoolResponseData
	// Filter by name, as API doesn't support this :(
	assert.True(t,
		functional.Fold(ls.GetItems(), func(found bool, x dataplatform.NodePoolResponseData) bool {
			if *x.Properties.Name == uniqueResourceName {
				foundNodepool = x
				return true
			}
			return found
		}, false),
		fmt.Sprintf("Couldn't filter the dataplatform nodepool by name (%s) that was supposed to be created by the tested command", uniqueResourceName),
	)
	assert.Equal(t, int32(1), *(foundNodepool.Properties.NodeCount))
	createdNodepoolId = *foundNodepool.Id

}

func testClusterIdentifyRequiredNotSet(t *testing.T) {
	viper.Set(constants.FlagOutput, "text")
	viper.Set(constants.FlagCols, "Name")
	viper.Set(constants.FlagNoHeaders, true)

	c := cluster.ClusterCreateCmd()
	c.Command.Flags().Set(constants.FlagName, uniqueResourceName)
	c.Command.Flags().Set(constants.FlagMaintenanceDay, "Friday")
	c.Command.Flags().Set(constants.FlagMaintenanceTime, "10:00:00")

	err := c.Command.Execute()
	assert.ErrorContains(t, err, constants.FlagDatacenterId)
}

func setup() error {
	// make sure datacenter exists
	dcs, resp, err := client.Must().CloudClient.DataCentersApi.DatacentersGet(context.Background()).Filter("name", uniqueResourceName).Depth(1).Execute()
	if resp.HttpNotFound() || len(*dcs.Items) < 1 {
		dc, _, err := client.Must().CloudClient.DataCentersApi.DatacentersPost(context.Background()).Datacenter(sdkcompute.Datacenter{Properties: &sdkcompute.DatacenterProperties{Name: sdkcompute.PtrString(uniqueResourceName), Location: sdkcompute.PtrString("fr/par")}}).Execute()
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
	time.Sleep(15 * time.Second)
	return nil
}

func teardown() {
	_, _, err := client.Must().DataplatformClient.DataPlatformNodePoolApi.ClustersNodepoolsDelete(context.Background(),
		createdClusterId, createdNodepoolId).Execute()
	if err != nil {
		fmt.Printf("failed deleting nodepool: %v\n", err)
	}

	time.Sleep(30 * time.Second)

	_, _, err = client.Must().DataplatformClient.DataPlatformClusterApi.ClustersDelete(context.Background(), createdClusterId).Execute()
	if err != nil {
		fmt.Printf("failed deleting cluster: %v\n", err)
	}

	time.Sleep(30 * time.Second)

	_, err = client.Must().CloudClient.DataCentersApi.DatacentersDelete(context.Background(), createdDcId).Execute()
	if err != nil {
		fmt.Printf("failed deleting dc: %v\n", err)
	}
}
