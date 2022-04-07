#!/bin/bash

# Regenerate mocks

# For Auth Resources
pushd "services/auth-v1/resources" >/dev/null || exit

GO111MODULE=off go get -d github.com/golang/mock/mockgen

mkdir -p mocks

mockgen -source client.go >mocks/ClientService.go
mockgen -source token.go >mocks/TokenService.go

pushd >/dev/null || exit
