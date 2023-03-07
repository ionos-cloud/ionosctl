package main

import (
	"fmt"
	"github.com/ionos-cloud/ionosctl/v6/commands/dbaas/mongo"
	"github.com/ionos-cloud/ionosctl/v6/pkg/doc"
	"os"
)

const (
	EnvVarDocsOutFolder = "DOCS_OUT_DBAAS_MONGO"
)

func main() {
	dir := os.Getenv(EnvVarDocsOutFolder)
	if dir == "" {
		panic(fmt.Errorf("%s environment variable not set", EnvVarDocsOutFolder))
	}
	if _, err := os.Stat(dir); err != nil {
		panic(fmt.Errorf("error getting directory stat: %w", err))
	}

	err := doc.WriteDocs(mongo.DBaaSMongoCmd(), dir)
	if err != nil {
		panic(fmt.Errorf("error writing docs %w", err))
	}
}
