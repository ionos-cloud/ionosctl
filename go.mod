module github.com/ionos-cloud/ionosctl

go 1.17

require (
	github.com/cheggaaa/pb/v3 v3.0.8
	github.com/fatih/color v1.13.0
	github.com/fatih/structs v1.1.0
	github.com/golang/mock v1.6.0
	github.com/ionos-cloud/sdk-go-auth v1.0.4
	github.com/ionos-cloud/sdk-go-dbaas-postgres v1.0.3
	github.com/ionos-cloud/sdk-go/v6 v6.1.0
	github.com/mitchellh/go-homedir v1.1.0
	github.com/spf13/cobra v1.3.0
	github.com/spf13/pflag v1.0.5
	github.com/spf13/viper v1.10.0
	github.com/stretchr/testify v1.7.0
	go.uber.org/multierr v1.6.0
	golang.org/x/crypto v0.0.0-20210921155107-089bfa567519
	golang.org/x/term v0.0.0-20210927222741-03fcf44c2211

)

require (
	github.com/VividCortex/ewma v1.1.1 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/fsnotify/fsnotify v1.5.1 // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/hashicorp/hcl v1.0.0 // indirect
	github.com/inconshreveable/mousetrap v1.0.0 // indirect
	github.com/magiconair/properties v1.8.5 // indirect
	github.com/mattn/go-colorable v0.1.12 // indirect
	github.com/mattn/go-isatty v0.0.14 // indirect
	github.com/mattn/go-runewidth v0.0.12 // indirect
	github.com/mitchellh/mapstructure v1.4.3 // indirect
	github.com/pelletier/go-toml v1.9.4 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/rivo/uniseg v0.2.0 // indirect
	github.com/spf13/afero v1.6.0 // indirect
	github.com/spf13/cast v1.4.1 // indirect
	github.com/spf13/jwalterweatherman v1.1.0 // indirect
	github.com/subosito/gotenv v1.2.0 // indirect
	go.uber.org/atomic v1.7.0 // indirect
	golang.org/x/net v0.0.0-20220722155237-a158d28d115b // indirect
	golang.org/x/oauth2 v0.0.0-20211104180415-d3ed0bb246c8 // indirect
	golang.org/x/sys v0.0.0-20220722155257-8c9f86f7a55f // indirect
	golang.org/x/text v0.3.7 // indirect
	google.golang.org/appengine v1.6.7 // indirect
	google.golang.org/protobuf v1.27.1 // indirect
	gopkg.in/ini.v1 v1.66.2 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b // indirect
)

require github.com/ionos-cloud/sdk-go-autoscaling v1.0.0

replace github.com/ionos-cloud/sdk-go-autoscaling => ../../../../../sdk/stackable
