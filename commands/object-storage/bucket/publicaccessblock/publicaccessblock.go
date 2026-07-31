package publicaccessblock

import (
	"github.com/spf13/cobra"

	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/table"
)

var allCols = []table.Column{
	{Name: "BlockPublicAcls", JSONPath: "BlockPublicAcls", Default: true},
	{Name: "IgnorePublicAcls", JSONPath: "IgnorePublicAcls", Default: true},
	{Name: "BlockPublicPolicy", JSONPath: "BlockPublicPolicy", Default: true},
	{Name: "RestrictPublicBuckets", JSONPath: "RestrictPublicBuckets", Default: true},
}

type publicAccessBlockInfo struct {
	BlockPublicAcls       bool `json:"BlockPublicAcls"`
	IgnorePublicAcls      bool `json:"IgnorePublicAcls"`
	BlockPublicPolicy     bool `json:"BlockPublicPolicy"`
	RestrictPublicBuckets bool `json:"RestrictPublicBuckets"`
}

func PublicAccessBlockCmd() *core.Command {
	cmd := &core.Command{
		Command: &cobra.Command{
			Use:     "public-access-block",
			Aliases: []string{"pab"},
			Short:   "Manage guardrails that keep a bucket private (overrides ACLs/policies)",
			Long: `Manage a bucket's Public Access Block: a set of safety guardrails that take PRECEDENCE over ACLs and bucket policies to prevent unintended public exposure. When these settings are on, they win even if a policy or ACL would otherwise grant public access - making them the reliable way to guarantee a bucket stays private.

Four independent booleans:
  BlockPublicAcls        Reject new requests that would set a public ACL on the bucket or its objects.
  IgnorePublicAcls       Ignore any public ACLs already present (treat them as if absent).
  BlockPublicPolicy      Reject bucket policies that would grant public access.
  RestrictPublicBuckets  Restrict access through any already-public policy to authorized principals only.

Subcommands: 'put' applies a configuration (all four flags), 'get' shows the current values, 'delete' removes the block (removal re-exposes whatever ACLs/policies alone would allow, so remove with care).`,
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

	return cmd
}
