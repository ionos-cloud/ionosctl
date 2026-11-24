package version

import "testing"

func TestOptionalDate(t *testing.T) {
	format := optionalDate("properties.deprecationDate")
	cases := []struct {
		name string
		item map[string]any
		want any
	}{
		{"set date is kept", map[string]any{"properties": map[string]any{"deprecationDate": "2027-08-01T00:00:00Z"}}, "2027-08-01T00:00:00Z"},
		{"zero date is blank", map[string]any{"properties": map[string]any{"deprecationDate": zeroDate}}, ""},
		{"missing date is blank", map[string]any{"properties": map[string]any{}}, ""},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := format(tc.item); got != tc.want {
				t.Errorf("got %v, want %v", got, tc.want)
			}
		})
	}
}
