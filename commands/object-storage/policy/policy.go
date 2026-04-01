package policy

import (
	"context"

	objectstorage "github.com/ionos-cloud/sdk-go-bundle/products/objectstorage/v2"
	"github.com/spf13/cobra"

	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/table"
)

var allCols = []table.Column{
	{Name: "Sid", JSONPath: "Sid", Default: true},
	{Name: "Effect", JSONPath: "Effect", Default: true},
	{Name: "Action", JSONPath: "Action", Default: true},
	{Name: "Resource", JSONPath: "Resource", Default: true},
	{Name: "Principal", JSONPath: "Principal", Default: true},
	{Name: "Condition", JSONPath: "Condition"},
}

var statusCols = []table.Column{
	{Name: "Bucket", JSONPath: "Bucket", Default: true},
	{Name: "IsPublic", JSONPath: "IsPublic", Default: true},
}

func PolicyCmd() *core.Command {
	cmd := &core.Command{
		Command: &cobra.Command{
			Use:              "policy",
			Aliases:          []string{"pol"},
			Short:            "Bucket policy operations for contract-owned object storage",
			TraverseChildren: true,
		},
	}

	cmd.Command.PersistentFlags().StringSlice(constants.ArgCols, nil, table.ColsMessage(allCols))
	_ = cmd.Command.RegisterFlagCompletionFunc(
		constants.ArgCols,
		func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return table.AllCols(allCols), cobra.ShellCompDirectiveNoFileComp
		},
	)

	cmd.AddCommand(GetCmd())
	cmd.AddCommand(PutCmd())
	cmd.AddCommand(DeleteCmd())
	cmd.AddCommand(StatusCmd())

	return cmd
}

// resolveRegionalClient resolves the bucket's region and returns a regional S3 client.
func resolveRegionalClient(ctx context.Context, name string) (*objectstorage.APIClient, error) {
	s3, err := client.GetObjectStorageClient("")
	if err != nil {
		return nil, err
	}

	loc, _, err := s3.BucketsApi.GetBucketLocation(ctx, name).Execute()
	if err != nil {
		return nil, err
	}

	region := ""
	if loc != nil {
		region = loc.GetLocationConstraint()
	}

	return client.GetObjectStorageClient(region)
}
