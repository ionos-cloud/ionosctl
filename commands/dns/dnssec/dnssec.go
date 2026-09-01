package dnssec

import (
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/table"
	"github.com/spf13/cobra"
)

var allCols = []table.Column{
	{Name: "Id", JSONPath: "id", Default: true},
	{Name: "KeyTag", Default: true, Format: metadataItemField("keyTag")},
	{Name: "DigestAlgorithmMnemonic", Default: true, Format: metadataItemField("digestAlgorithmMnemonic")},
	{Name: "Digest", Default: true, Format: metadataItemField("digest")},
	{Name: "Validity", JSONPath: "properties.validity", Default: true},
	{Name: "Flags", Format: metadataItemKeyDataField("flags")},
	{Name: "PubKey", Format: metadataItemKeyDataField("pubKey")},
	{Name: "ComposedKeyData", Format: metadataItemField("composedKeyData")},
	{Name: "Algorithm", JSONPath: "properties.keyParameters.algorithm"},
	{Name: "KskBits", JSONPath: "properties.keyParameters.kskBits"},
	{Name: "ZskBits", JSONPath: "properties.keyParameters.zskBits"},
	{Name: "NsecMode", JSONPath: "properties.nsecParameters.nsecMode"},
	{Name: "Nsec3Iterations", JSONPath: "properties.nsecParameters.nsec3Iterations"},
	{Name: "Nsec3SaltBits", JSONPath: "properties.nsecParameters.nsec3SaltBits"},
}

func metadataItemField(field string) table.FormatFunc {
	return func(item map[string]any) any {
		md, _ := item["metadata"].(map[string]any)
		if md == nil {
			return nil
		}
		items, _ := md["items"].([]any)
		if len(items) == 0 {
			return nil
		}
		first, _ := items[0].(map[string]any)
		if first == nil {
			return nil
		}
		return first[field]
	}
}

func metadataItemKeyDataField(field string) table.FormatFunc {
	return func(item map[string]any) any {
		md, _ := item["metadata"].(map[string]any)
		if md == nil {
			return nil
		}
		items, _ := md["items"].([]any)
		if len(items) == 0 {
			return nil
		}
		first, _ := items[0].(map[string]any)
		if first == nil {
			return nil
		}
		kd, _ := first["keyData"].(map[string]any)
		if kd == nil {
			return nil
		}
		return kd[field]
	}
}

func Root() *core.Command {
	cmd := &core.Command{
		Command: &cobra.Command{
			Use:     "dnssec",
			Aliases: []string{"sec", "dnskey", "key", "keys"},
			Short:   "Manage DNSSEC signing keys for a zone",
			Long: `Manage DNSSEC keys for a zone.

DNSSEC signs a zone so resolvers can cryptographically verify its answers were not forged. Enabling it generates two keys: a Key Signing Key (KSK) that signs the key set, and a Zone Signing Key (ZSK) that signs the records. To complete the chain of trust you must copy the resulting DS record (shown by 'dnssec get') to your domain's registrar/parent zone — until then, signing has no effect.

'create' enables DNSSEC, 'get' shows the keys and the DS record to publish, 'delete' turns it off.`,
			TraverseChildren: true,
		},
	}

	cmd.AddColsFlag(allCols)

	cmd.AddCommand(Get())
	cmd.AddCommand(Create())
	cmd.AddCommand(Delete())

	return cmd
}
