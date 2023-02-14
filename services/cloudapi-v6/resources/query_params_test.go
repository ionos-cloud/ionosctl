package resources

import (
	"github.com/ionos-cloud/ionosctl/internal/pointer"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestListQueryParams_SetMaxResults(t *testing.T) {
	tests := []struct {
		name string
		from ListQueryParams
		with int32
		want ListQueryParams
	}{
		{
			name: "Empty max results",
			from: ListQueryParams{},
			with: 0,
			want: ListQueryParams{},
		},
		{
			name: "Empty max results, depth set",
			from: ListQueryParams{QueryParams: QueryParams{Depth: pointer.From(int32(10))}},
			with: 0,
			want: ListQueryParams{QueryParams: QueryParams{Depth: pointer.From(int32(10))}},
		},
		{
			name: "Non empty max results",
			from: ListQueryParams{},
			with: 10,
			want: ListQueryParams{MaxResults: pointer.From(int32(10))},
		},
		{
			name: "Non empty max results, depth set",
			from: ListQueryParams{QueryParams: QueryParams{Depth: pointer.From(int32(10))}},
			with: 10,
			want: ListQueryParams{MaxResults: pointer.From(int32(10)), QueryParams: QueryParams{Depth: pointer.From(int32(10))}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q := tt.from.SetMaxResults(tt.with)
			assert.Equalf(t, tt.want, q, "SetMaxResults(%v)", tt.with)
		})
	}
}
