package completer

import (
	"context"
	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/pkg/functional"
	cdn "github.com/ionos-cloud/sdk-go-cdn"
)

// DistributionsProperty returns a list of properties of all distributions matching the given filters
func DistributionsProperty[V any](f func(cdn.Distribution) V, fs ...Filter) []V {
	recs, err := Distributions(fs...)
	if err != nil {
		return nil
	}
	return functional.Map(*recs.Items, f)
}

// Distributions returns all distributions matching the given filters
func Distributions(fs ...Filter) (cdn.Distributions, error) {
	req := client.Must().CDNClient.DistributionsApi.DistributionsGet(context.Background())
	for _, f := range fs {
		var err error
		req, err = f(req)
		if err != nil {
			return cdn.Distributions{}, err
		}
	}

	ls, _, err := req.Execute()
	if err != nil {
		return cdn.Distributions{}, err
	}
	return ls, nil
}

type Filter func(request cdn.ApiDistributionsGetRequest) (cdn.ApiDistributionsGetRequest, error)
