//go:build integration
// +build integration

package dns_test

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/ionos-cloud/ionosctl/v6/commands/dns/record"
	"github.com/ionos-cloud/ionosctl/v6/commands/dns/utils"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	utils2 "github.com/ionos-cloud/ionosctl/v6/internal/utils"
	"github.com/ionos-cloud/ionosctl/v6/pkg/functional"
	dns "github.com/ionos-cloud/sdk-go-bundle/products/dns/v2"

	"github.com/cilium/fake"
	"github.com/ionos-cloud/ionosctl/v6/commands/dns/zone"
	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

var (
	sharedZ         dns.ZoneRead
	cl              *client.Client
	GoodToken       string
	tokCreationTime time.Time
)

func TestDNSCommands(t *testing.T) {
	err := setup()
	if err != nil {
		assert.FailNow(t, err.Error())
	}
	TestZone(t)   // sets sharedZ
	TestRecord(t) // uses sharedZ

	t.Cleanup(Cleanup)
}

// TODO: Improve me with the new config PR
func setup() error {
	if GoodToken = os.Getenv("IONOS_TOKEN"); GoodToken != "" {
		cl = client.NewClient("", "", GoodToken, "")
		return nil
	}

	// Otherwise, generate a token, since DNS doesn't function without it, only with username & password

	GoodUsername := os.Getenv("IONOS_USERNAME")
	GoodPassword := os.Getenv("IONOS_PASSWORD")

	if GoodUsername == "" || GoodPassword == "" {
		return fmt.Errorf("empty user/pass")
	}

	tempClNoToken := client.NewClient(GoodUsername, GoodPassword, "", "")
	tok, _, err := tempClNoToken.AuthClient.TokensApi.TokensGenerate(context.Background()).Execute()
	if err != nil {
		return fmt.Errorf("failed using username & password to generate a token for DNS tests: %w", err)
	}

	if tok.Token == nil {
		return fmt.Errorf("tok is nil")
	}

	GoodToken = *tok.Token
	tokCreationTime = time.Now().In(time.UTC).Add(-10 * time.Second)

	cl = client.NewClient("", "", *tok.Token, "")

	return nil
}

// Returns DNS Zone ID
func TestZone(t *testing.T) {
	var err error
	viper.Set(constants.ArgOutput, "text")
	viper.Set(constants.ArgToken, GoodToken)

	// === `ionosctl dns z create`
	c := zone.ZonesPostCmd()

	// // Verify name is required for zone creation
	// err := c.Command.Execute()
	// assert.ErrorContains(t, err, fmt.Sprintf("\"%s\" not set", constants.FlagName))

	// Generate a zone
	randName := fmt.Sprintf("%s%s.%s.space", fake.Adjective(), fake.Noun(), fake.AlphaNum(4))
	randDesc := fake.AlphaNum(32)
	c.Command.Flags().Set(constants.FlagName, randName)
	c.Command.Flags().Set(constants.FlagDescription, randDesc)
	err = c.Command.Execute()
	assert.NoError(t, err)

	// Try to find the zone created by the command
	zoneByName, _, err := client.Must().DnsClient.ZonesApi.ZonesGet(context.Background()).FilterZoneName(randName).Limit(1).Execute()

	assert.NoError(t, err)
	assert.NotEmpty(t, zoneByName.Items)
	sharedZ = (zoneByName.Items)[0]
	assert.NotEmpty(t, sharedZ.Properties)
	assert.Equal(t, randDesc, *sharedZ.Properties.Description)

	resolvedId, err := utils.ZoneResolve(randName)
	assert.NoError(t, err)
	assert.Equal(t, sharedZ.Id, resolvedId)

	// === `ionosctl dns z get`
	c = zone.ZonesFindByIdCmd()
	// // Verify id is required for zone get
	// err = c.Command.Execute()
	// assert.ErrorContains(t, err, fmt.Sprintf("\"%s\" not set", constants.FlagZone))

	// Try to see if ionosctl zone get finds newly created zone, using ID
	c.Command.Flags().Set(constants.FlagZone, sharedZ.Properties.ZoneName)
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
	c.Command.Flags().Set(constants.FlagZone, sharedZ.Properties.ZoneName)
	err = c.Command.Execute()
	assert.NoError(t, err)

	zoneThroughSdk, _, err := client.Must().DnsClient.ZonesApi.ZonesFindById(context.Background(), *sharedZ.Id).Execute()
	assert.NoError(t, err)
	assert.Equal(t, randDesc, *zoneThroughSdk.Properties.Description)

	resolvedId, err = utils.ZoneResolve(randName)
	assert.NoError(t, err)
	assert.Equal(t, *sharedZ.Id, resolvedId)
}

func TestRecord(t *testing.T) {
	var err error
	viper.Set(constants.CfgToken, GoodToken)
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
	c.Command.Flags().Set(constants.FlagZone, *sharedZ.Id)
	err = c.Command.Execute()
	assert.NoError(t, err)

	// Try to find the record created by the command
	recByName, _, err := client.Must().DnsClient.RecordsApi.RecordsGet(context.Background()).FilterName(randName).
		FilterZoneId(*sharedZ.Id).Limit(1).Execute()
	assert.NoError(t, err)
	assert.NotEmpty(t, recByName.Items)
	r := (recByName.Items)[0]
	assert.NotEmpty(t, r.Properties)
	assert.Equal(t, randIp, r.Properties.Content)

	// also test record.ZoneResolve
	resolvedId, err := record.Resolve(randName)
	assert.NoError(t, err)
	assert.Equal(t, r.Id, resolvedId)

	// `ionosctl dns r update`
	c = record.ZonesRecordsPutCmd()
	// // check prereqs
	// err = c.Command.Execute()
	// assert.ErrorContains(t, err, fmt.Sprintf("\"%s\", \"%s\" not set", constants.FlagRecord, constants.FlagZone))

	// try changing content of prev record
	randIp = fake.IP(fake.WithIPv4())
	c.Command.Flags().Set(constants.FlagContent, randIp)
	c.Command.Flags().Set(constants.FlagZone, *sharedZ.Id)
	c.Command.Flags().Set(constants.FlagRecord, r.Properties.Name) // test that querying by name works too
	err = c.Command.Execute()
	assert.NoError(t, err)
}

func Cleanup() {
	viper.Set(constants.ArgOutput, "text")

	ls, _, err := cl.DnsClient.ZonesApi.ZonesGet(context.Background()).Execute()
	if err != nil {
		log.Printf("Failed deletion: %s", err.Error())
	}

	err = functional.ApplyAndAggregateErrors(ls.Items,
		func(z dns.ZoneRead) error {
			_, _, err2 := cl.DnsClient.ZonesApi.ZonesDelete(context.Background(), *z.Id).Execute()
			return err2
		},
	)
	if err != nil {
		log.Printf("Failed deletion: %s", err.Error())
	}

	cleanupTokensCreatedAfterDate(tokCreationTime)
}

// TODO: Make some util func for me! It would also be useful for ionosctl users.
func cleanupTokensCreatedAfterDate(taym time.Time) {
	toks, _, err := cl.AuthClient.TokensApi.TokensGet(context.Background()).Execute()

	if err != nil {
		panic(err)
	}

	// Delete tokens generated since setup
	for _, t := range *toks.Tokens {
		date, err := utils2.ParseDate(*t.CreatedDate)
		if err != nil {
			panic(fmt.Errorf("couldn't parse date %s: %w", *t.CreatedDate, err))
		}

		// Delete the token if it was created after setup
		if date.After(taym) {
			_, _, err := cl.AuthClient.TokensApi.TokensDeleteById(context.Background(), *t.Id).Execute()
			if err != nil {
				panic(err)
			}
			break
		}
	}
}
