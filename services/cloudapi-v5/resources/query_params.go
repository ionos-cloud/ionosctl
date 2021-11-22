package resources

type ListQueryParams struct {
	Filters    *map[string]string `json:"Filters,omitempty"`
	OrderBy    *string            `json:"OrderBy,omitempty"`
	MaxResults *int32             `json:"MaxResults,omitempty"`
}

func (q ListQueryParams) SetFilters(filters map[string]string) ListQueryParams {
	q.Filters = &filters
	return q
}

func (q ListQueryParams) SetOrderBy(orderBy string) ListQueryParams {
	q.OrderBy = &orderBy
	return q
}

func (q ListQueryParams) SetMaxResults(maxResults int32) ListQueryParams {
	q.MaxResults = &maxResults
	return q
}

type QueryParams struct {
	Depth  *int32 `json:"Depth,omitempty"`
	Pretty *bool  `json:"Pretty,omitempty"`
}

func (d QueryParams) SetDepth(depth int32) QueryParams {
	d.Depth = &depth
	return d
}

func (d QueryParams) SetPretty(pretty bool) QueryParams {
	d.Pretty = &pretty
	return d
}
