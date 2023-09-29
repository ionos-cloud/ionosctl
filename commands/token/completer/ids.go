package completer

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/services/auth-v1/resources"
)

func TokensIds() []string {
	tokenSvc := resources.NewTokenService(client.Must(), context.Background())
	tokens, _, err := tokenSvc.List(0)
	if err != nil {
		return nil
	}
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
