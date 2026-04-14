package backup

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/ionos-cloud/ionosctl/v6/commands/dbaas/mariadb/cluster"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/table"
	"github.com/ionos-cloud/ionosctl/v6/pkg/functional"
	"github.com/ionos-cloud/sdk-go-bundle/products/dbaas/mariadb/v2"

	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/spf13/cobra"
)

func Root() *core.Command {
	cmd := &core.Command{
		Command: &cobra.Command{
			Use:              "backup",
			Aliases:          []string{"b"},
			Short:            "MariaDB Backup Operations",
			Long:             "The sub-commands of `ionosctl dbaas mariadb backup` allow you to manage the MariaDB Cluster Backups under your account.",
			TraverseChildren: true,
		},
	}

	cmd.Command.PersistentFlags().StringSlice(constants.ArgCols, nil, table.ColsMessage(allCols))
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return table.AllCols(allCols), cobra.ShellCompDirectiveNoFileComp
	})

	cmd.AddCommand(List())
	cmd.AddCommand(Get())

	return cmd

}

var allCols = []table.Column{
	{Name: "BackupId", JSONPath: "id", Default: true},
	{Name: "ClusterId", JSONPath: "properties.clusterId", Default: true},
	{Name: "Size", Default: true, Format: func(item map[string]any) any {
		v := table.Navigate(item, "properties.size")
		if v == nil {
			return nil
		}
		f, ok := v.(float64)
		if !ok {
			return v
		}
		return fmt.Sprintf("%d MiB", int(f))
	}},
	{Name: "Items", Default: true, Format: func(item map[string]any) any {
		baseBackups, ok := table.Navigate(item, "properties.baseBackups").([]any)
		if !ok || len(baseBackups) == 0 {
			return nil
		}
		var parts []string
		for _, bb := range baseBackups {
			m, ok := bb.(map[string]any)
			if !ok {
				continue
			}
			createdStr, _ := m["created"].(string)
			size, _ := m["size"].(float64)
			if createdStr == "" {
				continue
			}
			t, err := time.Parse(time.RFC3339, createdStr)
			if err != nil {
				continue
			}
			parts = append(parts, fmt.Sprintf("%s (%d MiB)", humanize.Time(t), int(size)))
		}
		return strings.Join(parts, ", ")
	}},
}

func Backups(fs ...Filter) (mariadb.BackupList, error) {
	cs, err := cluster.Clusters()

	if len(cs.Items) == 0 {
		return mariadb.BackupList{}, nil // no clusters -> empty response with no error
	}

	req := client.Must().MariaClient.BackupsApi.BackupsGet(context.Background())
	for _, f := range fs {
		req = f(req)
	}

	bs, _, err := req.Execute()

	return bs, err
}

type Filter func(request mariadb.ApiBackupsGetRequest) mariadb.ApiBackupsGetRequest

func BackupsProperty[V any](f func(c mariadb.BackupResponse) V, fs ...Filter) []V {
	recs, err := Backups(fs...)
	if err != nil {
		return nil
	}
	return functional.Map(recs.Items, f)
}
