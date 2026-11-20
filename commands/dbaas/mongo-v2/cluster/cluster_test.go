package cluster

import (
	"fmt"
	"testing"
	"time"

	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/pkg/pointer"
	mongo "github.com/ionos-cloud/sdk-go-bundle/products/dbaas/mongo/v3"
)

func TestParseRecoveryTime(t *testing.T) {
	t.Run("keywords and empty default to ~now", func(t *testing.T) {
		for _, in := range []string{"", "now", "NOW", "latest", " Latest "} {
			got, err := parseRecoveryTime(in)
			if err != nil {
				t.Fatalf("parseRecoveryTime(%q) unexpected err: %v", in, err)
			}
			if d := time.Since(got); d < 0 || d > time.Minute {
				t.Errorf("parseRecoveryTime(%q) = %v, want ~now", in, got)
			}
			if got.Location() != time.UTC {
				t.Errorf("parseRecoveryTime(%q) not UTC: %v", in, got.Location())
			}
		}
	})

	cases := map[string]string{
		"2025-01-02T15:04:05Z":      "2025-01-02T15:04:05Z",
		"2025-01-02T15:04:05+02:00": "2025-01-02T13:04:05Z",
		"2025-01-02T15:04:05":       "2025-01-02T15:04:05Z",
		"2025-01-02 15:04:05":       "2025-01-02T15:04:05Z",
		"2025-01-02 15:04":          "2025-01-02T15:04:00Z",
		"2025-01-02":                "2025-01-02T00:00:00Z",
	}
	for in, want := range cases {
		t.Run("parses "+in, func(t *testing.T) {
			got, err := parseRecoveryTime(in)
			if err != nil {
				t.Fatalf("unexpected err: %v", err)
			}
			if got.Format(time.RFC3339) != want {
				t.Errorf("parseRecoveryTime(%q) = %s, want %s", in, got.Format(time.RFC3339), want)
			}
		})
	}

	t.Run("garbage errors", func(t *testing.T) {
		if _, err := parseRecoveryTime("yesterday"); err == nil {
			t.Error("expected error for unparseable input")
		}
	})
}

func TestToInt32GB(t *testing.T) {
	cases := map[string]int32{
		"10":   10,
		"8GB":  8,
		"16GB": 16,
		"1TB":  1024,
	}
	for in, want := range cases {
		got, err := toInt32GB(in, constants.FlagRam)
		if err != nil {
			t.Fatalf("toInt32GB(%q) unexpected err: %v", in, err)
		}
		if got != want {
			t.Errorf("toInt32GB(%q) = %d, want %d", in, got, want)
		}
	}

	if _, err := toInt32GB("not-a-size", constants.FlagRam); err == nil {
		t.Error("expected error for unparseable size")
	}
}

func TestValidateRetentionDays(t *testing.T) {
	for _, d := range []int32{1, 7, 365} {
		if err := validateRetentionDays(d); err != nil {
			t.Errorf("validateRetentionDays(%d) unexpected err: %v", d, err)
		}
	}
	for _, d := range []int32{0, -1, 366} {
		if err := validateRetentionDays(d); err == nil {
			t.Errorf("validateRetentionDays(%d) expected err, got nil", d)
		}
	}
}

func TestValidateInstanceAddresses(t *testing.T) {
	cases := []struct {
		name      string
		addresses int
		instances int32
		sharded   bool
		wantErr   bool
	}{
		{"one address per instance", 3, 3, false, false},
		{"single instance", 1, 1, false, false},
		{"too few addresses", 1, 3, false, true},
		{"too many addresses", 2, 1, false, true},
		{"sharded takes three", 3, 2, true, false},
		{"sharded ignores the instance count", 3, 5, true, false},
		{"sharded with one address", 1, 3, true, true},
		{"sharded with one address per instance", 5, 5, true, true},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			addresses := make([]string, c.addresses)
			for i := range addresses {
				addresses[i] = fmt.Sprintf("192.168.1.%d/24", 10+i)
			}
			err := validateInstanceAddresses(addresses, c.instances, c.sharded)
			if (err != nil) != c.wantErr {
				t.Errorf("validateInstanceAddresses(%d addresses, %d instances, sharded=%v) err = %v, wantErr %v",
					c.addresses, c.instances, c.sharded, err, c.wantErr)
			}
		})
	}
}

func TestValidateShardIncrease(t *testing.T) {
	sharded := mongo.MongoSharding{Enabled: pointer.From(true), Shards: pointer.From(int32(3))}
	cases := []struct {
		name    string
		current mongo.MongoSharding
		shards  int32
		wantErr bool
	}{
		{"same count", sharded, 3, false},
		{"adds shards", sharded, 4, false},
		{"removes a shard", sharded, 2, true},
		{"sharding disabled", mongo.MongoSharding{Enabled: pointer.From(false)}, 2, true},
		{"no sharding", mongo.MongoSharding{}, 2, true},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if err := validateShardIncrease(c.current, c.shards); (err != nil) != c.wantErr {
				t.Errorf("validateShardIncrease(%d) err = %v, wantErr %v", c.shards, err, c.wantErr)
			}
		})
	}
}

func TestRecoveryWindowEnd(t *testing.T) {
	end := time.Date(2026, 11, 18, 15, 0, 0, 0, time.FixedZone("CET", 3600))
	closed := mongo.Snapshot{LatestRecoveryTargetTime: &mongo.NullableIonosTime{NullableTime: *mongo.NewNullableTime(&end)}}
	if got := recoveryWindowEnd(closed); !got.Equal(end) || got.Location() != time.UTC {
		t.Errorf("recoveryWindowEnd(closed window) = %v, want %v in UTC", got, end)
	}

	for name, s := range map[string]mongo.Snapshot{
		"open window":   {LatestRecoveryTargetTime: &mongo.NullableIonosTime{}},
		"field missing": {},
	} {
		if d := time.Since(recoveryWindowEnd(s)); d < 0 || d > time.Minute {
			t.Errorf("recoveryWindowEnd(%s) is not ~now", name)
		}
	}
}

func TestEditionProps(t *testing.T) {
	for name, cluster := range map[string]mongo.Cluster{
		"enterprise": {EnterpriseCluster: &mongo.EnterpriseCluster{Name: "enterprise"}},
		"business":   {BusinessCluster: &mongo.BusinessCluster{Name: "business"}},
	} {
		props, err := editionProps(&cluster)
		if err != nil {
			t.Fatalf("editionProps(%s) unexpected err: %v", name, err)
		}
		if props.GetName() != name {
			t.Errorf("editionProps(%s).GetName() = %q", name, props.GetName())
		}
	}

	if _, err := editionProps(&mongo.Cluster{}); err == nil {
		t.Error("expected error for a cluster without edition properties")
	}
}

func TestClusterEdition(t *testing.T) {
	if got := clusterEdition(mongo.Cluster{EnterpriseCluster: &mongo.EnterpriseCluster{}}); got != mongo.MONGOEDITION_ENTERPRISE {
		t.Errorf("clusterEdition(enterprise) = %s", got)
	}
	if got := clusterEdition(mongo.Cluster{BusinessCluster: &mongo.BusinessCluster{}}); got != mongo.MONGOEDITION_BUSINESS {
		t.Errorf("clusterEdition(business) = %s", got)
	}
}

func TestClearCredentials(t *testing.T) {
	user := mongo.MongoUser{Username: "admin", Database: "testdb"}
	cluster := mongo.Cluster{
		EnterpriseCluster: &mongo.EnterpriseCluster{Credentials: &user},
		BusinessCluster:   &mongo.BusinessCluster{Credentials: &user},
	}
	clearCredentials(&cluster)
	if cluster.EnterpriseCluster.Credentials != nil || cluster.BusinessCluster.Credentials != nil {
		t.Error("credentials were not cleared")
	}
}

func TestInstanceCount(t *testing.T) {
	cases := map[string]struct {
		item map[string]any
		want any
	}{
		"enterprise": {map[string]any{"properties": map[string]any{"instances": map[string]any{"count": 3.0}}}, 3.0},
		"business":   {map[string]any{"properties": map[string]any{"instanceCount": 1.0}}, 1.0},
		"missing":    {map[string]any{"properties": map[string]any{}}, nil},
	}
	for name, c := range cases {
		if got := instanceCount(c.item); got != c.want {
			t.Errorf("instanceCount(%s) = %v, want %v", name, got, c.want)
		}
	}
}
