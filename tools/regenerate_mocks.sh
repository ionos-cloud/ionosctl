#!/bin/bash

# Regenerate mocks

pushd "pkg/printer" >/dev/null || exit

GO111MODULE=off go get -d github.com/golang/mock/mockgen

mockgen -source printer.go >mocks/PrintService.go

pushd >/dev/null || exit
