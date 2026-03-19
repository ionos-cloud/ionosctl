package filterprops

import "testing"

func TestFirstLower(t *testing.T) {
	tests := []struct {
		in, want string
	}{
		{"ImageType", "imageType"},
		{"Name", "name"},
		{"A", "a"},
		{"", ""},
		{"already", "already"},
		{"ID", "iD"},
	}
	for _, tt := range tests {
		if got := FirstLower(tt.in); got != tt.want {
			t.Errorf("FirstLower(%q) = %q, want %q", tt.in, got, tt.want)
		}
	}
}

func TestAllCloudAPIv6Types_NotEmpty(t *testing.T) {
	types := AllCloudAPIv6Types()
	if len(types) == 0 {
		t.Fatal("AllCloudAPIv6Types() returned empty slice")
	}
}
