package client

import (
	"testing"
)

// mockSDKConfig implements sdkConfiguration for testing
type mockSDKConfig struct {
	params map[string]string
}

func (m *mockSDKConfig) AddDefaultQueryParam(key, val string) {
	if m.params == nil {
		m.params = make(map[string]string)
	}
	m.params[key] = val
}

func TestSetFilters_EmptyList(t *testing.T) {
	cfg := &mockSDKConfig{}
	setFilters(cfg, []string{})

	if len(cfg.params) > 0 {
		t.Errorf("expected no params for empty filter list, got %v", cfg.params)
	}
}

func TestSetFilters_LowercaseKeys(t *testing.T) {
	cfg := &mockSDKConfig{}
	filters := []string{"name=gpu", "location=us"}
	setFilters(cfg, filters)

	if cfg.params["filter.name"] != "gpu" {
		t.Errorf("expected filter.name=gpu, got %v", cfg.params["filter.name"])
	}
	if cfg.params["filter.location"] != "us" {
		t.Errorf("expected filter.location=us, got %v", cfg.params["filter.location"])
	}
}

func TestSetFilters_CamelCaseKeysPreserved(t *testing.T) {
	cfg := &mockSDKConfig{}
	filters := []string{"imageType=HDD", "imageAliases=ubuntu:latest"}
	setFilters(cfg, filters)

	if cfg.params["filter.imageType"] != "HDD" {
		t.Errorf("expected filter.imageType=HDD, got %v", cfg.params["filter.imageType"])
	}
	if cfg.params["filter.imageAliases"] != "ubuntu:latest" {
		t.Errorf("expected filter.imageAliases=ubuntu:latest, got %v", cfg.params["filter.imageAliases"])
	}
}

func TestSetFilters_MultipleValuesPerKey(t *testing.T) {
	cfg := &mockSDKConfig{}
	filters := []string{"status=active", "status=pending"}
	setFilters(cfg, filters)

	expected := "active,pending"
	if cfg.params["filter.status"] != expected {
		t.Errorf("expected filter.status=%s (comma-separated), got %v", expected, cfg.params["filter.status"])
	}
}

func TestSetFilters_DifferentCaseKeysTreatedSeparately(t *testing.T) {
	cfg := &mockSDKConfig{}
	filters := []string{"Name=value1", "name=value2"}
	setFilters(cfg, filters)

	// Different cases are now separate keys (API property names are case-sensitive)
	if cfg.params["filter.Name"] != "value1" {
		t.Errorf("expected filter.Name=value1, got %v", cfg.params["filter.Name"])
	}
	if cfg.params["filter.name"] != "value2" {
		t.Errorf("expected filter.name=value2, got %v", cfg.params["filter.name"])
	}
}

func TestSetFilters_SkipsInvalidFilters(t *testing.T) {
	cfg := &mockSDKConfig{}
	filters := []string{
		"valid=value",
		"invalid_no_equals",
		"=value_no_key",
		"name=gpu",
	}
	setFilters(cfg, filters)

	if cfg.params["filter.valid"] != "value" {
		t.Errorf("expected filter.valid=value, got %v", cfg.params["filter.valid"])
	}
	if cfg.params["filter.name"] != "gpu" {
		t.Errorf("expected filter.name=gpu, got %v", cfg.params["filter.name"])
	}
	if _, exists := cfg.params["filter."]; exists {
		t.Errorf("should not have created filter. for empty key")
	}
	// Should have exactly 2 filters
	if len(cfg.params) != 2 {
		t.Errorf("expected 2 filters, got %d", len(cfg.params))
	}
}

func TestSetFilters_ComplexValue(t *testing.T) {
	cfg := &mockSDKConfig{}
	filters := []string{"Name=gpu-datacenter-01"}
	setFilters(cfg, filters)

	if cfg.params["filter.Name"] != "gpu-datacenter-01" {
		t.Errorf("expected filter.Name=gpu-datacenter-01, got %v", cfg.params["filter.Name"])
	}
}

func TestSetFilters_ValueWithEquals(t *testing.T) {
	cfg := &mockSDKConfig{}
	filters := []string{"filter=key=value"}
	setFilters(cfg, filters)

	// SplitN with limit 2 should preserve the second "=" in the value
	if cfg.params["filter.filter"] != "key=value" {
		t.Errorf("expected filter.filter=key=value, got %v", cfg.params["filter.filter"])
	}
}

func TestSetFilters_KeyCasePreserved(t *testing.T) {
	tests := []struct {
		name     string
		filters  []string
		expected map[string]string
	}{
		{
			name:    "all lowercase",
			filters: []string{"name=test", "location=us"},
			expected: map[string]string{
				"filter.name":     "test",
				"filter.location": "us",
			},
		},
		{
			name:    "camelCase preserved",
			filters: []string{"imageType=HDD", "licenceType=LINUX", "imageAliases=ubuntu"},
			expected: map[string]string{
				"filter.imageType":    "HDD",
				"filter.licenceType":  "LINUX",
				"filter.imageAliases": "ubuntu",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &mockSDKConfig{}
			setFilters(cfg, tt.filters)

			for expectedKey, expectedVal := range tt.expected {
				if actual, exists := cfg.params[expectedKey]; !exists || actual != expectedVal {
					t.Errorf("expected %s=%s, got %v", expectedKey, expectedVal, actual)
				}
			}
			if len(cfg.params) != len(tt.expected) {
				t.Errorf("expected %d filters, got %d", len(tt.expected), len(cfg.params))
			}
		})
	}
}
