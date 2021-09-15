#!/bin/bash

# Regenerate mocks

pushd "internal/printer" >/dev/null || exit

GO111MODULE=off go get -u github.com/golang/mock/mockgen

mockgen -source printer.go >mocks/PrintService.go

pushd >/dev/null || exit
