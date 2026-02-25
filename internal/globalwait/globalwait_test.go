package globalwait

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func useFastPolling(t *testing.T) {
	t.Helper()
	orig := pollInterval
	pollInterval = 50 * time.Millisecond
	t.Cleanup(func() { pollInterval = orig })
}

func TestExtractHref_SingleResource(t *testing.T) {
	data := map[string]any{
		"id":   "test-id",
		"type": "datacenter",
		"href": "/cloudapi/v6/datacenters/test-id",
		"metadata": map[string]any{
			"state": "BUSY",
		},
		"properties": map[string]any{
			"name": "test",
		},
	}
	assert.Equal(t, "/cloudapi/v6/datacenters/test-id", ExtractHref(data))
}

func TestExtractHref_ListResponse(t *testing.T) {
	data := map[string]any{
		"id":   "datacenters",
		"type": "collection",
		"href": "/cloudapi/v6/datacenters",
		"items": []any{
			map[string]any{
				"id":   "dc-1",
				"href": "/cloudapi/v6/datacenters/dc-1",
			},
		},
	}
	// Should return empty - list responses are skipped
	assert.Equal(t, "", ExtractHref(data))
}

func TestExtractHref_NoHref(t *testing.T) {
	data := map[string]any{
		"message": "Resource deleted",
	}
	assert.Equal(t, "", ExtractHref(data))
}

func TestExtractHref_InvalidData(t *testing.T) {
	assert.Equal(t, "", ExtractHref(nil))
	assert.Equal(t, "", ExtractHref("not a json object"))
	assert.Equal(t, "", ExtractHref(42))
}

func TestExtractHref_FullURL(t *testing.T) {
	data := map[string]any{
		"id":   "zone-1",
		"href": "https://dns.de-fra.ionos.com/zones/zone-1",
	}
	assert.Equal(t, "https://dns.de-fra.ionos.com/zones/zone-1", ExtractHref(data))
}

func TestCaptureAndGetHref(t *testing.T) {
	Reset()

	data := map[string]any{
		"id":   "test-id",
		"href": "/cloudapi/v6/datacenters/test-id",
	}
	CaptureHref(data)
	assert.Equal(t, "/cloudapi/v6/datacenters/test-id", GetHref())
}

func TestCaptureHref_OverwritesPrevious(t *testing.T) {
	Reset()

	data1 := map[string]any{"id": "1", "href": "/first"}
	data2 := map[string]any{"id": "2", "href": "/second"}

	CaptureHref(data1)
	CaptureHref(data2)
	assert.Equal(t, "/second", GetHref())
}

func TestCaptureHref_SkipsEmptyHref(t *testing.T) {
	Reset()

	data1 := map[string]any{"id": "1", "href": "/first"}
	data2 := map[string]any{"message": "no href here"}

	CaptureHref(data1)
	CaptureHref(data2)
	// Should still have the first href since data2 has no href
	assert.Equal(t, "/first", GetHref())
}

func TestReset(t *testing.T) {
	data := map[string]any{"id": "1", "href": "/test"}
	CaptureHref(data)
	assert.NotEmpty(t, GetHref())

	Reset()
	assert.Empty(t, GetHref())
}

func TestPoll_AvailableImmediately(t *testing.T) {
	useFastPolling(t)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]any{
			"metadata": map[string]any{"state": "AVAILABLE"},
		})
	}))
	defer server.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := Poll(ctx, server.URL, "test-token", "", "")
	assert.NoError(t, err)
}

func TestPoll_TransitionsToAvailable(t *testing.T) {
	useFastPolling(t)

	var callCount atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		count := callCount.Add(1)
		state := "BUSY"
		if count >= 3 {
			state = "AVAILABLE"
		}
		json.NewEncoder(w).Encode(map[string]any{
			"metadata": map[string]any{"state": state},
		})
	}))
	defer server.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := Poll(ctx, server.URL, "test-token", "", "")
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, callCount.Load(), int32(3))
}

func TestPoll_Failed(t *testing.T) {
	useFastPolling(t)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]any{
			"metadata": map[string]any{"state": "FAILED"},
		})
	}))
	defer server.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := Poll(ctx, server.URL, "test-token", "", "")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "FAILED")
}

func TestPoll_Timeout(t *testing.T) {
	useFastPolling(t)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]any{
			"metadata": map[string]any{"state": "BUSY"},
		})
	}))
	defer server.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	err := Poll(ctx, server.URL, "test-token", "", "")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "timeout")
}

func TestPoll_VPNStatusField(t *testing.T) {
	useFastPolling(t)

	// VPN uses "status" instead of "state"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]any{
			"metadata": map[string]any{"status": "AVAILABLE"},
		})
	}))
	defer server.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := Poll(ctx, server.URL, "test-token", "", "")
	assert.NoError(t, err)
}

func TestPoll_ActiveState(t *testing.T) {
	useFastPolling(t)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]any{
			"metadata": map[string]any{"state": "ACTIVE"},
		})
	}))
	defer server.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := Poll(ctx, server.URL, "test-token", "", "")
	assert.NoError(t, err)
}

func TestPoll_AuthHeaders(t *testing.T) {
	useFastPolling(t)

	var capturedAuth string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedAuth = r.Header.Get("Authorization")
		json.NewEncoder(w).Encode(map[string]any{
			"metadata": map[string]any{"state": "AVAILABLE"},
		})
	}))
	defer server.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := Poll(ctx, server.URL, "my-token", "", "")
	assert.NoError(t, err)
	assert.Equal(t, "Bearer my-token", capturedAuth)
}

func TestPoll_BasicAuth(t *testing.T) {
	useFastPolling(t)

	var capturedUser, capturedPass string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedUser, capturedPass, _ = r.BasicAuth()
		json.NewEncoder(w).Encode(map[string]any{
			"metadata": map[string]any{"state": "AVAILABLE"},
		})
	}))
	defer server.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := Poll(ctx, server.URL, "", "user", "pass")
	assert.NoError(t, err)
	assert.Equal(t, "user", capturedUser)
	assert.Equal(t, "pass", capturedPass)
}

func TestPoll_DepthParam(t *testing.T) {
	useFastPolling(t)

	var capturedDepth string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedDepth = r.URL.Query().Get("depth")
		json.NewEncoder(w).Encode(map[string]any{
			"metadata": map[string]any{"state": "AVAILABLE"},
		})
	}))
	defer server.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := Poll(ctx, server.URL+"?depth=1", "token", "", "")
	assert.NoError(t, err)
	assert.Equal(t, "1", capturedDepth)
}

func TestPoll_TransientErrors(t *testing.T) {
	useFastPolling(t)

	var callCount atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		count := callCount.Add(1)
		if count < 3 {
			// Return invalid JSON to simulate transient error
			w.Write([]byte("invalid json"))
			return
		}
		json.NewEncoder(w).Encode(map[string]any{
			"metadata": map[string]any{"state": "AVAILABLE"},
		})
	}))
	defer server.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := Poll(ctx, server.URL, "token", "", "")
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, callCount.Load(), int32(3))
}

func TestBuildFullURL(t *testing.T) {
	tests := []struct {
		name     string
		href     string
		expected string
	}{
		{
			name:     "full URL",
			href:     "https://dns.de-fra.ionos.com/zones/zone-1",
			expected: "https://dns.de-fra.ionos.com/zones/zone-1?depth=1",
		},
		{
			name:     "URL with existing query params",
			href:     "https://example.com/resource?foo=bar",
			expected: "https://example.com/resource?foo=bar&depth=1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := buildFullURL(tt.href)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestBuildFullURL_RelativePath(t *testing.T) {
	result := buildFullURL("/cloudapi/v6/datacenters/test-id")
	// Should prepend the base URL and add depth param
	assert.Contains(t, result, "/cloudapi/v6/datacenters/test-id")
	assert.Contains(t, result, "depth=1")
	assert.True(t, len(result) > len("/cloudapi/v6/datacenters/test-id?depth=1"))
}

func TestCaptureRenderInfo(t *testing.T) {
	Reset()

	prefix := "items"
	mapping := map[string]string{"Name": "properties.name", "State": "metadata.state"}
	cols := []string{"Name", "State"}

	CaptureRenderInfo(prefix, mapping, cols)
	ri := GetRenderInfo()

	assert.NotNil(t, ri)
	assert.Equal(t, "items", ri.Prefix)
	assert.Equal(t, mapping, ri.Mapping)
	assert.Equal(t, cols, ri.Cols)
}

func TestCaptureRenderInfo_NilAfterReset(t *testing.T) {
	CaptureRenderInfo("", map[string]string{"A": "a"}, []string{"A"})
	assert.NotNil(t, GetRenderInfo())

	Reset()
	assert.Nil(t, GetRenderInfo())
}

func TestIsRerendering(t *testing.T) {
	Reset()
	assert.False(t, IsRerendering())

	SetRerendering(true)
	assert.True(t, IsRerendering())

	SetRerendering(false)
	assert.False(t, IsRerendering())
}

func TestReset_ClearsAll(t *testing.T) {
	CaptureHref(map[string]any{"id": "1", "href": "/test"})
	CaptureRenderInfo("", map[string]string{"A": "a"}, []string{"A"})
	SetRerendering(true)

	Reset()

	assert.Empty(t, GetHref())
	assert.Nil(t, GetRenderInfo())
	assert.False(t, IsRerendering())
}

func TestFetchResource_NoHref(t *testing.T) {
	Reset()

	_, err := FetchResource()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no href captured")
}

func TestFetchJSON_ReturnsFullResponse(t *testing.T) {
	responseData := map[string]any{
		"id":   "dc-1",
		"href": "/test",
		"metadata": map[string]any{
			"state": "AVAILABLE",
		},
		"properties": map[string]any{
			"name":     "My DC",
			"location": "de/txl",
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(responseData)
	}))
	defer server.Close()

	result, err := fetchJSON(server.URL, "test-token", "", "")
	assert.NoError(t, err)
	assert.NotNil(t, result)

	m, ok := result.(map[string]any)
	assert.True(t, ok)
	assert.Equal(t, "dc-1", m["id"])

	props, ok := m["properties"].(map[string]any)
	assert.True(t, ok)
	assert.Equal(t, "My DC", props["name"])
}

func TestFetchJSON_AuthHeaders(t *testing.T) {
	var capturedAuth string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedAuth = r.Header.Get("Authorization")
		json.NewEncoder(w).Encode(map[string]any{"id": "1"})
	}))
	defer server.Close()

	_, err := fetchJSON(server.URL, "my-token", "", "")
	assert.NoError(t, err)
	assert.Equal(t, "Bearer my-token", capturedAuth)
}

func TestFetchJSON_BasicAuth(t *testing.T) {
	var capturedUser, capturedPass string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedUser, capturedPass, _ = r.BasicAuth()
		json.NewEncoder(w).Encode(map[string]any{"id": "1"})
	}))
	defer server.Close()

	_, err := fetchJSON(server.URL, "", "user", "pass")
	assert.NoError(t, err)
	assert.Equal(t, "user", capturedUser)
	assert.Equal(t, "pass", capturedPass)
}
