package completer

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/internal/client"
)

func TokensIds() []string {
	tokens, _, err := client.Must().AuthClient.TokensApi.TokensGet(context.Background()).Execute()
	if err != nil {
		return nil
	}
	tokenIds := make([]string, 0)
	if items, ok := tokens.GetTokensOk(); ok && items != nil {
		for _, item := range items {
			if itemId, ok := item.GetIdOk(); ok && itemId != nil {
				tokenIds = append(tokenIds, *itemId)
			}
		}
	} else {
		return nil
	}
	return tokenIds
}
