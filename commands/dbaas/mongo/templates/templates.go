package templates

import (
	"context"
	"fmt"
	"strings"

	"github.com/gofrs/uuid/v5"
	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/pkg/functional"
	mongo "github.com/ionos-cloud/sdk-go-bundle/products/dbaas/mongo/v2"
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
	allCols = []string{"TemplateId", "Name", "Edition", "Cores", "StorageSize", "Ram"}
)

// List retrieves a list of templates, optionally filtered by a given funcs
func List(filters ...func(x mongo.TemplateResponse) bool) ([]mongo.TemplateResponse, error) {
	xs, _, err := client.Must().MongoClient.TemplatesApi.TemplatesGet(context.Background()).Execute()
	if err != nil {
		return nil, fmt.Errorf("failed getting templates: %w", err)
	}

	if len(filters) == 0 {
		return xs.GetItems(), nil
	}

	filteredTemplates := xs.GetItems()
	for _, f := range filters {
		filteredTemplates = functional.Filter(filteredTemplates, f)
	}

	return filteredTemplates, nil
}

// Find returns the first template for which found() returns true
func Find(found func(x mongo.TemplateResponse) bool) (mongo.TemplateResponse, error) {
	filteredTemplates, err := List(found)
	if err != nil {
		return mongo.TemplateResponse{}, err
	}

	if len(filteredTemplates) > 0 {
		return filteredTemplates[0], nil
	}
	return mongo.TemplateResponse{}, fmt.Errorf("no matching template found")
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
		templateMatchingWholeWordIgnoreCase, err := Find(func(x mongo.TemplateResponse) bool {
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
