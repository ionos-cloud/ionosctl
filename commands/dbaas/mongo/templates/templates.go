package templates

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/fatih/structs"
	"github.com/gofrs/uuid/v5"
	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/functional"
	"github.com/ionos-cloud/ionosctl/v6/pkg/constants"
	"github.com/ionos-cloud/ionosctl/v6/pkg/core"
	"github.com/ionos-cloud/ionosctl/v6/pkg/printer"
	"github.com/ionos-cloud/ionosctl/v6/pkg/utils"
	ionoscloud "github.com/ionos-cloud/sdk-go-dbaas-mongo"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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

func getTemplatesPrint(c *core.CommandConfig, ls *[]ionoscloud.TemplateResponse) printer.Result {
	r := printer.Result{}
	if c != nil && ls != nil {
		r.OutputJSON = ls
		r.KeyValue = getTemplateRows(ls)                                                                                   // map header -> rows
		r.Columns = printer.GetHeadersAllDefault(allCols, viper.GetStringSlice(core.GetFlagName(c.NS, constants.ArgCols))) // headers
	}
	return r
}

type TemplatePrint struct {
	TemplateId  string `json:"TemplateId,omitempty"`
	Name        string `json:"Name,omitempty"`
	Edition     string `json:"Edition,omitempty"`
	Cores       int32  `json:"Cores,omitempty"`
	StorageSize string `json:"StorageSize,omitempty"`
	Ram         string `json:"Ram,omitempty"`
}

var allCols = structs.Names(TemplatePrint{})

func getTemplateRows(ls *[]ionoscloud.TemplateResponse) []map[string]interface{} {
	if ls == nil {
		return nil
	}

	out := make([]map[string]interface{}, 0, len(*ls))
	for _, t := range *ls {
		var cols TemplatePrint

		if t.Id != nil {
			cols.TemplateId = *t.Id
		}

		properties := t.Properties
		if properties != nil {
			if properties.Cores != nil {
				cols.Cores = *properties.Cores
			}
			if properties.StorageSize != nil {
				cols.StorageSize = fmt.Sprintf("%d GB", *properties.StorageSize)
			}
			if properties.Name != nil {
				cols.Name = *properties.Name
			}
			if properties.Edition != nil {
				cols.Edition = *properties.Edition
			}
			if properties.Ram != nil {
				ramGb, err := utils.ConvertToGB(strconv.Itoa(int(*properties.Ram)), utils.MegaBytes)
				if err == nil {
					cols.Ram = fmt.Sprintf("%d GB", ramGb)
				}
			}
		}

		o := structs.Map(cols)
		out = append(out, o)
	}
	return out
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
func Resolve(nameOrId string) string {
	if dumbWord := strings.ToLower(nameOrId); dumbWord == "mongodb" || dumbWord == "business" {
		// Save the user from himself by throwing away queries that result in very vague and unwanted expensive stuff
		return ""
	}

	uid, errParseUuid := uuid.FromString(nameOrId)
	id := uid.String()
	if errParseUuid != nil {
		// It's a name
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
			return ""
		}

		id = *templateMatchingWholeWordIgnoreCase.Id
	}
	return id
}
