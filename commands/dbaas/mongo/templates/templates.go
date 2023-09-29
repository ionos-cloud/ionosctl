package templates

import (
	"context"
	"fmt"
	"strings"

	"github.com/gofrs/uuid/v5"
	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/functional"
	"github.com/ionos-cloud/ionosctl/v6/pkg/convbytes"
	"github.com/ionos-cloud/ionosctl/v6/pkg/core"
	"github.com/ionos-cloud/ionosctl/v6/pkg/json2table"
	ionoscloud "github.com/ionos-cloud/sdk-go-dbaas-mongo"
	"github.com/spf13/cobra"
)

func TemplatesCmd() *core.Command {
	cmd := &core.Command{
		Command: &cobra.Command{
			Use:              "templates",
			Aliases:          []string{"t"},
			Short:            "Mongo Templates Operations",
			Long:             "Templates can be used to create MongoDB clusters; they contain properties such as number of cores, RAM, and the storage size",
			TraverseChildren: true,
		},
	}

	cmd.AddCommand(TemplatesListCmd())
	return cmd
}

var (
	allJSONPaths = map[string]string{
		"TemplateId": "id",
		"Name":       "properties.name",
		"Edition":    "properties.edition",
		"Cores":      "properties.cores",
	}

	allCols = []string{"TemplateId", "Name", "Edition", "Cores", "StorageSize", "Ram"}
)

func convertTemplateToTable(template ionoscloud.TemplateResponse) ([]map[string]interface{}, error) {
	properties, ok := template.GetPropertiesOk()
	if !ok || properties == nil {
		return nil, fmt.Errorf("could not retrieve Mongo Template properties")
	}

	ram, ok := properties.GetRamOk()
	if !ok || ram == nil {
		return nil, fmt.Errorf("could not retrieve Mongo Template RAM")
	}

	storage, ok := properties.GetStorageSizeOk()
	if !ok || storage == nil {
		return nil, fmt.Errorf("could not retrieve Mongo Template storage")
	}

	temp, err := json2table.ConvertJSONToTable("", allJSONPaths, template)
	if err != nil {
		return nil, fmt.Errorf("could not convert from JSON to Table format: %w", err)
	}

	temp[0]["RAM"] = fmt.Sprintf("%d GB", convbytes.Convert(int64(*ram), convbytes.MB, convbytes.GB))
	temp[0]["StorageSize"] = fmt.Sprintf("%d GB", convbytes.Convert(int64(*storage), convbytes.MB, convbytes.GB))

	return temp, nil
}

func convertTemplatesToTable(templates ionoscloud.TemplateList) ([]map[string]interface{}, error) {
	items, ok := templates.GetItemsOk()
	if !ok || items == nil {
		return nil, fmt.Errorf("could not retrieve Templates items")
	}

	var templatesConverted []map[string]interface{}
	for _, item := range *items {
		temp, err := convertTemplateToTable(item)
		if err != nil {
			return nil, err
		}

		templatesConverted = append(templatesConverted, temp...)
	}

	return templatesConverted, nil
}

// List retrieves a list of templates, optionally filtered by a given funcs
func List(filters ...func(x ionoscloud.TemplateResponse) bool) ([]ionoscloud.TemplateResponse, error) {
	xs, _, err := client.Must().MongoClient.TemplatesApi.TemplatesGet(context.Background()).Execute()
	if err != nil {
		return nil, fmt.Errorf("failed getting templates: %w", err)
	}

	if len(filters) == 0 {
		return *xs.GetItems(), nil
	}

	filteredTemplates := *xs.GetItems()
	for _, f := range filters {
		filteredTemplates = functional.Filter(filteredTemplates, f)
	}

	return filteredTemplates, nil
}

// Find returns the first template for which found() returns true
func Find(found func(x ionoscloud.TemplateResponse) bool) (ionoscloud.TemplateResponse, error) {
	filteredTemplates, err := List(found)
	if err != nil {
		return ionoscloud.TemplateResponse{}, err
	}

	if len(filteredTemplates) > 0 {
		return filteredTemplates[0], nil
	}
	return ionoscloud.TemplateResponse{}, fmt.Errorf("no matching template found")
}

// Resolve resolves nameOrId to the ID of the template.
// If it's an ID, it's returned as is. If it's not, then it's a name, and we try to resolve it
// with a case sensitive "whole word match" operation for the name of the template
//
// e.g.:
// - Resolve("S") -> "id of MongoDB Business S template" (note that 4XL_S is correctly ignored in this case)
// - Resolve("id of MongoDB Business L template") -> "id of MongoDB Business L template"
func Resolve(nameOrId string) (string, error) {
	if dumbWord := strings.ToLower(nameOrId); dumbWord == "mongodb" || dumbWord == "business" {
		// Save the user from himself by throwing away queries that result in very vague and unwanted expensive stuff
		return "", fmt.Errorf("the words used to select your template (%s) are too vague. please be more specific", dumbWord)
	}

	uid, errParseUuid := uuid.FromString(nameOrId)
	id := uid.String()
	if errParseUuid != nil {
		// It's a name

		// Why doesn't the API have a FindByID or something :(
		templateMatchingWholeWordIgnoreCase, err := Find(func(x ionoscloud.TemplateResponse) bool {
			if x.Properties == nil || x.Properties.Name == nil {
				return false
			}

			words := strings.Split(*x.Properties.Name, " ")
			for _, word := range words {
				if strings.ToLower(word) == strings.ToLower(nameOrId) {
					return true
				}
			}
			return false
		})

		if err != nil {
			return "", fmt.Errorf("failed finding a template with any word of its name case-insensitively matching %s: %w", nameOrId, err)
		}

		id = *templateMatchingWholeWordIgnoreCase.Id
	}
	return id, nil
}
