package server

import (
	"context"
	"strconv"

	"github.com/ionos-cloud/ionosctl/v6/commands/compute/completer"
	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	cloudapiv6 "github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func ServerCreateCmd() *core.Command {
	create := core.NewCommand(context.TODO(), nil, core.CommandBuilder{
		Namespace: "server",
		Resource:  "server",
		Verb:      "create",
		Aliases:   []string{"c"},
		ShortDesc: "Create an ENTERPRISE, VCPU, CUBE or GPU Server in a Virtual Data Center",
		LongDesc: `Use this command to create a Server in a Virtual Data Center (--datacenter-id). The compute shape you must describe depends on the server --type, and this is the single most important relationship to get right:

  * ENTERPRISE (the default type): you size the machine explicitly. --cores and --ram are REQUIRED. --cpu-family is optional (see below).
  * VCPU: you size the machine explicitly with --cores and --ram (REQUIRED), but you CANNOT set --cpu-family — the platform selects it.
  * CUBE: a fixed bundle. --template-id is REQUIRED and defines the cores, RAM and included NVMe Direct Attached Storage (DAS) boot volume; do NOT pass --cores/--ram. A DAS boot volume is created and attached automatically.
  * GPU: a fixed bundle with an attached GPU. --template-id is REQUIRED; --cpu-family is not accepted (the AMD_TURIN family is assigned automatically). A DAS boot volume is created and attached automatically.

RAM (--ram): must be a multiple of 256. The default unit is MB, so --ram 256 means 256MB; you may also give a unit, e.g. --ram 1GB. Minimum 256MB; the maximum depends on your contract limits. Applies to ENTERPRISE and VCPU only.

CPU family (--cpu-family): for ENTERPRISE, values such as INTEL_SKYLAKE, INTEL_XEON and AMD_OPTERON are valid, but availability differs per datacenter location. Run ` + "`" + `ionosctl compute location` + "`" + ` to see which families a location offers. If you omit it, the first family available in the datacenter's location is chosen for you.

Boot media / OS: for CUBE and GPU the DAS boot volume's OS comes from --image-id / --image-alias (or a bare --licence-type, default LINUX, for an empty volume). When you boot from a PUBLIC image you must also set --password or --ssh-key-paths so you can log in; PRIVATE images (which already contain credentials) need neither. To boot from a CD-ROM/ISO instead, attach one after creation with ` + "`" + `server cdrom attach` + "`" + `.

Confidential Computing (--confidential): creates an AMD SEV-SNP encrypted-memory VM. It is ENTERPRISE-only and requires --image-id pointing at a private, SEV-SNP-capable image; --cores and --cpu-family are NOT allowed because both are derived from the image's launch-config.json. Size its boot volume with --size and --storage-type.

Use ` + "`" + `--wait` + "`" + ` (` + "`" + `-w` + "`" + `) to wait for the Server to reach AVAILABLE state before the command returns.`,
		Example: `# BASIC: an ENTERPRISE server with 2 cores and 2GB RAM (CPU family auto-selected for the datacenter's location)
ionosctl compute server create --datacenter-id DATACENTER_ID --name web-01 --cores 2 --ram 2GB

# ADVANCED: an ENTERPRISE server pinned to a CPU family, booting a public image with an SSH key, waiting until it is AVAILABLE
ionosctl compute server create --datacenter-id DATACENTER_ID --name db-01 --cores 4 --ram 8GB \
  --cpu-family INTEL_SKYLAKE --image-id IMAGE_ID --ssh-key-paths ~/.ssh/id_rsa.pub --wait

# CUBE server: sizing comes entirely from the template bundle (do not pass --cores/--ram)
ionosctl compute server create --datacenter-id DATACENTER_ID --type CUBE --template-id TEMPLATE_ID

# VCPU server: explicit cores/RAM, no --cpu-family
ionosctl compute server create --datacenter-id DATACENTER_ID --type VCPU --cores 2 --ram 4GB

# GPU server: fixed template bundle; CPU family is set automatically
ionosctl compute server create --datacenter-id DATACENTER_ID --type GPU --template-id TEMPLATE_ID`,
		PreCmdRun:  PreRunServerCreate,
		CmdRun:     RunServerCreate,
		InitClient: true,
	})
	create.AddUUIDFlag(cloudapiv6.ArgDataCenterId, "", "", cloudapiv6.DatacenterId, core.RequiredFlagOption())
	_ = create.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataCentersIds(), cobra.ShellCompDirectiveNoFileComp
	})
	create.AddStringFlag(cloudapiv6.ArgName, cloudapiv6.ArgNameShort, "Unnamed Server", "Display name of the Server. Type-specific defaults are used if omitted (e.g. 'Unnamed Server', 'Unnamed Cube')")
	create.AddIntFlag(constants.FlagCores, "", cloudapiv6.DefaultServerCores, "The number of CPU cores. REQUIRED for ENTERPRISE and VCPU; ignored for CUBE/GPU (fixed by --template-id) and for --confidential (derived from the image). Maximum depends on your contract resource limits", core.RequiredFlagOption())
	create.AddStringFlag(constants.FlagRam, "", "", "Memory size, in multiples of 256. Default unit is MB (e.g. --ram 256 = 256MB); a unit may be given (e.g. --ram 2GB). Minimum 256MB, maximum per contract limit. REQUIRED for ENTERPRISE and VCPU; not used for CUBE/GPU (fixed by --template-id)", core.RequiredFlagOption())
	_ = create.Command.RegisterFlagCompletionFunc(constants.FlagRam, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"256MB", "512MB", "1024MB", "2GB", "3GB", "4GB", "5GB", "10GB", "16GB"}, cobra.ShellCompDirectiveNoFileComp
	})
	create.AddStringFlag(constants.FlagCpuFamily, "", cloudapiv6.DefaultServerCPUFamily,
		"CPU family for ENTERPRISE Servers, e.g. INTEL_SKYLAKE, INTEL_XEON, AMD_OPTERON. Availability varies by datacenter location (see `ionosctl compute location`). "+
			"Leave as AUTO to have the API pick the first family available in the datacenter's location. "+
			"Not accepted for VCPU (platform-chosen), GPU (always AMD_TURIN) or --confidential (image-derived); the API also rejects it for CUBE")
	_ = create.Command.RegisterFlagCompletionFunc(constants.FlagCpuFamily, func(*cobra.Command, []string, string) ([]string, cobra.ShellCompDirective) {
		datacenterId := viper.GetString(core.GetFlagName(create.NS, cloudapiv6.ArgDataCenterId))
		return completer.DatacenterCPUFamilies(create.Command.Context(), datacenterId), cobra.ShellCompDirectiveNoFileComp
	})
	create.AddBoolFlag(constants.FlagConfidential, "", false,
		"Create a Confidential Computing (SEV-SNP) VM from a confidential boot image. Requires --type ENTERPRISE and --image-id (a private, SEV-SNP image). "+
			"Do not set --cores or --cpu-family: both are derived from the image's launch-config.json. "+
			"A boot volume is created from --image-id and attached automatically; size it with --size and --storage-type.")
	create.AddStringFlag(cloudapiv6.ArgSize, "", strconv.Itoa(cloudapiv6.DefaultVolumeSize), "[Confidential] Size of the confidential boot volume. Default unit is GB, e.g. --size 10 or --size 10GB")
	create.AddSetFlag(constants.FlagStorageType, "", "HDD", []string{"HDD", "SSD", "SSD Standard", "SSD Premium"}, "[Confidential] Storage backing of the confidential boot volume. SSD tiers offer higher performance than HDD")
	create.AddBoolFlag(constants.FlagNICMultiQueue, "", false, constants.FlagNICMultiQueueDescription)
	create.AddStringFlag(constants.FlagAvailabilityZone, constants.FlagAvailabilityZoneShort, "AUTO", "Physical availability zone the Server is placed in. AUTO lets the platform choose; ZONE_1 / ZONE_2 pin it, letting you spread servers across zones for fault tolerance")
	_ = create.Command.RegisterFlagCompletionFunc(constants.FlagAvailabilityZone, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"AUTO", "ZONE_1", "ZONE_2"}, cobra.ShellCompDirectiveNoFileComp
	})
	create.AddUUIDFlag(cloudapiv6.ArgTemplateId, "", "", "[CUBE/GPU Server] Id of the instance-size Template that fixes cores, RAM and the included DAS boot volume. REQUIRED for --type CUBE and GPU. List available templates with `ionosctl compute template list`", core.RequiredFlagOption())
	_ = create.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgTemplateId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.TemplatesIds(), cobra.ShellCompDirectiveNoFileComp
	})
	create.AddSetFlag(constants.FlagType, "", "ENTERPRISE", []string{"ENTERPRISE", "CUBE", "VCPU", "GPU"}, "Compute model of the Server, which decides the required sizing flags. ENTERPRISE/VCPU need --cores and --ram; CUBE/GPU need --template-id instead (see the long description)")
	_ = create.Command.RegisterFlagCompletionFunc(constants.FlagType, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"ENTERPRISE", "CUBE", "VCPU", "GPU"}, cobra.ShellCompDirectiveNoFileComp
	})
	create.AddBoolFlag(constants.FlagPromoteVolume, "", false, "[CUBE/GPU Server] After creation, promote the auto-attached DAS volume to be the Server's boot volume. Requires --wait (the promotion is a follow-up PATCH once the server is AVAILABLE)")

	// Volume Properties - for DAS Volume associated with Cube Server
	create.AddStringFlag(cloudapiv6.ArgVolumeName, "N", "Unnamed Direct Attached Storage", "[CUBE Server] Display name of the included Direct Attached Storage (DAS) boot volume")
	create.AddStringFlag(cloudapiv6.ArgBus, "", "VIRTIO", "[CUBE Server] Bus the DAS volume is exposed on. VIRTIO is faster and recommended for modern OSes; IDE maximises compatibility with older ones")
	_ = create.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgBus, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"VIRTIO", "IDE"}, cobra.ShellCompDirectiveNoFileComp
	})
	create.AddSetFlag(cloudapiv6.ArgLicenceType, "l", "LINUX", constants.EnumLicenceType, "[CUBE Server] OS licence type of the DAS boot volume, used for an empty volume when no --image-id/--image-alias is given. Determines OS-specific billing/handling (e.g. LINUX, WINDOWS)")
	create.AddStringFlag(cloudapiv6.ArgImageAlias, cloudapiv6.ArgImageAliasShort, "", "[CUBE Server] Human-friendly alias of a public image to install on the DAS boot volume (e.g. ubuntu:latest), as an alternative to --image-id. Public images require --password or --ssh-key-paths")
	create.AddUUIDFlag(cloudapiv6.ArgImageId, "", "", "[CUBE Server] Id of an image or snapshot to install on the DAS boot volume. Public images require --password or --ssh-key-paths; private images (and snapshots) carry their own credentials. Also used as the confidential boot image when --confidential is set")
	_ = create.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgImageId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		confidential := viper.GetBool(core.GetFlagName(create.NS, constants.FlagConfidential))
		imageIds := completer.ImageIds(func(r ionoscloud.ApiImagesGetRequest) ionoscloud.ApiImagesGetRequest {
			// Completer for HDD images that are in the same location as the datacenter
			chosenDc, _, err := client.Must().CloudClient.DataCentersApi.DatacentersFindById(context.Background(),
				viper.GetString(core.GetFlagName(create.NS, cloudapiv6.ArgDataCenterId))).Execute()
			if err != nil || chosenDc.Properties == nil || chosenDc.Properties.Location == nil {
				return ionoscloud.ApiImagesGetRequest{}
			}

			r = r.Filter("location", *chosenDc.Properties.Location).Filter("imageType", "HDD")
			if confidential {
				// Only private, Confidential-Computing-capable images can boot a Confidential VM.
				r = r.Filter("public", "false").Filter("requiredFeatures", "SEV-SNP")
			}
			return r
		})

		if confidential {
			// Snapshots can't be confidential boot images; don't offer them.
			return imageIds, cobra.ShellCompDirectiveNoFileComp
		}

		snapshotIds := completer.SnapshotIds()

		return append(imageIds, snapshotIds...), cobra.ShellCompDirectiveNoFileComp
	})
	create.AddStringFlag(constants.ArgPassword, constants.ArgPasswordShort, "", "[CUBE Server] Root/Administrator password to set on the installed OS. Applies to PUBLIC images only, is set once at creation (not modifiable afterwards), and accepts characters a-z, A-Z, 0-9. Provide this and/or --ssh-key-paths when booting a public image")
	create.AddStringSliceFlag(cloudapiv6.ArgSshKeyPaths, cloudapiv6.ArgSshKeyPathsShort, []string{""}, "[CUBE Server] Paths to SSH public key files to inject into the DAS boot volume's OS (public images only). Comma-separate multiple paths. An alternative or complement to --password for logging in")

	return create
}
