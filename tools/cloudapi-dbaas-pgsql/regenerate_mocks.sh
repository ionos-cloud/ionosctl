#!/bin/bash

# Regenerate Mocks

# For CloudApi DbaasPgsql Resources
pushd "services/cloudapi-dbaas-pgsql/resources" >/dev/null || exit

GO111MODULE=off go get -u github.com/golang/mock/mockgen

mockgen -source client.go >mocks/ClientService.go
mockgen -source backup.go >mocks/BackupService.go
mockgen -source cluster.go >mocks/ClusterService.go
mockgen -source version.go >mocks/VersionService.go
mockgen -source info.go >mocks/InfoService.go
mockgen -source logs.go >mocks/LogService.go
mockgen -source quota.go >mocks/QuotaService.go
mockgen -source restore.go >mocks/RestoreService.go

pushd >/dev/null || exit
