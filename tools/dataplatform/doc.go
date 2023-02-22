package main

import (
	"fmt"
	"github.com/ionos-cloud/ionosctl/commands/dataplatform"
	"github.com/ionos-cloud/ionosctl/pkg/doc"
	"os"
)

const (
	EnvVarDocsOutFolder = "DOCS_OUT_DATAPLATFORM"
)

func main() {
	dir := os.Getenv(EnvVarDocsOutFolder)
	if dir == "" {
		panic(fmt.Errorf("%s environment variable not set", EnvVarDocsOutFolder))
	}
	if _, err := os.Stat(dir); err != nil {
		panic(fmt.Errorf("error getting directory stat: %w", err))
	}

	err := doc.WriteDocs(dataplatform.DataplatformCmd(), dir)
	if err != nil {
		panic(fmt.Errorf("error writing docs %w", err))
	}
}
