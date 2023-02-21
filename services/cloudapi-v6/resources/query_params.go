package resources

type ListQueryParams struct {
	Filters     *map[string][]string `json:"Filters,omitempty"`
	OrderBy     *string              `json:"OrderBy,omitempty"`
	MaxResults  *int32               `json:"MaxResults,omitempty"`
	QueryParams QueryParams          `json:"QueryParams,omitempty"`
}

func (q ListQueryParams) SetFilters(filters map[string][]string) ListQueryParams {
	q.Filters = &filters
	return q
}

func (q ListQueryParams) SetOrderBy(orderBy string) ListQueryParams {
	if orderBy != "" {
		q.OrderBy = &orderBy
	}
	return q
}

func (q ListQueryParams) SetMaxResults(maxResults int32) ListQueryParams {
	if maxResults > 0 {
		q.MaxResults = &maxResults
	}
	return q
}

func (q ListQueryParams) SetDepth(depth int32) ListQueryParams {
	q.QueryParams.Depth = &depth
	return q
}

func (q ListQueryParams) SetPretty(pretty bool) ListQueryParams {
	q.QueryParams.Pretty = &pretty
	return q
}

// TODO: Merge ListQueryParams into QueryParams
// TODO: Once Compute namespace is added, add all of QueryParam's flags as Global flags, to reduce duplication

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
