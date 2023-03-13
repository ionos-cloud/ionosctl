package main

import (
	"fmt"
	"os"

	"github.com/ionos-cloud/ionosctl/v6/commands/container-registry"
	"github.com/ionos-cloud/ionosctl/v6/pkg/doc"
)

const (
	EnvVarDocsOutFolder = "DOCS_OUT_CONTAINER_REGISTRY"
)

func main() {
	dir := os.Getenv(EnvVarDocsOutFolder)
	if dir == "" {
		panic(fmt.Errorf("%s environment variable not set", EnvVarDocsOutFolder))
	}
	if _, err := os.Stat(dir); err != nil {
		panic(fmt.Errorf("error getting directory stat: %w", err))
	}

	err := doc.WriteDocs(container_registry.ContainerRegistryCmd(), dir)
	if err != nil {
		panic(fmt.Errorf("error writing docs %w", err))
	}
}
