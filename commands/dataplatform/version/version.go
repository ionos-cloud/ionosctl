package version

import (
	"context"
	"sort"
	"strconv"
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
	for i := 0; i < len(parts1) && i < len(parts2); i++ {
		num1, err1 := strconv.Atoi(parts1[i])
		num2, err2 := strconv.Atoi(parts2[i])

		// fall back to string comparison
		if err1 != nil || err2 != nil {
			return parts1[i] < parts2[i]
		}

		if num1 != num2 {
			return num1 < num2
		}
	}

	// if all parts equal, the version with fewer parts is considered older
	return len(parts1) < len(parts2)
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

	var versions []string
	if ls.Items != nil && len(*ls.Items) > 0 {
		versions = *(ls.Items)
	}

	return versions, err
}

func Versions() []string {
	ls, err := VersionsE()
	if err != nil {
		return nil
	}
	return ls
}
