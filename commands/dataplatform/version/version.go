package version

import (
	"context"
	"sort"
	"strings"

	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/spf13/cobra"
)

func Root() *core.Command {
	cmd := &core.Command{
		Command: &cobra.Command{
			Use:              "version",
			Aliases:          []string{"v"},
			Short:            "Dataplatform Versions Operations",
			Long:             "This command allows you to view interact with versions for dataplatform clusters",
			TraverseChildren: true,
		},
	}

	cmd.AddCommand(List())
	cmd.AddCommand(Active())

	return cmd
}

func compareVersions(v1, v2 string) bool {
	parts1 := strings.Split(v1, ".")
	parts2 := strings.Split(v2, ".")
	if parts1[0] == parts2[0] {
		return parts1[1] < parts2[1]
	}
	return parts1[0] < parts2[0]
}

func Latest(versions []string) string {
	if len(versions) == 0 {
		return ""
	}
	sort.Slice(versions, func(i, j int) bool {
		return compareVersions(versions[i], versions[j])
	})
	return versions[len(versions)-1]
}

func VersionsE() ([]string, error) {
	ls, _, err := client.Must().DataplatformClient.DataPlatformMetaDataApi.VersionsGet(context.Background()).Execute()
	return ls, err
}

func Versions() []string {
	ls, err := VersionsE()
	if err != nil {
		return nil
	}
	return ls
}
