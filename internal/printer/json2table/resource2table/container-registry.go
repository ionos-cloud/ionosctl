package resource2table

import (
	"fmt"

	"github.com/ionos-cloud/ionosctl/v6/internal/printer/json2table"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/json2table/jsonpaths"
	"github.com/ionos-cloud/sdk-go-bundle/products/containerregistry/v2"
)

func ConvertContainerRegistryVulnerabilitiesToTable(vulnerabilities containerregistry.ArtifactVulnerabilityReadList) (
	[]map[string]interface{}, error,
) {
	items, ok := vulnerabilities.GetItemsOk()
	if !ok || items == nil {
		return nil, fmt.Errorf("could not retrieve Container Registry Vulnerabilities items")
	}

	var convertedVulnerabilities []map[string]interface{}
	for _, vulnerability := range *items {
		convertedVulnerability, err := ConvertContainerRegistryVulnerabilityToTable(vulnerability)
		if err != nil {
			return nil, err
		}

		convertedVulnerabilities = append(convertedVulnerabilities, convertedVulnerability...)
	}

	return convertedVulnerabilities, nil
}

func ConvertContainerRegistryVulnerabilityToTable(vulnerability containerregistry.VulnerabilityRead) (
	[]map[string]interface{}, error,
) {
	properties, ok := vulnerability.GetPropertiesOk()
	if !ok || properties == nil {
		return nil, fmt.Errorf("could not retrieve Container Registry Vulnerability properties")
	}

	datasource, ok := properties.GetDataSourceOk()
	if !ok || datasource == nil {
		return nil, fmt.Errorf("could not retrieve Container Registry Vulnerability datasource")
	}

	affects, ok := properties.GetAffectsOk()
	if !ok || affects == nil {
		return nil, fmt.Errorf("could not retrieve Container Registry Vulnerability affects")
	}

	var affectsFormatted []interface{}
	for _, affect := range *affects {
		name, ok := affect.GetNameOk()
		if !ok || name == nil {
			return nil, fmt.Errorf("could not retrieve Container Registry Vulnerability affects name")
		}

		version, ok := affect.GetVersionOk()
		if !ok || version == nil {
			return nil, fmt.Errorf("could not retrieve Container Registry Vulnerability affects version")
		}

		affectsFormatted = append(affectsFormatted, fmt.Sprintf("%v (%v)", *name, *version))
	}

	convertedVulnerability, err := json2table.ConvertJSONToTable(
		"", jsonpaths.ContainerRegistryVulnerability,
		vulnerability,
	)
	if err != nil {
		return nil, err
	}

	datasourceId, idOk := datasource.GetIdOk()
	datasourceUrl, urlOk := datasource.GetUrlOk()

	if idOk && urlOk && datasourceId != nil && datasourceUrl != nil {
		convertedVulnerability[0]["DataSource"] = fmt.Sprintf("%v (%v)", *datasourceId, *datasourceUrl)
	}
	convertedVulnerability[0]["Affects"] = affectsFormatted

	return convertedVulnerability, nil
}
