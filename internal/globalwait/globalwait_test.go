package globalwait

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestExtractHref(t *testing.T) {
	tests := []struct {
		name     string
		input    any
		expected string
	}{
		{
			name:     "single resource with href",
			input:    map[string]any{"href": "https://api.ionos.com/cloudapi/v6/datacenters/abc", "id": "abc"},
			expected: "https://api.ionos.com/cloudapi/v6/datacenters/abc",
		},
		{
			name:     "list response with items - skip",
			input:    map[string]any{"href": "https://api.ionos.com/cloudapi/v6/datacenters", "items": []any{}},
			expected: "",
		},
		{
			name:     "no href field",
			input:    map[string]any{"id": "abc", "name": "test"},
			expected: "",
		},
		{
			name:     "nil input",
			input:    nil,
			expected: "",
		},
		{
			name:     "non-map input",
			input:    "just a string",
			expected: "",
		},
		{
			name:     "relative href",
			input:    map[string]any{"href": "/cloudapi/v6/datacenters/abc"},
			expected: "/cloudapi/v6/datacenters/abc",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ExtractHref(tt.input)
			assert.Equal(t, tt.expected, got)
		})
	}
}

func TestCaptureAndGetHref(t *testing.T) {
	Reset()
	assert.Empty(t, GetHref())

	CaptureHref("https://api.ionos.com/test")
	assert.Equal(t, "https://api.ionos.com/test", GetHref())
}

func TestCaptureAndGetRerenderable(t *testing.T) {
	Reset()
	r, cols := GetRerenderable()
	assert.Nil(t, r)
	assert.Nil(t, cols)

	mock := &mockRerenderable{}
	CaptureRerenderable(mock, []string{"Col1", "Col2"})

	r, cols = GetRerenderable()
	assert.Equal(t, mock, r)
	assert.Equal(t, []string{"Col1", "Col2"}, cols)
}

func TestIsRerendering(t *testing.T) {
	Reset()
	assert.False(t, IsRerendering())

	SetRerendering(true)
	assert.True(t, IsRerendering())

	SetRerendering(false)
	assert.False(t, IsRerendering())
}

func TestReset(t *testing.T) {
	CaptureHref("test-href")
	CaptureRerenderable(&mockRerenderable{}, []string{"col"})
	SetRerendering(true)

	Reset()

	assert.Empty(t, GetHref())
	r, cols := GetRerenderable()
	assert.Nil(t, r)
	assert.Nil(t, cols)
	assert.False(t, IsRerendering())
}

func TestPoll_Available(t *testing.T) {
	callCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callCount++
		state := "BUSY"
		if callCount >= 2 {
			state = "AVAILABLE"
		}
		json.NewEncoder(w).Encode(map[string]any{
			"metadata": map[string]any{"state": state},
		})
	}))
	defer server.Close()

	// Speed up for testing
	oldInterval := pollInterval
	pollInterval = 100 * time.Millisecond
	defer func() { pollInterval = oldInterval }()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := Poll(ctx, server.URL, "", "", "")
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, callCount, 2)
}

func TestPoll_Failed(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]any{
			"metadata": map[string]any{"state": "FAILED"},
		})
	}))
	defer server.Close()

	oldInterval := pollInterval
	pollInterval = 100 * time.Millisecond
	defer func() { pollInterval = oldInterval }()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := Poll(ctx, server.URL, "", "", "")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "FAILED")
}

func TestPoll_Timeout(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]any{
			"metadata": map[string]any{"state": "BUSY"},
		})
	}))
	defer server.Close()

	oldInterval := pollInterval
	pollInterval = 100 * time.Millisecond
	defer func() { pollInterval = oldInterval }()

	ctx, cancel := context.WithTimeout(context.Background(), 350*time.Millisecond)
	defer cancel()

	err := Poll(ctx, server.URL, "", "", "")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "timeout")
}

func TestPoll_StatusField(t *testing.T) {
	// VPN uses "status" instead of "state"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]any{
			"metadata": map[string]any{"status": "ACTIVE"},
		})
	}))
	defer server.Close()

	oldInterval := pollInterval
	pollInterval = 100 * time.Millisecond
	defer func() { pollInterval = oldInterval }()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := Poll(ctx, server.URL, "", "", "")
	assert.NoError(t, err)
}

func TestPoll_AuthHeaders(t *testing.T) {
	var gotAuth string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotAuth = r.Header.Get("Authorization")
		json.NewEncoder(w).Encode(map[string]any{
			"metadata": map[string]any{"state": "AVAILABLE"},
		})
	}))
	defer server.Close()

	oldInterval := pollInterval
	pollInterval = 100 * time.Millisecond
	defer func() { pollInterval = oldInterval }()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := Poll(ctx, server.URL, "test-token", "", "")
	assert.NoError(t, err)
	assert.Equal(t, "Bearer test-token", gotAuth)
}

func TestBuildFullURL(t *testing.T) {
	tests := []struct {
		name     string
		href     string
		expected string
	}{
		{
			name:     "full URL",
			href:     "https://api.ionos.com/cloudapi/v6/datacenters/abc",
			expected: "https://api.ionos.com/cloudapi/v6/datacenters/abc?depth=1",
		},
		{
			name:     "full URL with query params",
			href:     "https://api.ionos.com/cloudapi/v6/datacenters/abc?pretty=true",
			expected: "https://api.ionos.com/cloudapi/v6/datacenters/abc?pretty=true&depth=1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := buildFullURL(tt.href)
			assert.Equal(t, tt.expected, got)
		})
	}
}

func TestPollWithJSONLog(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]any{
			"metadata": map[string]any{"state": "AVAILABLE"},
		})
	}))
	defer server.Close()

	oldInterval := pollInterval
	pollInterval = 100 * time.Millisecond
	defer func() { pollInterval = oldInterval }()

	var buf bytes.Buffer
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := pollWithJSONLog(ctx, &buf, server.URL, "", "", "")
	assert.NoError(t, err)
	assert.Contains(t, buf.String(), "Waiting for state")
	assert.Contains(t, buf.String(), "DONE")
}

// mockRerenderable implements Rerenderable for testing
type mockRerenderable struct {
	extractCalled bool
	renderCalled  bool
	data          any
}

func (m *mockRerenderable) Extract(sourceData any) error {
	m.extractCalled = true
	m.data = sourceData
	return nil
}

func (m *mockRerenderable) Render(visibleCols []string) (string, error) {
	m.renderCalled = true
	return fmt.Sprintf("rendered:%v", m.data), nil
}

func TestParentHref(t *testing.T) {
	tests := []struct {
		name     string
		href     string
		expected string
	}{
		{
			name:     "volume to server (CloudAPI)",
			href:     "https://api.ionos.com/cloudapi/v6/datacenters/aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee/servers/11111111-2222-3333-4444-555555555555/volumes/66666666-7777-8888-9999-000000000000",
			expected: "https://api.ionos.com/cloudapi/v6/datacenters/aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee/servers/11111111-2222-3333-4444-555555555555",
		},
		{
			name:     "server to datacenter (CloudAPI)",
			href:     "https://api.ionos.com/cloudapi/v6/datacenters/aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee/servers/11111111-2222-3333-4444-555555555555",
			expected: "https://api.ionos.com/cloudapi/v6/datacenters/aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee",
		},
		{
			name:     "datacenter stop (CloudAPI root)",
			href:     "https://api.ionos.com/cloudapi/v6/datacenters/aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee",
			expected: "",
		},
		{
			name:     "VPN peer to gateway (regional)",
			href:     "https://vpn.de-fra.ionos.com/wireguardgateways/aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee/peers/11111111-2222-3333-4444-555555555555",
			expected: "https://vpn.de-fra.ionos.com/wireguardgateways/aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee",
		},
		{
			name:     "VPN gateway stop (regional root)",
			href:     "https://vpn.de-fra.ionos.com/wireguardgateways/aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee",
			expected: "",
		},
		{
			name:     "too short",
			href:     "https://api.ionos.com/cloudapi/v6",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parentHref(tt.href)
			assert.Equal(t, tt.expected, got)
		})
	}
}

func TestResourceAndParentURLs(t *testing.T) {
	urls := resourceAndParentURLs("https://api.ionos.com/cloudapi/v6/datacenters/aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee/servers/11111111-2222-3333-4444-555555555555/volumes/66666666-7777-8888-9999-000000000000")
	assert.Equal(t, []string{
		"https://api.ionos.com/cloudapi/v6/datacenters/aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee/servers/11111111-2222-3333-4444-555555555555/volumes/66666666-7777-8888-9999-000000000000",
		"https://api.ionos.com/cloudapi/v6/datacenters/aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee/servers/11111111-2222-3333-4444-555555555555",
		"https://api.ionos.com/cloudapi/v6/datacenters/aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee",
	}, urls)
}

func TestResourceAndParentURLs_TopLevel(t *testing.T) {
	urls := resourceAndParentURLs("https://api.ionos.com/cloudapi/v6/datacenters/aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee")
	assert.Equal(t, []string{
		"https://api.ionos.com/cloudapi/v6/datacenters/aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee",
	}, urls)
}

func TestResourceAndParentURLs_Regional(t *testing.T) {
	urls := resourceAndParentURLs("https://vpn.de-fra.ionos.com/wireguardgateways/aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee/peers/11111111-2222-3333-4444-555555555555")
	assert.Equal(t, []string{
		"https://vpn.de-fra.ionos.com/wireguardgateways/aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee/peers/11111111-2222-3333-4444-555555555555",
		"https://vpn.de-fra.ionos.com/wireguardgateways/aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee",
	}, urls)
}

func TestCaptureRequestURL(t *testing.T) {
	Reset()

	// CaptureRequestURL sets href when none captured from table output
	CaptureRequestURL("https://api.ionos.com/cloudapi/v6/datacenters/aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee/servers/11111111-2222-3333-4444-555555555555/volumes/66666666-7777-8888-9999-000000000000")
	assert.Equal(t, "https://api.ionos.com/cloudapi/v6/datacenters/aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee/servers/11111111-2222-3333-4444-555555555555/volumes/66666666-7777-8888-9999-000000000000", GetHref())

	// If href already set (from table output), CaptureRequestURL doesn't overwrite
	Reset()
	CaptureHref("https://api.ionos.com/cloudapi/v6/datacenters/aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee")
	CaptureRequestURL("https://api.ionos.com/cloudapi/v6/datacenters/aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee/servers/11111111-2222-3333-4444-555555555555")
	assert.Equal(t, "https://api.ionos.com/cloudapi/v6/datacenters/aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee", GetHref())
}
