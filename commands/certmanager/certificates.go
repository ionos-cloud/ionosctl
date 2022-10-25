package certificates

import (
	"github.com/ionos-cloud/ionosctl/pkg/core"
	"github.com/spf13/cobra"
)

func CertCmd() *core.Command {
	clusterCmd := &core.Command{
		Command: &cobra.Command{
			Use:              "certificate-manager",
			Short:            "Certificate Manager Operations",
			Long:             "The sub-commands of `ionosctl certificate-manager` allow you to manage the SSL Certificates under your account.",
			TraverseChildren: true,
		},
	}

	return clusterCmd
}
