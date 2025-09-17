# Changelog

## upcoming

### Fixed
- Fixed missing column 'DestinationIp' for 'firewall rule' command
- Fixed missing columns 'RAM', 'PostgresVersion', 'Cores' for 'dbaas postgres' command
- Remove 'dbaas mongo template list' StorageSize column conversion to GB as the API now already returns GB values

### Changed
- Changed 'Ram' to 'RAM' for 'server', 'template' commands for consistency.
- Changed default RAM size to the new minimum of 4GB for 'dbaas postgres'.

## [v6.9.3] – September 2025

### Fixed
- Fixed a bug where 'ionosctl dbaas postgres cluster create' would result in a 404.

### Changed
- Under 'dbass postgres', the shorthand flags of
     "--datacenter-id" ("D"), "instances", ("I"), "backup-location" ("B"),
     "maintenance-time" ("T"), "maintenance-day" ("d"), "version" ("V"),
     "recovery-time" ("R"), "backup-id" ("b"), "db-username" ("U"), "db-password" ("P")
     have been deprecated and will be removed in a future release

- Under 'token delete', the shorthand flag of "--all" ("A") has been deprecated and will be removed in a future release,
     while the standard shorthand for "--all" ("a") has been added.

- Under 'token', the shorthand flags of "--current" ("C") and "--expired" ("E")
     have been deprecated and will be removed in a future release

### Dependencies
- All dependencies bumped to their latest versions

## [v6.9.2] – August 2025

### Added
- Added support for Observability Monitoring
- Added support for 'central' under 'ionosctl logging-service'
- Added support for deleting all labels of all compatible resources with `ionosctl label delete --all`, without specifying a resource type.
- Added missing locations for all regional APIs
- Added support for 'de/fra/2' location

### Dependencies
- Bump all dependencies, including SDKs, to their latest versions
- Bump minimum go version to 1.24.5

## Changed
- Rework image upload `--location` logic to support both short and API-style location identifiers (e.g. `vit` and `es/vit`) when resolving FTP and API endpoints, which also allows using 'de/fra/2' location.
- Improve verbose messages for 'image upload' command.
- Improve 'dbaas mongo cluster update' help text by @fepape-ionos.
- Removed 'dataplatform' commands, as the Dataplatform API has been sunsetted.

## [v6.9.1] – July 2025

### Changed
- Changed config generation product name from 'compute' to 'cloud' for CloudAPI, though existing configurations using 'compute' key will continue to work.

### Fixed
- Fixed a bug where the fallback to IONOS_CONFIG_FILE and ~/.ionos/config was not working correctly
- Fixed 'ionosctl mongo cluster update' sending a payload containing certain nil/empty values to the API, which would cause a 500 Internal Server Error.
- Fixed a bug where config overrides were ignored for certain CloudAPI commands


## [v6.9.0] – July 2025

> [!WARNING]
> This version introduces important changes to the config subsystem.
> **Legacy JSON configs** (`~/.ionosctl/config.json` or `~/snap/…/config.json`) are no longer supported and should be removed.
> On first run, existing credentials will be ***copied over*** into the new YAML format, but you should run:
> ```bash
> ionosctl config login
> ```
> to regenerate a fresh config and take full advantage of the new features.
> as well as
> ```bash
> ionosctl logout --only-purge-old'
> ```
> to delete the old config file.
'

### Impact and Migration Notes

- Credentials found in any legacy `config.json` ***will be carried over***.
- Although you can continue to use the CLI as is, it is highly recommended to delete all `config.json` files and re-run `ionosctl login` to produce a YAML config with all the new features.

### Added

- Added support for the new **SDK YAML config** layout:
  - Generate a complete config to a file via `ionosctl config login`.
  - Generate an example config to stdout via `ionosctl config login --example`.
  - Migrate existing credentials from `config.json` with a one-time deprecation warning.
  - Support ***multiple profiles*** and ***multiple environments*** with ***per-product and per-location URL overrides*** in a single YAML.
  - For now, use IONOS_CURRENT_PROFILE to set the current profile, though better CLI support will be added in the future.
  - Adds fallbacks to IONOS_CONFIG_FILE (SDKs config env var) and ~/.ionos/config (SDKs config location) for config file location, if not found at flag '--config'.

- Added `--example` flag to `ionosctl config login`:
  - Prints a sample YAML config to stdout without authenticating or writing any file.

- Added a spinner loading animation for `ionosctl config login` only if generating to a file and API index polling takes more than 250ms.

- Added **filtering** options to `ionosctl config login`:
  - `--filter-version`
  - `--filter-visibility`
  - `--filter-gate`
  - `--whitelist` / `--blacklist`
  - `--custom-names` for remapping API names in the generated config.

### Changed

- Changed (subcommands) README broken link path to "For more information, see **SUBCOMMANDS** section for respective products." in Introduction and README.md file.
- Changed **authentication precedence** and updated `whoami` to reflect:
  1. `IONOS_TOKEN` env var
  2. `IONOS_USERNAME` + `IONOS_PASSWORD` env vars
  3. Token in YAML config
  4. Username/password in YAML config

- Changed `ionosctl config logout` to:
  - Clear credentials from the YAML.
  - Detect side-by-side `config.json`, prompt for deletion, and remove on confirmation.

- Changed `ionosctl config location` resolution order to:
  1. `--config` flag
  2. `IONOS_CONFIG_FILE` env var
  3. SDK default (`~/.ionos/config.yaml`)
  - If no file exists, it still prints the `--config` value to avoid breaking changes.

- Changed Certificate Manager API to /v2
- Certificate Manager commands are now nested under `certificate` resource, but the old commands are still available (though hidden in the helptext) for backwards compatibility
- Certificate Manager command `api-version` is now no longer available, has been hidden and deprecated, and using it will print a warning as well as a dummy value `v2.0`
- Added a few friendly certificate manager aliases

## [v6.8.6] (June 2025)

### Added
- Added '--version' flag to 'dbaas mongo cluster update' to allow updating the cluster version (#532 - @fepape-ionos)
- Added tab-completions for '--version' flag for 'dbaas mongo cluster create' and 'dbaas mongo cluster update' commands (#532 - @fepape-ionos)

### Changed
- Changed default for 'dbaas mongo cluster create --version' to 7.0 (#532 - @fepape-ionos)
- Changed logic for 'delete --all' to ask for confirmation before deleting each resource

## [v6.8.5] (June 2025)

### Added
- Added support for 'dbaas in-memory-db' commands

### Others
- Removed mentions of deprecated flag in 'image upload' command help text (#529)


## [v6.8.4] (May 2025)

### Added
- Added support for ExposeSerial, RequireLegacyBios, ApplicationType properties to image commands

### Dependencies
- a minimum go version of 1.23.8 is now required to build ionosctl

## [v6.8.3] (May 2025)

### Added
- Added API Gateway support

### Fixed
- Fixed DBaaS Mongo basepath when using a custom API URL
- Fixed 'ionosctl token delete --current' command (#513)

## [v6.8.2] (April 2025)

### Added
- Added --server-type for managed kubernetes nodepools
- Added support for labeling images

### Changed
- Changed login message to include the username you have logged in as, if possible.

### Fixed
- Various fixes for 'container-registry' commands:
  - Added missing container-registry columns
  - The garbage collection days and time are now randomly generated to be a random working day during hours 10:00 - 16:00
  - Fixed listing container-registry token scopes
  - Fixed various panics related to nil ExpiryDate on container-registry tokens including #505

## [v6.8.1] (February 2025)

### Added

- Added parent resource ID to `list --all` command for: `container-registry artifacts`, `container-registry tokens`,
`kafka topic`

### Fixed

- Fixed not being able to use snapshot IDs for --image-id
- Fixed 404 when waiting for backupunits to be ready

### Changed

- When using config files for authentication, print explicit errors if the file cannot be read

## [v6.8.0] (January 2025)

### Added
- Added support for `vpn` commands
- Added support for `kafka` commands
- Added `--ttl` flag for creating tokens. You can specify the time-to-live for the token i.e. `--ttl 1h30m`
- Added a default location for Regional APIs when tab-completing

## [v6.7.9] (December 2024)

### Added
- Added support for CDN operations

### Fixed
- Removed Viper binding for `json-properties` and `json-properties-example` (fixes [463](https://github.com/ionos-cloud/ionosctl/issues/463))
- Fixed '--api-url' flag being overriden in certain cases
- Fixed MariaDB commands failing if a cluster is in the new "CREATING" state
- Regional APIs are now more robust (db mariadb, dns, cdn, logging service)
  - Missing locations to all of these APIs have been added
  - Added a '--location' flag which changes the '--api-url' accordingly
  - You can now force a specific location for these APIs even if it is not marked as supported by the `--location` flag
  - You cannot set both '--api-url' and '--location' at the same time

## [v6.7.8] (October 2024)

### Added
- Added support for `manpages` generation via `ionosctl man` command
- Added support for DNS `secondary-zones` and `zone files`

### Changed
- Added authentication warning to `image upload` help text

### Fixed
- Fixed examples for `image upload`

## [v6.7.7] (June 2024)

### Added
- Added support for DNS resources:
  - `ionosctl dns keys` commands which allows you to enable/disable DNSSEC and manage DNSKEY records.
  - `ionosctl dns quota` commands which allows you to get the DNS quota for your account.
  - `ionosctl dns reverse-record` commands which allows you to manage reverse DNS records.

### Fixed
- Fixed the column path mapping for 'server' resource to display the actual server's type ('CUBE'/'ENTERPRISE'),
  not the CloudAPI resource type ('server').
- Fixed a bug with 'nic create' not creating the LAN when missing and instead returning a 404
- Fixed json to table conversion error for `k8s nodepool lan list`
- Fixed error message on failed tokens API call by @avorima in #449

## [v6.7.6] (April 2024)

### Fixed
- Fixed a bug for `image upload` where using a custom `--ftp-url` and no `--location` would silently fail the operation

### Added
- Added support for 'dbaas mariadb'
- Added 'mg' alias for 'dbaas mongo'

### Deprecated
- Deprecated 'm' and 'mdb' aliases for 'dbaas mongo'. Using this alias will print a warning to stderr

## [v6.7.5] (March 2024)

### Added
- Added user and database commands to Postgres

### Changed
- Changed help text order to match terminal reading patterns: command-specific information moved lower, global/general help moved higher.
- Changed `--cpu-family` flag to always use the first valid CPU_FAMILY for the chosen location, previously it would always try using `AMD_OPTERON` family.

### Fixed
- Kubernetes node pool auto-scaling can be disabled by setting `--min-node-count` and `--max-node-count` to 0.
- Fixed missing version in `ionosctl version` output when installed via `go install`.
- Fixed a bug where explicit slices of strings would fail to be converted to slices of `interface{}` causing a panic for certain `--cols`

## [v6.7.4] (January 2024)

### Added
- Added `version` resource for Dataplatform API, with `list` and `latest` subcommands
- Added support for Container-Registry Vulnerabilities
  - New `--vulnerability-scanning` flag added to `registry create` and `registry update` commands
  - New `artifacts` and `vulnerabilities` commands under `container-registry`
  - `repository` command functionality will eventually be moved to `repository delete`. For the time being, both commands are available.

### Changed
- When creating a Dataplatform cluster, now the latest version will be used by default

### Fixed
- Fixed `--cols` for server: `server.cols` Viper variable being used both by `server` and `vm-autoscaling server` commands.

## [v6.7.3] (December 2023)

### Added
- Added support for private Kubernetes clusters
  - Use `--public=false` when creating a Kubernetes Cluster to use this feature
- Added support for VM Autoscaling API
- Added `shell` command for an interactive shell powered by [go-prompt](https://github.com/elk-language/go-prompt) via [comptplus](https://github.com/ionoscloudsdk/comptplus/) offering a new layer of interactivity and ease-of-use.
  - The shell is context-aware and will offer suggestions based on the current command.
  - This shell supports autocompletion for commands, flags, and flag values.
  - User input is currently unsupported, and commands with user input will fail and ask for `--force` to be set.

### Changed
- Improved help text, error handling & examples for `image upload`
- Deprecated `--image-alias` in favor of `--rename` for `image upload`
  - setting `--image-alias` will simply set `--rename`.

## [v6.7.2] (November 2023)

### Added
- Added support for Logging Service API
- Added `--json-properties` and `--json-properties-example` to `k8s nodepool create` which allows creation of nodepools using a JSON file. This is useful for creating nodepools with a large number of properties.
  - `--json-properties` is used to specify the path to the JSON file containing the nodepool properties.
  - `--json-properties-example` is used to generate a JSON file containing all the nodepool properties and their default values. This file can be used as a template for creating nodepools using JSON files.

### Fixed
- ionosctl will now exit with code 0 when no resources found for `image list`, `request list`.
- fix cluster k8sVersion column extraction JSON path by @printminion in #407
- Fixed `backupunit list` columns
- Fixed `backupunit get-sso-url` characters being treated as format placeholders
- Fixed various json paths (for certain columns extraction) for `user`, `location`, `nic`, `k8s cluster`, `dbaas postgres logs`

## [v6.7.1] (October 2023)

### Added

* Added `URN` column in `label` subcommands

### Fixed

* Fixed `ResourceId` and `ResourceType` columns not being printed in `label` subcommands
* Fixed `--no-headers` flag value being ignored - now treated as a global flag

### Changed

* Changed how `request targets` are printed for better readability
* In help text & documentation, `Private Cross Connect` has been renamed to `Cross Connect`, and an alias of `cc` has been added to the `pcc` command

## [v6.7.0] (October 2023)

### Added
* Added a new namespace, `cfg`, for commands related to the user's config file.
* Added the `cfg logout` convenience command which deletes sensitive data in the config file
* Added the `cfg location` convenience command which shows the location of the config file
* Added the `cfg whoami` command which allows debugging authentication:
  * If logged in (either with username & password or by JWT), it will print the current user's email
  * If `--provenance` is set, it will show which api-url is used, the used authentication layer, as well as if using a JWT or username & password.
  * A failed authentication will forcefully set `--provenance`.

* Added completer descriptions for `--image-id` flag.
  * Only relevant, usable images are now completed (e.g. HDD images for `volume create --image-id`).
  * Completed images will also be ordered so that private (user-uploaded) images are shown first.

* Added `api-json` output type, which affects `list --all` outputs, grouping children resources by their parent resource.

* Added IPv6 support for Datacenter, LAN, NIC and Firewall Rules.



### Changed
* Changed the `ionosctl login` logic for generating config files:
  * If a username & password is provided, it will now use these credentials to generate a token, which will be stored in the config file instead of the username & password pair.
  * If you are unable to use the IONOS API to generate a token, you can use a pre-generated one with `login --token <JWT>`
  * The default API URL `api.ionos.com` is no longer saved to the config file if the user doesn't provide any API URL.
  * If using `login --token <JWT>` to directly provide a JWT, it will be validated before being saved to the config file, however the user can set `--skip-verify` to skip this validation.
* Reworked the authentication logic to be layer-based.
  * The authentication layers, in order of priority, are:
      1. Global Flags
      2. Environment Variables
      3. Config File Entries
  * Within each layer, a token takes precedence over a username and password combination. For instance, if a token and a username/password pair are both defined in environment variables, ionosctl will prioritize the token. However, higher layers can override the use of a token from a lower layer. For example, username and password environment variables will supersede a token found in the config file.
* Moved `login` command under the new `cfg` namespace.
  * Note that all `cfg` namespace commands except `cfg location`, are also available as root-level commands (i.e. `ionosctl login`) for backwards-compatibility reasons, however they are hidden within the help text.
* Empty columns will now be removed from the output

### Fixed
* Fixed #249: Added `-o json` missing fields (e.g. `_links`, `type`, `href`, etc)
* Fixed #297: `ionosctl login` not clearing previous credentials
  * The command will ask for confirmation if a config file already exists at the set path. The user can skip this check by using `--force`
* Fixed 404 on firewallrule delete command: flag values not properly sent to API
* Fixed `password` or `sshkey-path` being required for private images
* Fixed 400 Bad Request by default on `dbaas mongo cluster create` due to `SSD` being the default storage type.


## [v6.6.10] (September 2023)
### Fixed
* Fixed #359: `image update` using unset flags.
* Fixed empty columns on `request list` by increasing the default request depth to 2.


## [v6.6.9] (September 2023)
### Fixed
* Fixed #349: IONOS_API_URL env var value being overriden by default values of flags

## [v6.6.8] (August 2023)
### Added
* Added new flags for `group create` and `group update` commands: `--access-dns`, `--manage-dbaas`, `--manage-registry`, `--manage-dataplatform`

### Fixed
* Changed default for `dataplatform cluster create --version` to 23.7 as 23.4 is no longer supported

## [v6.6.7] (August 2023)
### Added
* Added support for MongoDB Enterprise Edition (#340)
* Added support for completions helptext for Datacenter IDs, Lan IDs, Mongo resources (Templates, Clusters, Snapshots)
* Added support for resolving DBaaS Mongo Templates via a full word of their name (e.g. `--template playground`, `--template XS` is valid)
* Added default template for business edition: MongoDB Business XS (1 core, 50 GB storage, 2 GB RAM)
* Added inferred flag values to make `dbaas mongo` commands easier to use. For instance, setting `--shards` infers `--type sharded-cluster`, etc.
* Added defaults for DBaaS Mongo `--maintenance-day` and `--maintenance-window`: not setting these flags will result in a random day/time Monday-Friday between 10:00-16:00
* Added context-aware completions for `dbaas mongo user create --roles` to make this flag easier to use

## [v6.6.6] (August 2023)
### Added
* Added support for VCPU server type
* Added `token parse` command, which you can use to verify your token's privileges or see more details about your JWT

### Fixed
* Fixed #333: Fix flags --no-headers and --cols for `snapshot` command
* @avorima (#332): Avoid wrapping non-Result objects (i.e. arrays) in a `Message` JSON object when using `-o json`

## [v6.6.5] (July 2023)

### Fixed
* Changed default for `dataplatform cluster create --version` to 23.4 as 22.11 is no longer supported
* Fixed a bug for DNS completions regarding overriding the default value for `api-url`
* Fixed filters breaking for camelcase properties (e.g. `imageAlias`)
* Fixed missing filters for `RequestStatusMetadata` for command `request list` i.e. now you can also filter by `message, status`
* Removed the hardcoded `INTEL_SKYLAKE` value for `CPU_FAMILY` if creating a CUBE server. Now, by default for CUBE servers, this field is sent as nil to the API.

### Dependencies
* Various dependencies bumps (see #320). Most importantly:
  * Bump Cobra to v1.7.0
  * Bump Auth SDK to v1.0.6
  * Bump Container Registry SDK to v1.0.1, which fixes #326
  * Refactored away some dependencies (e.g. `google/uuid`, `uber/multierr`)

## [v6.6.4] (July 2023)

### Added
* Added support for DNS API
* Added RHEL license type
* Added the possibility of getting or deleting a token using the JWT directly: `--token`

### Fixed
* Deprecated warnings only show if the deprecated flag is being used

## [v6.6.3] (May 2023)

### Fixed
* Fixed token docs
* Fixed maintenance default, now maintenance is disabled by default for targetgroup target add
* Fix #288: improve client, config errors
* Fix #289  nodepool lan add --network flag using only last network

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
