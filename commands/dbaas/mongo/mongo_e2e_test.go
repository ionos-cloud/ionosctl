package mongo

import (
	"context"
	"fmt"
	"github.com/cilium/fake"
	"github.com/ionos-cloud/ionosctl/commands/dbaas/mongo/cluster"
	"github.com/ionos-cloud/ionosctl/commands/dbaas/mongo/user"
	"github.com/ionos-cloud/ionosctl/pkg/config"
	"github.com/ionos-cloud/ionosctl/pkg/constants"
	sdkcompute "github.com/ionos-cloud/sdk-go/v6"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

var (
	uniqueResourceName = "ionosctl-mongo-cluster-test-" + fake.AlphaNum(8)
	cidr               = fake.IP(fake.WithIPv4(), fake.WithIPCIDR("192.168.0.0/16")) + "/24"
	client             *config.Client
	createdClusterId   string
	createdDcId        string
)

// If your test is failing because your credentials env var seem empty, try running with `godotenv -f <config-file> go test <test>`
func TestMongoCommands(t *testing.T) {
	var err error
	client, err = config.GetClient()
	assert.NoError(t, err)
	go testMongoClusterCreateIdentifyRequiredNotSet(t)
	dcId, lanId, err := setupTestMongoCommands()
	if err != nil {
		t.Fatalf("Failed setting up Mongo required resources: %s", err)
	}
	//t.Cleanup(teardownTestMongoCommands)
	testMongoClusterCreate(t, dcId, lanId)
	testMongoUser(t)
}

func testMongoUser(t *testing.T) {
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

	c = user.UserDeleteCmd()
	c.Command.Flags().Set(constants.FlagClusterId, createdClusterId)
	c.Command.Flags().Set(constants.FlagName, name)
	err = c.Command.Execute()
	assert.NoError(t, err)
	createdUsers, _, err = client.MongoClient.UsersApi.ClustersUsersGet(context.Background(), createdClusterId).Execute()
	assert.NoError(t, err)
	assert.Empty(t, *createdUsers.GetItems())
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

func testMongoClusterCreateIdentifyRequiredNotSet(t *testing.T) {
	viper.Set(constants.ArgOutput, "text")
	viper.Set(constants.ArgCols, "Name")
	viper.Set(constants.ArgNoHeaders, true)
	fmt.Printf(viper.GetString(constants.ArgCols))

	c := cluster.ClusterCreateCmd()
	// Intentionally leave FlagDatacenterId unset, to see if err points us to this missing flag
	c.Command.Flags().Set(constants.FlagLanId, "foo")
	c.Command.Flags().Set(constants.FlagName, uniqueResourceName)
	c.Command.Flags().Set(constants.FlagInstances, "1")
	c.Command.Flags().Set(constants.FlagTemplateId, getPlaygroundTemplateUuid())
	c.Command.Flags().Set(constants.FlagMaintenanceDay, "Friday")
	c.Command.Flags().Set(constants.FlagMaintenanceTime, "10:00:00")
	c.Command.Flags().Set(constants.FlagCidr, cidr)

	err := c.Command.Execute()
	assert.ErrorContains(t, err, constants.FlagDatacenterId) // Assert that the error screams something about FlagDatacenterId
}

func setupTestMongoCommands() (string, string, error) {
	// make sure datacenter exists
	dcs, resp, err := client.CloudClient.DataCentersApi.DatacentersGet(context.Background()).Filter("name", uniqueResourceName).Depth(1).Execute()
	if resp.HttpNotFound() || len(*dcs.Items) < 1 {
		dc, _, err := client.CloudClient.DataCentersApi.DatacentersPost(context.Background()).Datacenter(sdkcompute.Datacenter{Properties: &sdkcompute.DatacenterProperties{Name: sdkcompute.PtrString(uniqueResourceName), Location: sdkcompute.PtrString("de/fra")}}).Execute()
		if err != nil {
			return createdDcId, "", fmt.Errorf("failed creating dc %w", err)
		}
		createdDcId = *dc.Id
		time.Sleep(10 * time.Second)
	} else if err != nil {
		return createdDcId, "", fmt.Errorf("failed getting dc %w", err)
	} else {
		createdDcId = *(*dcs.GetItems())[0].GetId()
	}

	fmt.Printf("dcId: %s\n", createdDcId)
	// make sure lan exists
	var lanId string
	lans, resp, err := client.CloudClient.LANsApi.DatacentersLansGet(context.Background(), createdDcId).Filter("name", uniqueResourceName).Depth(1).Execute()
	if resp.HttpNotFound() || len(*lans.Items) < 1 {
		lan, _, err := client.CloudClient.LANsApi.DatacentersLansPost(context.Background(), createdDcId).Lan(sdkcompute.LanPost{Properties: &sdkcompute.LanPropertiesPost{Name: sdkcompute.PtrString(uniqueResourceName), Public: sdkcompute.PtrBool(false)}}).Execute()
		if err != nil {
			return createdDcId, lanId, fmt.Errorf("failed creating lan: %w", err)
		}
		lanId = *lan.Id
		time.Sleep(10 * time.Second)
	} else if err != nil {
		return createdDcId, lanId, fmt.Errorf("failed getting lan: %w", err)
	} else {
		lanId = *(*lans.GetItems())[0].GetId()
	}
	fmt.Printf("lanId: %s\n", lanId)

	return createdDcId, lanId, nil
}

func teardownTestMongoCommands() {
	_, _, err := client.MongoClient.ClustersApi.ClustersDelete(context.Background(), createdClusterId).Execute()
	if err != nil {
		fmt.Printf("failed deleting cluster: %v\n", err)
	}

	time.Sleep(30 * time.Second) // Some clusters take longer to delete, and they still delete-protect datacenters. TODO: Use proxy API to slow down the deletion asynchronously for ~300secs

	_, err = client.CloudClient.DataCentersApi.DatacentersDelete(context.Background(), createdDcId).Execute()
	if err != nil {
		fmt.Printf("failed deleting dc: %v\n", err)
	}
}

// getPlaygroundTemplateUuid gets template ID of template "MongoDB Playground". In case of error, returns fallback UUID to that template as of 13 feb 2023
// If using this in a real production environment, you should remove 'fallbackUUID' and do some proper err handling
func getPlaygroundTemplateUuid() string {
	const fallbackUuid = "33457e53-1f8b-4ed2-8a12-2d42355aa759"

	client, err := config.GetClient()
	if err != nil {
		return fallbackUuid
	}
	ls, _, err := client.MongoClient.TemplatesApi.TemplatesGet(context.Background()).Execute()
	if err != nil {
		return fallbackUuid
	}
	for _, t := range *ls.Items {
		if *t.Edition == "playground" {
			return *t.Id
		}
	}

	return fallbackUuid
}
