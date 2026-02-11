package resource2table

import (
	"fmt"

	"github.com/ionos-cloud/ionosctl/v6/internal/printer/json2table"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/json2table/jsonpaths"
	psql "github.com/ionos-cloud/sdk-go-dbaas-psql"
)

func ConvertDbaasPostgresVersionsToTable(versions psql.PostgresVersionReadList) ([]map[string]interface{}, error) {
	items, ok := versions.GetItemsOk()
	if !ok || items == nil {
		return nil, nil
	}

	var versionsConverted []map[string]interface{}
	for _, item := range items {
		temp, err := ConvertDbaasPostgresVersionToTable(item)
		if err != nil {
			return nil, err
		}

		versionsConverted = append(versionsConverted, temp...)
	}

	return versionsConverted, nil
}

func ConvertDbaasPostgresVersionToTable(version psql.PostgresVersionRead) ([]map[string]interface{}, error) {
	table, err := json2table.ConvertJSONToTable("", jsonpaths.DbaasPostgresVersion, version)
	if err != nil {
		return nil, fmt.Errorf("failed getting table representation of version: %w", err)
	}

	return table, nil
}
