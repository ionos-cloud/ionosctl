# Changelog

## [5.0.8]

- Added `--verbose` flag to all commands to see step-by-step process when running a command
- Updated Cobra version to [v1.2.1](https://github.com/spf13/cobra/releases/tag/v1.2.0), improving completions with descriptions
- Fixed bug [#99](https://github.com/ionos-cloud/ionosctl/issues/99)
- Added support for `IONOS_API_URL` environment variable to overwrite default [URL](https://api.ionos.com)

## [5.0.7]
 
- Added `api-subnets` and `s3bucket` options for `cluster` command
- Updated `request` command to print target resources
- Updated go version to 1.16
- Updated authentication mechanism to first check environment variables over config file
- Updated `sdk-go` version to v5.1.4
- Updated `image`, `request` commands to support fetching the latest N Images/Requests
- Updated `nodepool` command to be able to set multiple LAN Ids to a NodePool

## [5.0.6]

- Add support for `--ram` option for Server, NodePool resources. e.g. `--ram 256` or `--ram 256MB`
- Add options for auto-completion for `--ram` option
- Renamed `--ram-size` flag to `--ram`

## [5.0.5]

- Added commands aliases
- Added flags aliases
- Renamed flags
- Improved `--cols` option for output

## [5.0.4]

- Updated sdk-go version to v5.1.0
- Added commands for IpFailover, IpConsumer, CD-ROM commands
- Added missing properties for resources (e.g. `State`)

## [5.0.3]

- Updated sdk-go to v5.0.3
- Fixed typo `K8sFindBySClusterId`

## [5.0.2]

- Added commands for Kubernetes, BackupUnit, Private Cross-Connect, Contract Resources, User Management
- Updated commands structure to: `ionosctl server volume attach`, `ionosctl loadbalancer nic attach`
- Updated documentation structure
- Added `--wait-for-request` and `--wait-for-state` options
- Renamed `--ignore-stdin` flag to `--force`

## [5.0.1]

- Added commands for image, snapshot, ip block, firewall rule, label
- Added support for token authentication
- Updated `attach` commands for volume and nic

## [5.0.0]

- Added commands for data center, server, volume, nic, lan, load balancer, request
- Added completion support for flags and commands for Zsh, Fish, PowerShell and Bash terminals
- Added login command for SDK authentication
- Added `ionosctl` boilerplate
