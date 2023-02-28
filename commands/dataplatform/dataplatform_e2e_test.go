package dataplatform

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/cilium/fake"
	"github.com/ionos-cloud/ionosctl/commands/dbaas/mongo/cluster"
	"github.com/ionos-cloud/ionosctl/commands/dbaas/mongo/user"
	"github.com/ionos-cloud/ionosctl/pkg/config"
	"github.com/ionos-cloud/ionosctl/pkg/constants"
	sdkcompute "github.com/ionos-cloud/sdk-go/v6"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

var (
	uniqueResourceName = "ionosctl-dataplatform-cluster-test-" + fake.AlphaNum(8)
	client             *config.Client
	createdClusterId   string
	createdDcId        string
)

// If your test is failing because your credentials env var seem empty, try running with `godotenv -f <config-file> go test <test>`
func TestDataplatformCmd(t *testing.T) {
	var err error
	client, err = config.GetClient()
	assert.NoError(t, err)
	if setup() != nil {
		t.Fatalf("Failed setting up Dataplatform required resources: %s", err)
	}
	t.Cleanup(teardownTestMongoCommands)
	testMongoClusterCreate(t, dcId, lanId)
	testMongoUserCreate(t)
}

func testMongoUserCreate(t *testing.T) {
	viper.Set(constants.ArgOutput, "text")
	viper.Set(constants.ArgCols, "Name")
	viper.Set(constants.ArgNoHeaders, true)
	fmt.Printf(viper.GetString(constants.ArgCols))

	name := fake.Name()
	c := user.UserCreateCmd()
	c.Command.Flags().Set(constants.FlagClusterId, createdClusterId)
	c.Command.Flags().Set(user.FlagRoles, fake.DeploymentTier()+"="+"readWrite,"+fake.DeploymentTier()+"="+"readWrite")
	c.Command.Flags().Set(constants.FlagName, name)
	c.Command.Flags().Set(constants.ArgPassword, fake.AlphaNum(12))

	err := c.Command.Execute()
	assert.NoError(t, err)

	createdUsers, _, err := client.MongoClient.UsersApi.ClustersUsersGet(context.Background(), createdClusterId).Execute()
	assert.NoError(t, err)
	assert.Equal(t, name, *(*createdUsers.GetItems())[0].GetProperties().Username)
}

func testMongoClusterCreate(t *testing.T, dcId, lanId string) {
	viper.Set(constants.ArgOutput, "text")
	viper.Set(constants.ArgCols, "Name")
	viper.Set(constants.ArgNoHeaders, true)
	fmt.Printf(viper.GetString(constants.ArgCols))

	c := cluster.ClusterCreateCmd()
	c.Command.Flags().Set(constants.FlagDatacenterId, dcId)
	c.Command.Flags().Set(constants.FlagLanId, lanId)
	c.Command.Flags().Set(constants.FlagName, uniqueResourceName)
	c.Command.Flags().Set(constants.FlagInstances, "1")
	c.Command.Flags().Set(constants.FlagTemplateId, getPlaygroundTemplateUuid())
	c.Command.Flags().Set(constants.FlagMaintenanceDay, "Friday")
	c.Command.Flags().Set(constants.FlagMaintenanceTime, "10:00:00")
	c.Command.Flags().Set(constants.FlagCidr, cidr)

	err := c.Command.Execute()
	assert.NoError(t, err)

	createdCluster, _, err := client.MongoClient.ClustersApi.ClustersGet(context.Background()).FilterName(uniqueResourceName).Execute()
	createdClusterId = *(*createdCluster.GetItems())[0].GetId()
	assert.NoError(t, err)
	assert.Equal(t, uniqueResourceName, *(*createdCluster.Items)[0].Properties.DisplayName)
	assert.Equal(t, "de/fra", *(*createdCluster.Items)[0].Properties.Location)
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
		client.DataplatformClient.GetConfig().
		return fmt.Errorf("failed getting dc %w", err)
	} else {
		createdDcId = *(*dcs.GetItems())[0].GetId()
	}

	fmt.Printf("dcId: %s\n", createdDcId)
	return nil
}

func teardown() {
	_, _, err := client.MongoClient.ClustersApi.ClustersDelete(context.Background(), createdClusterId).Execute()
	if err != nil {
		fmt.Printf("failed deleting cluster: %v\n", err)
	}

	_, err = client.CloudClient.DataCentersApi.DatacentersDelete(context.Background(), createdDcId).Execute()
	if err != nil {
		fmt.Printf("failed deleting dc: %v\n", err)
	}
}
