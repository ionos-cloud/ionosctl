# Changelog

## [v6.6.4] (June 2023)

## Added
* Added support for DNS API
* Added the possibility of getting or deleting a token using the JWT directly: `--token`

### Fixed
- Deprecated warnings only show if the deprecated flag is being used

## [v6.6.3] (May 2023)

## Fixed
* Fixed token docs
* Fixed maintenance default, now maintenance is disabled by default for targetgroup target add
* Fix #288: improve client, config errors
* Fix #289  nodepool lan add --network flag using only last network
* 
## [6.6.2] (April 2023)

### Added
- Added the ability to add multiple networks CIDRs / gateway IPs for `ionosctl k8s np lan add`
- Added the possibility of listing all columns: `--cols all`
- Added the possibility of filtering for multiple values per key. For example, `--filters name=hello,name=world` would list resources which contain either name `hello`, or name `world`.

### Fixed
- Fixed broken columns for `ionosctl k8s np lan list`
- Fixed certain flag descriptions for `ionosctl alb rule httprule add`

### Changed
- Changed `ionosctl container-registry token create` defaults to more closely resemble `ionosctl token create`: `--no-headers true`, `--cols CredentialsPassword`, such that token will be the only output by default.

## [6.6.1] (April 2023)

### Added
- Added support for Container Registry API

### Fixed
- Fixed multiple issues related to image upload:
  - Fixed: Timeout for image diff finding (continuous polling on GET `/images`) after FTP upload is no longer hardcoded. This polling now respects the `--timeout` flag and uses context.
  - Fixed: Silent failing due to timeout for image diff finding
  - Changed: no longer throw error if any of the values in `--location` is non-IONOS, because `--ftp-url` is customizable.
  - Changed: improved image diff finding complexity to use the new SDK multi-value-per-key filtering
- Fixed `-o json` renaming `items` to `Resources`
- Fixed various flag names in the ALB, Targetgroup examples

### Changed
- Various documentation generation improvements. `summary.md` is now generated automatically.

## [6.6.0] (March 2023)

### Added
- Added support for DBaaS Mongo API
- Added support for Dataplatform API

### Changed
- flag values for `cols`, `filters` are now case insensitive

## [6.5.2] (March 2023)

### Fixed
- Fixed go.mod: added v6 as the major version
- Fixed cols flag on certain commands e.g. `user list`
- Fixed `group user list` command

### Dependencies
- Updated SDK Postgres to v1.1.1

## [6.5.1] (February 2023)

### Changed
- Changed `ionosctl version` behaviour to only display the version of the CLI by default  e.g.
     ```bash
      $ ionosctl version
      v6.5.0
     ```
     You can use -v/--verbose flag to display SDK versions.

### Fixed
- Added MaxResults flag to commands where it was missing from: `user s3key list`, `location cpu list`, etc.
- Query Parameters MaxResults and OrderBy won't be sent to CloudAPI if their values are 0 or "".

### Dependencies
- Updated go version to 1.19
- Bump sdk-go to v6.1.4, bump sdk-go-dbaas-postgres to v1.0.6. Various other dependency updates.

## [6.5.0] (January 2023)

### Changed
- **Important (affects scripts):** Slice type printing has been improved. Before: `[property1 property2 property3]`, now: `property1,property2,property3`. This means you can direct ionosctl slice output back to its own commands. Thanks to @avorima.
- Warnings while using `-o text` are now also piped to stderr, to keep consistent with `-o json`. Thanks to @webner

### Added
- Added support for Certificate Manager API: `ionosctl certificate-manager`

### Fixed
- Fixed list commands for Groups (@webner), Group Shares.
- Fixed a number of commands which used Viper to get the value of ipslice flags, including `natgateway create` (fixes #225).


## [6.4.2] (November 2022)

### Fixed
- Fixed type of Cidr flag, for DBaaS Postgres commands `cluster create`, `cluster update`. Thanks to: @maboehm


## [6.4.1] (November 2022)

### Fixed
- Fixed ionosctl ignoring auth environment variables if no config file present

## [6.4.0] (October 2022)

### Added

- Added `image` resource commands:
  - Added `image update` and `image delete` which correspond to CloudAPI Image Patch and Image Delete routes.
  - Added `image upload` command, which uploads your image to the desired IONOS FTP servers. Each Ionos FTP server corresponds to a `location`. These uploads can run in parallel, and by default this command also runs a PATCH on the freshly uploaded image, in order to simulate a `create` command.

### Fixed

- Fixed `CLIHttpUserAgent` containing duplicated `v` characters for version

- Fixed config file username & password being ignored if environment variable IONOS_TOKEN is set, and IONOS_USERNAME and IONOS_PASSWORD not set

## [6.3.3] (October 2022)

### Fixed
- Fixed viper sometimes not binding with pflag QueryParams defaults.

## [6.3.2] (September 2022)

### Added

- Added latest commit hash to `ionosctl version`, when ionosctl was built from source (`make build` or `make install`).
- Added support for file descriptors when using `-o json`. Stdout will contain only the API response, while stderr will contain CLI messages. For example, you can use `2> /dev/null` in combination with `-o json` to get rid of CLI messages such as wait for state messages, verbose messages, and other status messages. (Default behaviour remains unchanged)
- Added `UUID`, `IP`, `IPSlice` flag verifications, IonosCTL will throw more verbose errors now when parameter types are not matching a certain format.

### Fixed
- Fixed various bugs with the label command (#194)
    - Fixed conditional flag requirements for label --resource-type flag: now errors will be more verbose about what flags are required in conjunction with this flag.
    - Fixed filtering, maxResults, orderBy for label list

### Dependency updates
- Updated go version to 1.18
- Updated cobra to 1.5.0
- Updated viper to 1.13.0
- Updated all Ionos GO SDKs to use latest versions
- Bumped various other dependencies


## [6.3.1] (August 2022)

### Fixes
- Verbose messages for query parameters are now consistent
- Flag defaults for depth, orderBy, maxResults now work correctly
- List `--all`, Delete `--all`, Detach `--all` and similar commands now all use minimum depth (0) for the parent resource.

### New package manager support
- Added support and instructions for installing via scoop for Windows https://scoop.sh/


## [6.3.0] (August 2022)

### Enhancements
- reduced default depth for LIST operations to 1 and all other operations to 0


### Features
- added `-a`/`-all` flag to list all contract-level resources of a specific type without the need of providing dependent resource ID
  - supported resources: `k8s nodepool`, `share`, `server`, `lan`, `volume`, `loadbalancer`, `networkloadbalancer`, `applicationloadbalancer`
  - example: `ionosctl server list --all -F "vmState"="RUNNING"`
- added flag `--depth` (short `-D`) to control depth response. Useful in combination with `-o json`.

    _**Note:** Short flag `-D` not yet available for `firewallrule` command (belongs to `--destination-ip` flag.)_

- added support and instructions to install ionosctl from `snap` package manager
- added support and instructions to install ionosctl from `brew` package manager

### Deprecation Notice

- Short flag `-D` for `--destination-ip` for `firewallrule` is considered deprecated and will be replaced by `-D` for `--depth` pending the next major release.

## [6.2.0] (June 2022)

### Features
- new service: **Application Load Balancer (ALB)**
  - added the CLI commands for Application Load Balancer under `applicationloadblanacer` and `targetgroup` namespaces (PR #155):
    - `ionosctl applicationloadbalancer`
    - `ionosctl applicationloadbalancer flowlog`
    - `ionosctl applicationloadbalancer rule`
    - `ionosctl applicationloadbalancer rule httprule`
    - `ionosctl targetgroup `
    - `ionosctl targetgroup target`

### Dependency update
- updated SDK-Go version from [6.0.4](https://github.com/ionos-cloud/sdk-go/releases/tag/v6.0.4) to [v6.1.0](https://github.com/ionos-cloud/sdk-go/releases/tag/v6.1.0)

## [6.1.7] (May 2022)

### Features
- updated `ionosctl version` command to print SDKs versions
- removed `--public` option from `ionosctl k8s cluster create` command
- removed `--gateway-ip` option from `ionosctl k8s nodepool create` command
- added option to do certificate pinning by using `IONOS_PINNED_CERT` environment variable for commands.
  - Note: Set the `IONOS_PINNED_CERT` environment variable to be the public sha256 fingerprint of the certificate to be pinned.

### Dependency-update
- updated SDK-Go-Auth version from [1.0.3](https://github.com/ionos-cloud/sdk-go-auth/releases/tag/v1.0.3) to [v1.0.4](https://github.com/ionos-cloud/sdk-go-auth/releases/tag/v1.0.4)
- updated SDK-Go version from [6.0.2](https://github.com/ionos-cloud/sdk-go/releases/tag/v6.0.2) to [v6.0.4](https://github.com/ionos-cloud/sdk-go/releases/tag/v6.0.4)
- updated github.com/spf13/cobra version from [v1.2.1](https://github.com/spf13/cobra/releases/tag/v1.2.1) to [v1.3.0](https://github.com/spf13/cobra/releases/tag/v1.3.0)

## [6.1.7] (May 2022)

### Features
- updated `ionosctl version` command to print SDKs versions
- removed `--public` option from `ionosctl k8s cluster create` command
- removed `--gateway-ip` option from `ionosctl k8s nodepool create` command
- added option to do certificate pinning by using `IONOS_PINNED_CERT` environment variable for commands.
  - Note: Set the `IONOS_PINNED_CERT` environment variable to be the public sha256 fingerprint of the certificate to be pinned.

### Dependency-update
- updated SDK-Go-Auth version from [1.0.3](https://github.com/ionos-cloud/sdk-go-auth/releases/tag/v1.0.3) to [v1.0.4](https://github.com/ionos-cloud/sdk-go-auth/releases/tag/v1.0.4)
- updated SDK-Go version from [6.0.2](https://github.com/ionos-cloud/sdk-go/releases/tag/v6.0.2) to [v6.0.4](https://github.com/ionos-cloud/sdk-go/releases/tag/v6.0.4)
- updated github.com/spf13/cobra version from [v1.2.1](https://github.com/spf13/cobra/releases/tag/v1.2.1) to [v1.3.0](https://github.com/spf13/cobra/releases/tag/v1.3.0)

## [6.1.6] (May 2022)

### Features
- added new values for `--storage-type` option for `ionosctl dbaas postgres cluster create` command: **SSD_PREMIUM**, **SSD_STANDARD**.
  - Note: Value **SSD** is deprecated. Use the equivalent **SSD_PREMIUM** instead.
- added option to do certificate pinning by using `IONOS_PINNED_CERT` environment variable for the `ionosctl dbaas postgres` commands.
  - Note: Set the `IONOS_PINNED_CERT` environment variable to be the public sha256 fingerprint of the certificate to be pinned.

### Dependency-update
- updated SDK-Go-DBaaS Postgres version from [v1.0.2](https://github.com/ionos-cloud/sdk-go-dbaas-postgres/releases/tag/v1.0.2) to [v1.0.3](https://github.com/ionos-cloud/sdk-go-dbaas-postgres/releases/tag/v1.0.3)

## [6.1.5] (April 2022)

### Fixes
- added `--labels`,`--annotations` options for `ionosctl k8s nodepool create` command to set one or multiple labels, annotations
- added `--labels`,`--annotations` options for `ionosctl k8s nodepool update` command to set one or multiple labels, annotations (fixes [164](https://github.com/ionos-cloud/ionosctl/issues/164))
- added `Annotations,Labels` values for `--cols` option for `ionosctl k8s nodepool` commands

### Deprecated
- marked as deprecated the following options: `--label-key`,`--label-value`,`--annotation-key`,`--annotation-value` for `ionosctl k8s nodepool update` command (use `--labels`,`--annotations` options instead).

## [6.1.4] (April 2022)

- enhancement: added `--backup-location` option for `ionosctl dbaas postgres cluster create` command
- enhancement: added `--direction` option for `ionosctl dbaas postgres logs list` command
- enhancement: added `--since` and `--until` option for `ionosctl dbaas postgres logs list` command, to easily specify timeframe for getting logs
- enhancement: added `--public` option for `ionosctl k8s cluster create` command
- enhancement: added `--gateway-ip` option for `ionosctl k8s nodepool create` command
- enhancement: added `BootServerId` value for `--cols` option for `ionosctl volume` commands
- dependency-update: added SDK-Go-DBaaS Postgres version [v1.0.2](https://github.com/ionos-cloud/sdk-go-dbaas-postgres/releases/tag/v1.0.2)
- dependency-update: added SDK-Go version [v6.0.2](https://github.com/ionos-cloud/sdk-go/releases/tag/v6.0.2)
- dependency-update: updated Go version from 1.16 to 1.17

## [6.1.3]

- enhancement: added `--no-headers` option for GET commands, for text output (PR #158)

## [6.1.2]

- enhancement: added `SynchronizationMode` as `--cols` option for dbaas postgres cluster commands
- enhancement: renamed methods for token commands according to the new updates from [v1.0.2](https://github.com/ionos-cloud/sdk-go-auth/releases/tag/v1.0.2)
- dependency-update: added SDK-Go-DBaaS Postgres version [v1.0.1](https://github.com/ionos-cloud/sdk-go-dbaas-postgres/releases/tag/v1.0.1)
- dependency-update: added SDK-Go-Auth version [v1.0.3](https://github.com/ionos-cloud/sdk-go-auth/releases/tag/v1.0.3)

## [6.1.1]

- bug-fix: `ionosctl k8s cluster` command supports now `--cols` option

## [6.1.0]

- new service: **Database as a Service (DBaaS) - Postgres**
  - added the CLI commands for DBaaS Postgres under `dbaas postgres` namespace (PR #155):
    - `ionosctl dbaas postgres cluster`
    - `ionosctl dbaas postgres logs`
    - `ionosctl dbaas postgres backup`
    - `ionosctl dbaas postgres version`
    - `ionosctl dbaas postgres api-version`
- dependency-update: added SDK-Go-DBaaS Postgres version [v1.0.0](https://github.com/ionos-cloud/sdk-go-dbaas-postgres/releases/tag/v1.0.0)

## [6.0.2]

- enhancement: added `--no-headers` option for list commands, for text output (PR #153)
- documentation: separate documentation per service (PR #148)

## [6.0.1]

- enhancement: `--all` option iterates through all resources even if it hits error
- enhancement: improved logs messages on `--all` option and request info and request info

## [6.0.0]

- feature: added `--password` on `ionosctl user update` command
- feature: updated code for `ionosctl k8s nodepool` commands in sync with the changes from SDK Go
- bug-fix: `ionosctl lan create` command supports now `--cols` option
- dependency-update: added SDK-Go version `v6.0.0-beta.9` to `v6.0.0`

## [6.0.0-beta.8]

- feature: added `token` commands, added support for Auth API, to generate, list, delete Tokens
- dependency-update: added SDK-Go-Auth version v1.0.1

## [6.0.0-beta.7]

- feature: added `--filters`, `--max-results`, `--order-by` options on all list commands
- feature: added `-all` option for remove and detach commands
- enhancement: added completion support for `--k8s-version` option
- dependency-update: SDK-Go version from v6.0.0-beta.6 to v6.0.0-beta.9

## [6.0.0-beta.6]

- bug-fix: fixed `login` command to support username and password or token authentication

## [6.0.0-beta.5]

- Added `--all` option on delete commands
- Updated SDK-Go version to v6.0.0-beta.6
- Added `--image-alias` option to volume commands
- Removed `--public` and `--gateway-ip` options from k8s cluster commands
- Renamed `--ssh-keys` to `--ssh-key-paths` on volume commands and support uploading SSH Keys from files
- Added BootVolume, `--volume-id` and BootCdrom, `--cdrom-id` to server update command
- Renamed `--target-ip` to `--destination-ip`, `--type` to `--direction` from firewall rule commands
- Updated documentation with usage of boolean flags

## [6.0.0-beta.4]

- Added usage message on required flags
- Improved pkg modularization
- Added request time on verbose print
- Fixed [#113](https://github.com/ionos-cloud/ionosctl/issues/113)

## [6.0.0-beta.3]

- Added K8s Cluster security improvements
- Renamed `--bucket-name` flag to `--s3bucket` flag
- Added `--verbose` flag
- Updated Cobra version to [v1.2.1](https://github.com/spf13/cobra/releases/tag/v1.2.0), improving completions with descriptions
- Updated Go version to 1.16
- Updated SDK-Go version to v6.0.0-beta.4

## [6.0.0-beta2]

- Added Template, FlowLog, NAT Gateway, Network Load Balancer commands
- Updated Server commands to support Server of type CUBE
- Updated Datacenter, Location, Group, Contract, Kubernetes Node Pool Lan properties
- Updated Image, Request commands to support fetching the latest N Images/Requests
