# Changelog

## \[5.0.5\]

* Added commands aliases
* Added flags aliases
* Renamed flags
* Improved `--cols` option for output

## \[5.0.4\]

* Updated sdk-go version to v5.1.0
* Added commands for IpFailover, IpConsumer, CD-ROM commands
* Added missing properties for resources \(e.g. `State`\)

## \[5.0.3\]

* Updated sdk-go to v5.0.3
* Fixed typo `K8sFindBySClusterId`

## \[5.0.2\]

* Added commands for Kubernetes, BackupUnit, Private Cross-Connect, Contract Resources, User Management
* Updated commands structure to: `ionosctl server volume attach`, `ionosctl loadbalancer nic attach`
* Updated documentation structure
* Added `--wait-for-request` and `--wait-for-state` options
* Renamed `--ignore-stdin` flag to `--force`

## \[5.0.1\]

* Added commands for image, snapshot, ip block, firewall rule, label
* Added support for token authentication
* Updated `attach` commands for volume and nic

## \[5.0.0\]

* Added commands for data center, server, volume, nic, lan, load balancer, request
* Added completion support for flags and commands for Zsh, Fish, PowerShell and Bash terminals
* Added login command for SDK authentication
* Added `ionosctl` boilerplate

