//go:build integration
// +build integration

package dns_test

import (
	"context"
	"log"
	"strconv"
	"testing"

	"github.com/ionos-cloud/ionosctl/v6/commands/dns/record"
	"github.com/ionos-cloud/ionosctl/v6/internal/functional"
	dns "github.com/ionos-cloud/sdk-go-dns"

	"github.com/cilium/fake"
	"github.com/ionos-cloud/ionosctl/v6/commands/dns/zone"
	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/pkg/constants"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

var (
	sharedZ dns.ZoneRead
)

func TestDNSCommands(t *testing.T) {
	TestZone(t)   // sets sharedZ
	TestRecord(t) // uses sharedZ

	Cleanup(t) // Cleanup using ionosctl commands, otherwise use SDK directly
}

// Returns DNS Zone ID
func TestZone(t *testing.T) {
	var err error
	viper.Set(constants.ArgOutput, "text")

	// === `ionosctl dns z create`
	c := zone.ZonesPostCmd()

	// // Verify name is required for zone creation
	// err := c.Command.Execute()
	// assert.ErrorContains(t, err, fmt.Sprintf("\"%s\" not set", constants.FlagName))

	// Generate a zone
	randName := fake.Adjective() + ".com"
	randDesc := fake.AlphaNum(32)
	c.Command.Flags().Set(constants.FlagName, randName)
	c.Command.Flags().Set(constants.FlagDescription, randDesc)
	err = c.Command.Execute()
	assert.NoError(t, err)

	// Try to find the zone created by the command
	zoneByName, _, err := client.Must().DnsClient.ZonesApi.ZonesGet(context.Background()).FilterZoneName(randName).Limit(1).Execute()
	assert.NoError(t, err)
	assert.NotEmpty(t, *zoneByName.Items)
	sharedZ = (*zoneByName.Items)[0]
	assert.NotEmpty(t, sharedZ.Properties)
	assert.Equal(t, randDesc, *sharedZ.Properties.Description)

	resolvedId, err := zone.Resolve(randName)
	assert.NoError(t, err)
	assert.Equal(t, *sharedZ.Id, resolvedId) // I added these 3 lines later - to test zone.Resolve too

	// === `ionosctl dns z get`
	c = zone.ZonesFindByIdCmd()
	// // Verify id is required for zone get
	// err = c.Command.Execute()
	// assert.ErrorContains(t, err, fmt.Sprintf("\"%s\" not set", constants.FlagZone))

	// Try to see if ionosctl zone get finds newly created zone, using ID
	c.Command.Flags().Set(constants.FlagZone, *sharedZ.Properties.ZoneName)
	err = c.Command.Execute()
	assert.NoError(t, err)
	// TODO: I can't change command output to a buffer and check correctness, because output buffer is hardcoded in command runner

	// === `ionosctl dns z update`
	c = zone.ZonesPutCmd()
	// // Check `ionosctl dns z update` prereqs
	// err = c.Command.Execute()
	// assert.ErrorContains(t, err, fmt.Sprintf("\"%s\" not set", constants.FlagZone))

	// Try changing desc using `ionosctl dns z update`
	randDesc = fake.AlphaNum(32)
	c.Command.Flags().Set(constants.FlagDescription, randDesc)
	c.Command.Flags().Set(constants.FlagZone, *sharedZ.Properties.ZoneName)
	err = c.Command.Execute()
	assert.NoError(t, err)

	zoneThroughSdk, _, err := client.Must().DnsClient.ZonesApi.ZonesFindById(context.Background(), *sharedZ.Id).Execute()
	assert.NoError(t, err)
	assert.Equal(t, randDesc, *zoneThroughSdk.Properties.Description)

	// Try changing name using `ionosctl dns z update`
	randName = fake.AlphaNum(32)
	c.Command.Flags().Set(constants.FlagName, randName)
	c.Command.Flags().Set(constants.FlagZone, *sharedZ.Properties.ZoneName)
	err = c.Command.Execute()
	assert.NoError(t, err)

	zoneThroughSdk, _, err = client.Must().DnsClient.ZonesApi.ZonesFindById(context.Background(), *sharedZ.Id).Execute()
	assert.NoError(t, err)
	assert.Equal(t, randName, *zoneThroughSdk.Properties.ZoneName)

	resolvedId, err = zone.Resolve(randName)
	assert.NoError(t, err)
	assert.Equal(t, *sharedZ.Id, resolvedId)
}

func TestRecord(t *testing.T) {
	var err error
	viper.Set(constants.ArgOutput, "text")

	// `ionosctl dns r create`
	c := record.ZonesRecordsPostCmd()
	// // Check reqs
	// err := c.Command.Execute()
	// assert.ErrorContains(t, err, fmt.Sprintf("\"%s\", \"%s\", \"%s\", \"%s\" not set", constants.FlagZone, constants.FlagName, constants.FlagContent, constants.FlagType))

	// Generate a record
	randIp := fake.IP(fake.WithIPv4())
	randName := fake.Adjective() + "-" + strconv.Itoa(int(fake.Port(fake.WithPortDynamic())))
	c.Command.Flags().Set(constants.FlagContent, randIp)
	c.Command.Flags().Set(constants.FlagType, "A")
	c.Command.Flags().Set(constants.FlagName, randName)
	c.Command.Flags().Set(constants.FlagZoneId, *sharedZ.Id)
	err = c.Command.Execute()
	assert.NoError(t, err)

	// Try to find the record created by the command
	recByName, _, err := client.Must().DnsClient.RecordsApi.RecordsGet(context.Background()).FilterName(randName).
		FilterZoneId(*sharedZ.Id).Limit(1).Execute()
	assert.NoError(t, err)
	assert.NotEmpty(t, *recByName.Items)
	r := (*recByName.Items)[0]
	assert.NotEmpty(t, r.Properties)
	assert.Equal(t, randIp, *r.Properties.Content)

	// also test record.Resolve
	resolvedId, err := record.Resolve(randName)
	assert.NoError(t, err)
	assert.Equal(t, *r.Id, resolvedId)

	// `ionosctl dns r update`
	c = record.ZonesRecordsPutCmd()
	// // check prereqs
	// err = c.Command.Execute()
	// assert.ErrorContains(t, err, fmt.Sprintf("\"%s\", \"%s\" not set", constants.FlagRecord, constants.FlagZone))

	// try changing content of prev record
	randIp = fake.IP(fake.WithIPv4())
	c.Command.Flags().Set(constants.FlagContent, randIp)
	c.Command.Flags().Set(constants.FlagZoneId, *sharedZ.Id)
	c.Command.Flags().Set(constants.FlagRecord, *r.Properties.Name) // test that querying by name works too
	err = c.Command.Execute()
	assert.NoError(t, err)
}

func Cleanup(t *testing.T) {
	viper.Set(constants.ArgOutput, "text")

	ls, _, err := client.Must().DnsClient.ZonesApi.ZonesGet(context.Background()).Execute()
	assert.NoError(t, err)

	err = functional.ApplyAndAggregateErrors(*ls.Items,
		func(z dns.ZoneRead) error {
			_, err2 := client.Must().DnsClient.ZonesApi.ZonesDelete(context.Background(), *z.Id).Execute()
			return err2
		},
	)
	if err != nil {
		log.Printf("Failed deletion: %s", err.Error())
	}
}
