package mongo

import (
	"bytes"
	"context"
	"fmt"
	"github.com/fatih/structs"
	"github.com/ionos-cloud/ionosctl/commands/dbaas/mongo/cluster"
	"github.com/ionos-cloud/ionosctl/pkg/config"
	"github.com/ionos-cloud/ionosctl/pkg/constants"
	sdkcompute "github.com/ionos-cloud/sdk-go/v6"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"testing"
)

const (
	uniqueResourceName = "ionosctl-mongo-cluster-test"
)

func TestMongoCommands(t *testing.T) {
	dcId, lanId, err := mongoClusterPrereqsMustExist() // Sets DatacenterId and LanId
	assert.NoError(t, err)
	defer cleanupMongoClusterPrereqs(dcId)

	viper.Set(constants.ArgOutput, "text")
	viper.Set(constants.ArgCols, structs.Names(cluster.ClusterPrint{}.DisplayName[0]))
	fmt.Printf(viper.GetString(constants.ArgCols))

	c := cluster.ClusterCreateCmd()
	c.Command.Flags().Set(constants.FlagDatacenterId, dcId)
	c.Command.Flags().Set(constants.FlagLanId, lanId)
	c.Command.Flags().Set(constants.FlagName, uniqueResourceName)
	c.Command.Flags().Set(constants.FlagInstances, "1")
	c.Command.Flags().Set(constants.FlagTemplateId, getPlaygroundTemplateUuid())
	c.Command.Flags().Set(constants.FlagMaintenanceDay, "Friday")
	c.Command.Flags().Set(constants.FlagMaintenanceTime, "10:00:00")
	c.Command.Flags().Set(constants.FlagCidr, "192.168.1.116/24")

	b := bytes.NewBufferString("")
	c.Command.SetOut(b)
	err = c.Command.Execute()
	assert.NoError(t, err)
}

func mongoClusterPrereqsMustExist() (string, string, error) {
	client, err := config.GetClient()
	if err != nil {
		return "", "", err
	}

	// make sure datacenter exists
	var dcId string
	dcs, resp, err := client.CloudClient.DataCentersApi.DatacentersGet(context.Background()).Filter("name", uniqueResourceName).Depth(1).Execute()
	if (err != nil && resp.HttpNotFound()) || len(*dcs.Items) < 1 {
		dc, _, err := client.CloudClient.DataCentersApi.DatacentersPost(context.Background()).Datacenter(sdkcompute.Datacenter{Properties: &sdkcompute.DatacenterProperties{Name: sdkcompute.PtrString(uniqueResourceName), Location: sdkcompute.PtrString("de/fra")}}).Execute()
		if err != nil {
			return dcId, "", nil
		}
		dcId = *dc.Id
	} else if err != nil {
		return dcId, "", nil
	} else {
		dcId = *(*dcs.GetItems())[0].GetId()
	}

	// make sure lan exists
	var lanId string
	lans, resp, err := client.CloudClient.LANsApi.DatacentersLansGet(context.Background(), dcId).Filter("name", uniqueResourceName).Depth(1).Execute()
	if (err != nil && resp.HttpNotFound()) || len(*lans.Items) < 1 {
		lan, _, err := client.CloudClient.LANsApi.DatacentersLansPost(context.Background(), dcId).Lan(sdkcompute.LanPost{Properties: &sdkcompute.LanPropertiesPost{Name: sdkcompute.PtrString(uniqueResourceName)}}).Execute()
		if err != nil {
			return dcId, lanId, nil
		}
		lanId = *lan.Id
	} else if err != nil {
		return dcId, lanId, nil
	} else {
		lanId = *(*lans.GetItems())[0].GetId()
	}

	return dcId, lanId, nil
}

func cleanupMongoClusterPrereqs(dcId string) error {
	client, err := config.GetClient()
	if err != nil {
		return err
	}
	_, err = client.CloudClient.DataCentersApi.DatacentersDelete(context.Background(), dcId).Execute()
	if err != nil {
		return err
	}
	return nil
}

// getPlaygroundTemplateUuid gets template ID of template "MongoDB Playground". In case of error, returns fallback UUID to that template as of 13 feb 2023
// If using this outside testing funcs, you should remove 'fallbackUUID' and do some proper err handling
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
