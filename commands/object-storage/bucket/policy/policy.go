package policy

import (
	"github.com/spf13/cobra"

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
			Use:     "policy",
			Aliases: []string{"pol"},
			Short:   "Manage the JSON access policy attached to a bucket",
			Long: `Manage a bucket's access policy: a single JSON document (S3/IAM policy syntax) that grants or denies access to the bucket and its objects. A policy is a list of statements, each with an Effect (Allow/Deny), a Principal (who: an account, or "*" for everyone/anonymous), an Action (which S3 operations, e.g. "s3:GetObject"), a Resource (which bucket/objects, as "arn:aws:s3:::<bucket>" or ".../*"), and optional Conditions.

The "s3:" action prefix and "arn:aws:s3:::" resource ARNs are the S3-compatible wire format required by the API - they are not references to AWS the company.

A public-access-block configuration takes precedence over a policy, so even an "Allow *" policy will not make a bucket public if public access is blocked.`,
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
