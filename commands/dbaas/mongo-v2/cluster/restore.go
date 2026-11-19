package cluster

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/ionos-cloud/ionosctl/v6/commands/dbaas/mongo-v2/completer"
	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/pkg/confirm"
	mongo "github.com/ionos-cloud/sdk-go-bundle/products/dbaas/mongo/v3"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// recoveryTimeLayouts are the timestamp formats accepted by --recovery-time, in
// addition to the "now"/"latest" keywords. Naive (no timezone) values are
// interpreted as UTC.
var recoveryTimeLayouts = []string{
	time.RFC3339,
	"2006-01-02T15:04:05",
	"2006-01-02 15:04:05",
	"2006-01-02 15:04",
	"2006-01-02",
}

func ClusterRestoreCmd() *core.Command {
	ctx := context.TODO()
	restoreCmd := core.NewCommand(ctx, nil, core.CommandBuilder{
		Namespace: "dbaas-mongo-v2",
		Resource:  "cluster",
		Verb:      "restore",
		Aliases:   []string{"r"},
		ShortDesc: "Restore a MongoDB Cluster in place to a point in time",
		LongDesc: `Use this command to trigger an in-place restore of the specified MongoDB Cluster from its own snapshots. The cluster's current data is overwritten with the restored data, and the cluster may be briefly unavailable while this happens.

Snapshots are not single points in time — each keeps the cluster's oplog, forming a continuous recovery WINDOW (see a snapshot's earliestRecoveryTargetTime / latestRecoveryTargetTime via ` + "`snapshot get`" + `). ` + "`--recovery-time`" + ` picks the exact moment inside that window to restore to. Omit ` + "`--recovery-time`" + ` (or use ` + "`now`" + `) to restore to the latest point. Accepted formats: ` + "`now`" + `, a date (` + "`2025-01-02`" + `), a date-time (` + "`\"2025-01-02 15:00\"`" + `), or a full RFC3339 timestamp (` + "`2025-01-02T15:00:00Z`" + `); values without a timezone are treated as UTC. Tab-completion suggests the window's boundaries.

To instead create a NEW cluster from a specific snapshot, use ` + "`cluster create --snapshot-id`" + `.

Required values to run command:

* Cluster Id`,
		Example:    "ionosctl dbaas mongo-v2 cluster restore --location <location> --cluster-id <cluster-id> --recovery-time 2025-01-02T15:00:00Z",
		PreCmdRun:  PreRunClusterRestore,
		CmdRun:     RunClusterRestore,
		InitClient: true,
	})
	restoreCmd.AddUUIDFlag(constants.FlagClusterId, constants.FlagIdShort, "", constants.DescCluster, core.RequiredFlagOption(),
		core.WithCompletion(completer.ClusterIds, constants.MongoDBApiRegionalURL, constants.MongoDBLocations),
	)
	restoreCmd.AddStringFlag(constants.FlagRecoveryTime, constants.FlagRecoveryTimeShortPsql, "now",
		"Point inside the recovery window to restore to: 'now', a date, a date-time, or an RFC3339 timestamp (no timezone = UTC). Defaults to the latest point",
		core.WithCompletionComplex(recoveryTimeCompletion, constants.MongoDBApiRegionalURL, constants.MongoDBLocations),
	)

	return restoreCmd
}

func PreRunClusterRestore(c *core.PreCommandConfig) error {
	return c.CheckRequiredFlagsAndLocation(constants.FlagClusterId)
}

// parseRecoveryTime turns a user-supplied --recovery-time into a UTC timestamp.
// It accepts the "now"/"latest" keywords (and empty, treated as now) plus the
// layouts in recoveryTimeLayouts; naive values are interpreted as UTC.
func parseRecoveryTime(s string) (time.Time, error) {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "", "now", "latest":
		return time.Now().UTC(), nil
	}
	for _, layout := range recoveryTimeLayouts {
		if t, err := time.Parse(layout, strings.TrimSpace(s)); err == nil {
			return t.UTC(), nil
		}
	}
	return time.Time{}, fmt.Errorf(
		"could not parse --recovery-time %q; use 'now', a date (2025-01-02), "+
			"a date-time (\"2025-01-02 15:00\"), or an RFC3339 timestamp (2025-01-02T15:00:00Z)", s)
}

// recoveryWindowEnd is the latest point a snapshot can restore to: its
// latestRecoveryTargetTime, or now while its window still extends to the present.
func recoveryWindowEnd(s mongo.Snapshot) time.Time {
	if l := s.LatestRecoveryTargetTime; l != nil && l.IsSet() && l.Get() != nil {
		return l.Get().UTC()
	}
	return time.Now().UTC()
}

// earliestRecoveryTime returns the earliest point any snapshot of the cluster can
// restore it to. ok is false when that cannot be determined (no snapshots yet, or
// the listing failed); the API then has the final say.
func earliestRecoveryTime(clusterId string) (earliest time.Time, ok bool) {
	snapshots, _, err := client.Must().MongoClientV2.SnapshotsApi.SnapshotsGet(context.Background()).FilterClusterId(clusterId).Execute()
	if err != nil {
		return time.Time{}, false
	}
	for _, s := range snapshots.Items {
		if e := s.Properties.EarliestRecoveryTargetTime; e != nil && (earliest.IsZero() || e.Before(earliest)) {
			earliest = e.Time
		}
	}
	return earliest.UTC(), !earliest.IsZero()
}

func RunClusterRestore(c *core.CommandConfig) error {
	clusterId := viper.GetString(core.GetFlagName(c.NS, constants.FlagClusterId))
	recoveryTime, err := parseRecoveryTime(viper.GetString(core.GetFlagName(c.NS, constants.FlagRecoveryTime)))
	if err != nil {
		return err
	}

	if earliest, ok := earliestRecoveryTime(clusterId); ok && recoveryTime.Before(earliest) {
		return fmt.Errorf("--%s %s is before the earliest point the cluster can be restored to (%s)",
			constants.FlagRecoveryTime, recoveryTime.Format(time.RFC3339), earliest.Format(time.RFC3339))
	}

	c.Verbose(constants.ClusterId, clusterId)

	if !confirm.FAsk(c.Command.Command.InOrStdin(), fmt.Sprintf("restore cluster with id: %v to %v, overwriting its current data", clusterId, recoveryTime.Format(time.RFC3339)), viper.GetBool(constants.ArgForce)) {
		return fmt.Errorf(confirm.UserDenied)
	}

	// Fetch existing cluster
	c.Verbose("Getting Cluster...")
	clusterRead, _, err := client.Must().MongoClientV2.ClustersApi.ClustersFindById(context.Background(), clusterId).Execute()
	if err != nil {
		return err
	}

	// A restore modifies the cluster, which must be ready for it. Surface a clear
	// message up front instead of an API error.
	if state := clusterRead.Metadata.State; state == nil || *state != mongo.MONGOCLUSTERSTATES_AVAILABLE {
		got := "unknown"
		if state != nil {
			got = string(*state)
		}
		return fmt.Errorf("cluster %s must be AVAILABLE to restore (current state: %s)", clusterId, got)
	}

	cluster := clusterRead.Properties
	props, err := editionProps(&cluster)
	if err != nil {
		return err
	}

	// In-place restore rolls the cluster back to a recovery timestamp within its own
	// snapshots' window; sourceSnapshotId is left unset (that variant is only valid
	// when creating a new cluster — see `cluster create --snapshot-id`).
	c.Verbose("Setting RecoveryTargetDatetime [RFC3339 format]: %v", recoveryTime.Format(time.RFC3339))
	props.SetRestoreFromSnapshot(mongo.MongoRestoreFromSnapshot{RecoveryTargetDatetime: &mongo.IonosTime{Time: recoveryTime}})

	// Leave the initial user unchanged: the fetched credentials lack the password.
	clearCredentials(&cluster)

	c.Verbose("Restoring Cluster in place...")

	clusterEnsure := mongo.NewClusterEnsure(clusterId, cluster)

	_, _, err = client.Must().MongoClientV2.ClustersApi.
		ClustersPut(context.Background(), clusterId).
		ClusterEnsure(*clusterEnsure).
		Execute()
	if err != nil {
		return err
	}

	c.Msg("MongoDB Cluster successfully restored")
	return nil
}

// recoveryTimeCompletion suggests useful --recovery-time values: "now" plus the
// earliest/latest recovery timestamps of the target cluster's snapshots.
func recoveryTimeCompletion(cmd *cobra.Command, _ []string, _ string) ([]string, cobra.ShellCompDirective) {
	suggestions := []string{"now"}

	clusterId, _ := cmd.Flags().GetString(constants.FlagClusterId)
	req := client.Must().MongoClientV2.SnapshotsApi.SnapshotsGet(context.Background())
	if clusterId != "" {
		req = req.FilterClusterId(clusterId)
	}
	snapshots, _, err := req.Execute()
	if err != nil {
		return suggestions, cobra.ShellCompDirectiveNoFileComp
	}

	seen := map[string]bool{}
	add := func(t time.Time) {
		ts := t.UTC().Format(time.RFC3339)
		if !seen[ts] {
			seen[ts] = true
			suggestions = append(suggestions, ts)
		}
	}
	for _, s := range snapshots.Items {
		if e := s.Properties.EarliestRecoveryTargetTime; e != nil {
			add(e.Time)
		}
		// latestRecoveryTargetTime is null while the window extends to the present.
		// Read the field directly: the SDK's GetLatestRecoveryTargetTimeOk
		// dereferences it without a nil check.
		if l := s.Properties.LatestRecoveryTargetTime; l != nil && l.IsSet() && l.Get() != nil {
			add(*l.Get())
		}
	}
	return suggestions, cobra.ShellCompDirectiveNoFileComp
}
