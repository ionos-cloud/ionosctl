package cluster

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/ionos-cloud/ionosctl/v6/commands/dbaas/mariadb-v2/completer"
	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/pkg/confirm"
	mariadb "github.com/ionos-cloud/sdk-go-bundle/products/dbaas/mariadb/v3"
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
		Namespace: "dbaas-mariadb-v2",
		Resource:  "cluster",
		Verb:      "restore",
		Aliases:   []string{"r"},
		ShortDesc: "Restore a MariaDB Cluster in place to a point in time",
		LongDesc: `Use this command to trigger an in-place restore of the specified MariaDB Cluster from its own backups.

Backups are not single points in time — they form a continuous recovery WINDOW (see a backup's earliestRecoveryTargetTime via ` + "`backup get`" + `). ` + "`--recovery-time`" + ` zooms into a specific moment inside that window; the cluster is rolled back to the nearest point at or before it. Omit ` + "`--recovery-time`" + ` (or use ` + "`now`" + `) to restore to the latest point. Accepted formats: ` + "`now`" + `, a date (` + "`2025-01-02`" + `), a date-time (` + "`\"2025-01-02 15:00\"`" + `), or a full RFC3339 timestamp (` + "`2025-01-02T15:00:00Z`" + `); values without a timezone are treated as UTC.

To instead create a NEW cluster from a specific backup, use ` + "`cluster create --backup-id`" + `.

Required values to run command:

* Cluster Id`,
		Example:    "ionosctl dbaas mariadb-v2 cluster restore --cluster-id <cluster-id> --recovery-time 2025-01-02T15:00:00Z",
		PreCmdRun:  PreRunClusterRestore,
		CmdRun:     RunClusterRestore,
		InitClient: true,
	})
	restoreCmd.AddUUIDFlag(constants.FlagClusterId, constants.FlagIdShort, "", constants.DescCluster, core.RequiredFlagOption(),
		core.WithCompletion(completer.ClusterIds, constants.MariaDBApiRegionalURL, constants.MariaDBLocations),
	)
	restoreCmd.AddStringFlag(constants.FlagRecoveryTime, constants.FlagRecoveryTimeShortPsql, "now",
		"Point inside the recovery window to restore to: 'now', a date, a date-time, or an RFC3339 timestamp (no timezone = UTC). The nearest point at or before this time is used; defaults to the latest",
		core.WithCompletionComplex(recoveryTimeCompletion, constants.MariaDBApiRegionalURL, constants.MariaDBLocations),
	)
	// Credentials are only re-sent when --password is given (see applyCredentialsFromFlags).
	restoreCmd.AddStringFlag(constants.ArgUser, "", "", "New username for the database user. Only applied together with --password")
	restoreCmd.AddStringFlag(constants.ArgPassword, "", "", "New password for the database user. When set, --user and --database must also be supplied (credentials are not returned by the API)")
	restoreCmd.AddStringFlag(constants.FlagDatabase, "", "", "Database for the credentials. Only applied together with --password")

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
	clusterRead, _, err := client.Must().MariaClientV2.ClustersApi.ClustersFindById(context.Background(), clusterId).Execute()
	if err != nil {
		return err
	}

	// The API rejects updates while the cluster is not ready with an opaque error.
	// Surface a clear message up front instead.
	if state := clusterRead.Metadata.State; state == nil || *state != mariadb.MARIADBCLUSTERSTATES_AVAILABLE {
		got := "unknown"
		if state != nil {
			got = string(*state)
		}
		return fmt.Errorf("cluster %s must be AVAILABLE to restore (current state: %s)", clusterId, got)
	}

	clusterProperties := clusterRead.Properties

	// In-place restore rolls the cluster back to a recovery timestamp within its own
	// backup window; sourceBackupId is left unset (that variant is only valid when
	// creating a new cluster — see `cluster create --backup-id`).
	restore := mariadb.NewMariadbRestoreClusterFromBackup()
	restore.RecoveryTargetDatetime = &mariadb.IonosTime{Time: recoveryTime}
	c.Verbose("Setting RecoveryTargetDatetime [RFC3339 format]: %v", recoveryTime.Format(time.RFC3339))
	clusterProperties.RestoreFromBackup = restore

	// Re-send credentials only if the user supplies a new --password (see helper).
	if err := applyCredentialsFromFlags(c, &clusterProperties); err != nil {
		return err
	}

	c.Verbose("Restoring Cluster in place...")

	clusterEnsure := mariadb.NewClusterEnsure(clusterId, clusterProperties)

	_, _, err = client.Must().MariaClientV2.ClustersApi.
		ClustersPut(context.Background(), clusterId).
		ClusterEnsure(*clusterEnsure).
		Execute()
	if err != nil {
		return err
	}

	c.Msg("MariaDB Cluster successfully restored")
	return nil
}

// recoveryTimeCompletion suggests useful --recovery-time values: "now" plus the
// earliest recovery timestamps of the target cluster's backups.
func recoveryTimeCompletion(cmd *cobra.Command, _ []string, _ string) ([]string, cobra.ShellCompDirective) {
	suggestions := []string{"now"}

	clusterId, _ := cmd.Flags().GetString(constants.FlagClusterId)
	req := client.Must().MariaClientV2.BackupsApi.BackupsGet(context.Background())
	if clusterId != "" {
		req = req.FilterClusterId(clusterId)
	}
	backups, _, err := req.Execute()
	if err != nil {
		return suggestions, cobra.ShellCompDirectiveNoFileComp
	}

	seen := map[string]bool{}
	for _, b := range backups.Items {
		if e := b.Properties.EarliestRecoveryTargetTime; e != nil {
			ts := e.Time.UTC().Format(time.RFC3339)
			if !seen[ts] {
				seen[ts] = true
				suggestions = append(suggestions, ts)
			}
		}
	}
	return suggestions, cobra.ShellCompDirectiveNoFileComp
}
