#!/bin/bash

# Regenerate Mocks

# For Data Platform Resources
pushd "services/dataplatform/resources" >/dev/null || exit

GO111MODULE=off go get -d github.com/golang/mock/mockgen

mockgen -source client.go > mocks/ClientService.go
mockgen -source cluster.go > mocks/ClusterService.go


pushd >/dev/null || exit
