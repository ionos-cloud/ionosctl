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

func TestSetFilters_CamelCasePreserved(t *testing.T) {
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

func TestSetFilters_UppercaseNormalizedToCamelCase(t *testing.T) {
	cfg := &mockSDKConfig{}
	filters := []string{"IMAGETYPE=gpu", "LOCATION=us"}
	setFilters(cfg, filters)

	// Uppercase keys should be normalized to correct camelCase
	if cfg.params["filter.imageType"] != "gpu" {
		t.Errorf("expected filter.imageType=gpu (from IMAGETYPE=gpu), got %v", cfg.params["filter.imageType"])
	}
	if cfg.params["filter.location"] != "us" {
		t.Errorf("expected filter.location=us (from LOCATION=us), got %v", cfg.params["filter.location"])
	}
}

func TestSetFilters_MixedCaseNormalizedAndMerged(t *testing.T) {
	cfg := &mockSDKConfig{}
	filters := []string{"Name=value1", "name=value2"}
	setFilters(cfg, filters)

	// Both should be normalized to "name" and merged
	expected := "value1,value2"
	if cfg.params["filter.name"] != expected {
		t.Errorf("expected filter.name=%s (merged from mixed case), got %v", expected, cfg.params["filter.name"])
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

func TestSetFilters_MultipleValuesPerKeyMixedCase(t *testing.T) {
	cfg := &mockSDKConfig{}
	filters := []string{"createdBy=user1", "CREATEDBY=user2"}
	setFilters(cfg, filters)

	// Both normalize to "createdBy" and merge
	expected := "user1,user2"
	if cfg.params["filter.createdBy"] != expected {
		t.Errorf("expected filter.createdBy=%s (merged mixed case), got %v", expected, cfg.params["filter.createdBy"])
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
	if len(cfg.params) != 2 {
		t.Errorf("expected 2 filters, got %d", len(cfg.params))
	}
}

func TestSetFilters_ComplexValue(t *testing.T) {
	cfg := &mockSDKConfig{}
	filters := []string{"name=gpu-datacenter-01"}
	setFilters(cfg, filters)

	if cfg.params["filter.name"] != "gpu-datacenter-01" {
		t.Errorf("expected filter.name=gpu-datacenter-01, got %v", cfg.params["filter.name"])
	}
}

func TestSetFilters_ValueWithEquals(t *testing.T) {
	cfg := &mockSDKConfig{}
	filters := []string{"filter=key=value"}
	setFilters(cfg, filters)

	if cfg.params["filter.filter"] != "key=value" {
		t.Errorf("expected filter.filter=key=value, got %v", cfg.params["filter.filter"])
	}
}

func TestSetFilters_CaseInsensitiveNormalization(t *testing.T) {
	tests := []struct {
		name     string
		filters  []string
		expected map[string]string
	}{
		{
			name:    "all uppercase normalized to camelCase",
			filters: []string{"NAME=test", "LOCATION=us"},
			expected: map[string]string{
				"filter.name":     "test",
				"filter.location": "us",
			},
		},
		{
			name:    "all lowercase preserved",
			filters: []string{"name=test", "location=us"},
			expected: map[string]string{
				"filter.name":     "test",
				"filter.location": "us",
			},
		},
		{
			name:    "mixed garbage case normalized to camelCase",
			filters: []string{"iMaGeTyPe=HDD", "LICENCETYPE=LINUX", "CreAtEdBY=admin"},
			expected: map[string]string{
				"filter.imageType":   "HDD",
				"filter.licenceType": "LINUX",
				"filter.createdBy":   "admin",
			},
		},
		{
			name:    "unknown key passed through as-is",
			filters: []string{"unknownFilter=value"},
			expected: map[string]string{
				"filter.unknownFilter": "value",
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
				t.Errorf("expected %d filters, got %d: %v", len(tt.expected), len(cfg.params), cfg.params)
			}
		})
	}
}

func TestNormalizeFilterKey(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"imageType", "imageType"},
		{"IMAGETYPE", "imageType"},
		{"imagetype", "imageType"},
		{"IMaGeTyPe", "imageType"},
		{"name", "name"},
		{"NAME", "name"},
		{"location", "location"},
		{"createdBy", "createdBy"},
		{"CREATEDBY", "createdBy"},
		// Unknown key passed through unchanged
		{"totallyUnknown", "totallyUnknown"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := normalizeFilterKey(tt.input)
			if result != tt.expected {
				t.Errorf("normalizeFilterKey(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}
