package config

const (
	// Global Flags for Root Command
	ArgConfig      = "config"
	ArgOutput      = "output"
	ArgQuiet       = "quiet"
	ArgIgnoreStdin = "ignore-stdin"
	ArgServerUrl   = "server-url"
	ArgVerbose     = "verbose"
	// Data Center Flags
	ArgDataCenterId          = "id"
	ArgDataCenterName        = "name"
	ArgDataCenterDescription = "description"
	ArgDataCenterRegion      = "location"

	// Default values for flags
	DefaultApiURL         = "https://api.ionos.com/cloudapi/v5"
	DefaultConfigFileName = "/ionosctl-config.json"
	DefaultOutputFormat   = "text"

	Username = "userdata.name"
	Password = "userdata.password"
)
