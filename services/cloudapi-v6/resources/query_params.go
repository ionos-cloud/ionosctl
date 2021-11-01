package resources

type ListQueryParams struct {
	Filters            *map[string]string `json:"Filters,omitempty"`
	OrderBy            *string            `json:"OrderBy,omitempty"`
	Offset             *int32             `json:"Offset,omitempty"`
	Limit              *int32             `json:"Limit,omitempty"`
	DefaultQueryParams *QueryParams       `json:"DefaultQueryParams,omitempty"`
}

func (q ListQueryParams) SetFilters(filters map[string]string) ListQueryParams {
	q.Filters = &filters
	return q
}

func (q ListQueryParams) SetOrderBy(orderBy string) ListQueryParams {
	q.OrderBy = &orderBy
	return q
}

func (q ListQueryParams) SetOffset(offset int32) ListQueryParams {
	q.Offset = &offset
	return q
}

func (q ListQueryParams) SetLimit(limit int32) ListQueryParams {
	q.Limit = &limit
	return q
}

func (q ListQueryParams) SetDefaultQueryParams(params QueryParams) ListQueryParams {
	q.DefaultQueryParams = &params
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
