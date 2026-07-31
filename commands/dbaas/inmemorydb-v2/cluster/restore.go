package cluster

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/ionos-cloud/ionosctl/v6/commands/dbaas/inmemorydb-v2/completer"
	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/pkg/confirm"
	inmemorydb "github.com/ionos-cloud/sdk-go-bundle/products/dbaas/inmemorydb/v3"
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
		Namespace: "dbaas-inmemorydb-v2",
		Resource:  "cluster",
		Verb:      "restore",
		Aliases:   []string{"r"},
		ShortDesc: "Restore an In-Memory DB Cluster in place to a point in time",
		LongDesc: `Use this command to trigger an in-place restore of the specified In-Memory DB Cluster from its own snapshots.

The cluster is restored to the most recent snapshot taken at or before ` + "`--recovery-time`" + `. If ` + "`--recovery-time`" + ` is omitted it defaults to now (the latest available snapshot). Accepted formats: ` + "`now`" + `, a date (` + "`2025-01-02`" + `), a date-time (` + "`\"2025-01-02 15:00\"`" + `), or a full RFC3339 timestamp (` + "`2025-01-02T15:00:00Z`" + `); values without a timezone are treated as UTC. Tab-completion suggests the cluster's available recovery window.

To instead create a NEW cluster from a specific snapshot, use ` + "`cluster create --snapshot-id`" + `.

Required values to run command:

* Cluster Id`,
		Example:    "ionosctl dbaas in-memory-db-v2 cluster restore --cluster-id <cluster-id> --password <password>",
		PreCmdRun:  PreRunClusterRestore,
		CmdRun:     RunClusterRestore,
		InitClient: true,
	})
	restoreCmd.AddUUIDFlag(constants.FlagClusterId, constants.FlagIdShort, "", constants.DescCluster, core.RequiredFlagOption(),
		core.WithCompletion(completer.ClusterIds, constants.InMemoryDBApiRegionalURL, constants.InMemoryDBLocations),
	)
	restoreCmd.AddStringFlag(constants.FlagRecoveryTime, "", "now",
		"When to restore the cluster to: 'now', a date, a date-time, or an RFC3339 timestamp (no timezone = UTC). The nearest snapshot at or before this time is used",
		core.WithCompletionComplex(recoveryTimeCompletion, constants.InMemoryDBApiRegionalURL, constants.InMemoryDBLocations),
	)
	restoreCmd.AddStringFlag(constants.ArgUser, "", "", "Username for the In-Memory DB user. Defaults to the cluster's current username")
	restoreCmd.AddStringFlag(constants.ArgPassword, "", "", "Password for the In-Memory DB user. Required because the API does not return it on GET requests. Plaintext is hashed (SHA-256) client-side", core.RequiredFlagOption())

	return restoreCmd
}

func PreRunClusterRestore(c *core.PreCommandConfig) error {
	return c.CheckRequiredFlagsAndLocation(constants.FlagClusterId, constants.ArgPassword)
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

func RunClusterRestore(c *core.CommandConfig) error {
	clusterId := viper.GetString(core.GetFlagName(c.NS, constants.FlagClusterId))
	recoveryTime, err := parseRecoveryTime(viper.GetString(core.GetFlagName(c.NS, constants.FlagRecoveryTime)))
	if err != nil {
		return err
	}

	c.Verbose(constants.ClusterId, clusterId)

	if !confirm.FAsk(c.Command.Command.InOrStdin(), fmt.Sprintf("restore cluster with id: %v to %v", clusterId, recoveryTime.Format(time.RFC3339)), viper.GetBool(constants.ArgForce)) {
		return fmt.Errorf(confirm.UserDenied)
	}

	// Fetch existing cluster
	c.Verbose("Getting Cluster...")
	clusterRead, _, err := client.Must().InMemoryDBClientV2.ClustersApi.ClustersFindById(context.Background(), clusterId).Execute()
	if err != nil {
		return err
	}

	// The API rejects updates while the cluster is not ready with an opaque 409.
	// Surface a clear message up front instead.
	if state := clusterRead.Metadata.State; state == nil || *state != inmemorydb.CLUSTERSTATE_AVAILABLE {
		got := "unknown"
		if state != nil {
			got = string(*state)
		}
		return fmt.Errorf("cluster %s must be AVAILABLE to restore (current state: %s)", clusterId, got)
	}

	clusterProperties := clusterRead.Properties

	// In-place restore of an existing cluster uses inPlaceRestore (a recovery
	// timestamp), NOT sourceSnapshotId — that variant is only valid when creating
	// a new cluster (see `cluster create --snapshot-id`).
	inPlace := inmemorydb.NewInPlaceRestoreClusterFromSnapshot(recoveryTime)
	c.Verbose("Setting RecoveryTargetTime [RFC3339 format]: %v", recoveryTime.Format(time.RFC3339))
	clusterProperties.RestoreFromSnapshot = &inmemorydb.ClusterRestoreFromSnapshot{InPlaceRestoreClusterFromSnapshot: inPlace}

	// The API does not return the password on GET, so the fetched cluster carries
	// a password with an empty algorithm that a PUT would reject. Rebuild
	// credentials from the required --password.
	if err := applyCredentialsFromFlags(c, &clusterProperties); err != nil {
		return err
	}

	c.Verbose("Restoring Cluster in place...")

	clusterEnsure := inmemorydb.NewClusterEnsure(clusterId, clusterProperties)

	_, _, err = client.Must().InMemoryDBClientV2.ClustersApi.
		ClustersPut(context.Background(), clusterId).
		ClusterEnsure(*clusterEnsure).
		Execute()
	if err != nil {
		return err
	}

	c.Msg("In-Memory DB Cluster successfully restored")
	return nil
}

// recoveryTimeCompletion suggests useful --recovery-time values: "now" plus the
// earliest/latest recovery timestamps of the target cluster's snapshots.
func recoveryTimeCompletion(cmd *cobra.Command, _ []string, _ string) ([]string, cobra.ShellCompDirective) {
	suggestions := []string{"now"}

	clusterId, _ := cmd.Flags().GetString(constants.FlagClusterId)
	req := client.Must().InMemoryDBClientV2.SnapshotsApi.SnapshotsGet(context.Background())
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
		if l := s.Properties.LatestRecoveryTargetTime; l != nil && l.IsSet() && l.Get() != nil {
			add(*l.Get())
		}
	}
	return suggestions, cobra.ShellCompDirectiveNoFileComp
}
