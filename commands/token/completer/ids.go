package completer

import (
	"context"
	"io"

	client2 "github.com/ionos-cloud/ionosctl/v6/internal/client"

	"github.com/ionos-cloud/ionosctl/v6/pkg/utils/clierror"
	"github.com/ionos-cloud/ionosctl/v6/services/auth-v1/resources"
)

func TokensIds(outErr io.Writer) []string {
	client, err := client2.Get()
	clierror.CheckErrorAndDie(err, outErr)
	tokenSvc := resources.NewTokenService(client, context.TODO())
	tokens, _, err := tokenSvc.List(0)
	clierror.CheckErrorAndDie(err, outErr)
	tokenIds := make([]string, 0)
	if items, ok := tokens.Tokens.GetTokensOk(); ok && items != nil {
		for _, item := range *items {
			if itemId, ok := item.GetIdOk(); ok && itemId != nil {
				tokenIds = append(tokenIds, *itemId)
			}
		}
	} else {
		return nil
	}
	return tokenIds
}
