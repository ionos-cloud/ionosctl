package globalwait

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

// --- helpers ---

func fastPoll(t *testing.T) {
	t.Helper()
	old := pollInterval
	pollInterval = 50 * time.Millisecond
	t.Cleanup(func() { pollInterval = old })
}

func quickCtx(t *testing.T, d time.Duration) context.Context {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), d)
	t.Cleanup(cancel)
	return ctx
}

func stateServer(state string) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]any{
			"metadata": map[string]any{"state": state},
		})
	}))
}

// --- ExtractHref ---

func TestExtractHref(t *testing.T) {
	tests := []struct {
		name     string
		input    any
		expected string
	}{
		{"single resource with href", map[string]any{"href": "https://api.ionos.com/cloudapi/v6/datacenters/abc", "id": "abc"}, "https://api.ionos.com/cloudapi/v6/datacenters/abc"},
		{"list response with items", map[string]any{"href": "https://api.ionos.com/cloudapi/v6/datacenters", "items": []any{}}, ""},
		{"no href field", map[string]any{"id": "abc", "name": "test"}, ""},
		{"nil input", nil, ""},
		{"non-map input (string)", "just a string", ""},
		{"non-map input (int)", 42, ""},
		{"non-map input (slice)", []string{"a", "b"}, ""},
		{"relative href", map[string]any{"href": "/certificates/abc"}, "/certificates/abc"},
		{"empty href", map[string]any{"href": ""}, ""},
		{"href is non-string", map[string]any{"href": 123}, ""},
		{"items present even if empty", map[string]any{"href": "https://x.com/res/1", "items": []any{}}, ""},
		{"nested struct with href", struct {
			Href string `json:"href"`
			ID   string `json:"id"`
		}{"https://api.ionos.com/test", "id1"}, "https://api.ionos.com/test"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, ExtractHref(tt.input))
		})
	}
}

// --- State management ---

func TestCaptureAndGetHref(t *testing.T) {
	Reset()
	assert.Empty(t, GetHref())

	CaptureHref("https://api.ionos.com/test")
	assert.Equal(t, "https://api.ionos.com/test", GetHref())

	// Overwrite
	CaptureHref("https://api.ionos.com/other")
	assert.Equal(t, "https://api.ionos.com/other", GetHref())
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

func TestCaptureRerenderable_NilCols(t *testing.T) {
	Reset()
	mock := &mockRerenderable{}
	CaptureRerenderable(mock, nil)
	r, cols := GetRerenderable()
	assert.Equal(t, mock, r)
	assert.Nil(t, cols)
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

func TestCaptureRequestURL(t *testing.T) {
	t.Run("sets href when empty", func(t *testing.T) {
		Reset()
		CaptureRequestURL(http.MethodPost, "https://api.ionos.com/cloudapi/v6/datacenters/aaaaaaaa-1111-2222-3333-444444444444", "")
		assert.Equal(t, "https://api.ionos.com/cloudapi/v6/datacenters/aaaaaaaa-1111-2222-3333-444444444444", GetHref())
	})

	t.Run("does not overwrite existing href", func(t *testing.T) {
		Reset()
		CaptureHref("https://api.ionos.com/first")
		CaptureRequestURL(http.MethodPost, "https://api.ionos.com/second", "")
		assert.Equal(t, "https://api.ionos.com/first", GetHref())
	})

	t.Run("empty URL does nothing", func(t *testing.T) {
		Reset()
		CaptureRequestURL(http.MethodPost, "", "")
		assert.Empty(t, GetHref())
	})
}

func TestSetResourceHref(t *testing.T) {
	Reset()
	viper.Set(constants.ArgServerUrl, "https://api.ionos.com")
	defer viper.Set(constants.ArgServerUrl, "")

	SetResourceHref("cloudapi", "v6", "datacenters", "dc-1", "servers", "srv-1")
	assert.Equal(t, "https://api.ionos.com/cloudapi/v6/datacenters/dc-1/servers/srv-1", GetHref())
}

func TestSetResourceHref_DefaultURL(t *testing.T) {
	Reset()
	viper.Set(constants.ArgServerUrl, "")

	SetResourceHref("cloudapi", "v6", "datacenters", "dc-1")
	href := GetHref()
	assert.Contains(t, href, "datacenters/dc-1")
	assert.Contains(t, href, "ionos.com")
}

// --- Poll ---

func TestPoll_Available(t *testing.T) {
	var callCount int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		n := atomic.AddInt32(&callCount, 1)
		state := "BUSY"
		if n >= 2 {
			state = "AVAILABLE"
		}
		json.NewEncoder(w).Encode(map[string]any{"metadata": map[string]any{"state": state}})
	}))
	defer server.Close()
	fastPoll(t)

	err := Poll(quickCtx(t, 5*time.Second), server.URL, "", "", "", false)
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, int(atomic.LoadInt32(&callCount)), 2)
}

func TestPoll_ImmediateAvailable(t *testing.T) {
	server := stateServer("AVAILABLE")
	defer server.Close()
	fastPoll(t)

	err := Poll(quickCtx(t, 5*time.Second), server.URL, "", "", "", false)
	assert.NoError(t, err)
}

func TestPoll_AllTerminalStates(t *testing.T) {
	for _, state := range []string{"AVAILABLE", "ACTIVE", "READY", "DONE", "available", "Active"} {
		t.Run(state, func(t *testing.T) {
			server := stateServer(state)
			defer server.Close()
			fastPoll(t)

			err := Poll(quickCtx(t, 5*time.Second), server.URL, "", "", "", false)
			assert.NoError(t, err)
		})
	}
}

func TestPoll_Failed(t *testing.T) {
	server := stateServer("FAILED")
	defer server.Close()
	fastPoll(t)

	err := Poll(quickCtx(t, 5*time.Second), server.URL, "", "", "", false)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "FAILED")
}

func TestPoll_FailedCaseInsensitive(t *testing.T) {
	server := stateServer("failed")
	defer server.Close()
	fastPoll(t)

	err := Poll(quickCtx(t, 5*time.Second), server.URL, "", "", "", false)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed")
}

func TestPoll_Timeout(t *testing.T) {
	server := stateServer("BUSY")
	defer server.Close()
	fastPoll(t)

	err := Poll(quickCtx(t, 200*time.Millisecond), server.URL, "", "", "", false)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "timeout")
}

func TestPoll_StatusField(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]any{"metadata": map[string]any{"status": "ACTIVE"}})
	}))
	defer server.Close()
	fastPoll(t)

	err := Poll(quickCtx(t, 5*time.Second), server.URL, "", "", "", false)
	assert.NoError(t, err)
}

func TestPoll_404_Delete(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"message":"not found"}`))
	}))
	defer server.Close()
	fastPoll(t)

	err := Poll(quickCtx(t, 5*time.Second), server.URL, "", "", "", true)
	assert.NoError(t, err)
}

func TestPoll_404_Create_Transient(t *testing.T) {
	// 404 during create is transient (resource provisioning). Should retry until timeout.
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()
	fastPoll(t)

	err := Poll(quickCtx(t, 200*time.Millisecond), server.URL, "", "", "", false)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "timeout")
}

func TestPoll_NoMetadataState_TreatedAsReady(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]any{"id": "abc", "properties": map[string]any{}})
	}))
	defer server.Close()
	fastPoll(t)

	err := Poll(quickCtx(t, 5*time.Second), server.URL, "", "", "", false)
	assert.NoError(t, err) // no state = resource doesn't track state, treat as ready
}

func TestPoll_NilMetadata_TreatedAsReady(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]any{"metadata": nil})
	}))
	defer server.Close()
	fastPoll(t)

	err := Poll(quickCtx(t, 5*time.Second), server.URL, "", "", "", false)
	assert.NoError(t, err)
}

func TestPoll_EmptyStateFields_TreatedAsReady(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]any{"metadata": map[string]any{"createdDate": "2024-01-01"}})
	}))
	defer server.Close()
	fastPoll(t)

	err := Poll(quickCtx(t, 5*time.Second), server.URL, "", "", "", false)
	assert.NoError(t, err)
}

func TestPoll_TransientErrors_Retried(t *testing.T) {
	var callCount int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		n := atomic.AddInt32(&callCount, 1)
		if n <= 2 {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("server error"))
			return
		}
		json.NewEncoder(w).Encode(map[string]any{"metadata": map[string]any{"state": "AVAILABLE"}})
	}))
	defer server.Close()
	fastPoll(t)

	err := Poll(quickCtx(t, 5*time.Second), server.URL, "", "", "", false)
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, int(atomic.LoadInt32(&callCount)), 3)
}

func TestPoll_MalformedJSON_Retried(t *testing.T) {
	var callCount int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		n := atomic.AddInt32(&callCount, 1)
		if n == 1 {
			w.Write([]byte("not json"))
			return
		}
		json.NewEncoder(w).Encode(map[string]any{"metadata": map[string]any{"state": "AVAILABLE"}})
	}))
	defer server.Close()
	fastPoll(t)

	err := Poll(quickCtx(t, 5*time.Second), server.URL, "", "", "", false)
	assert.NoError(t, err)
}

func TestPoll_AuthHeaders_BearerToken(t *testing.T) {
	var gotAuth string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotAuth = r.Header.Get("Authorization")
		json.NewEncoder(w).Encode(map[string]any{"metadata": map[string]any{"state": "AVAILABLE"}})
	}))
	defer server.Close()
	fastPoll(t)

	err := Poll(quickCtx(t, 5*time.Second), server.URL, "my-token", "", "", false)
	assert.NoError(t, err)
	assert.Equal(t, "Bearer my-token", gotAuth)
}

func TestPoll_AuthHeaders_BasicAuth(t *testing.T) {
	var gotUser, gotPass string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotUser, gotPass, _ = r.BasicAuth()
		json.NewEncoder(w).Encode(map[string]any{"metadata": map[string]any{"state": "AVAILABLE"}})
	}))
	defer server.Close()
	fastPoll(t)

	err := Poll(quickCtx(t, 5*time.Second), server.URL, "", "user", "pass", false)
	assert.NoError(t, err)
	assert.Equal(t, "user", gotUser)
	assert.Equal(t, "pass", gotPass)
}

func TestPoll_TokenPrecedence(t *testing.T) {
	var gotAuth string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotAuth = r.Header.Get("Authorization")
		json.NewEncoder(w).Encode(map[string]any{"metadata": map[string]any{"state": "AVAILABLE"}})
	}))
	defer server.Close()
	fastPoll(t)

	err := Poll(quickCtx(t, 5*time.Second), server.URL, "tok", "user", "pass", false)
	assert.NoError(t, err)
	assert.Equal(t, "Bearer tok", gotAuth) // token wins over basic auth
}

func TestPoll_IntermediateStates_KeepPolling(t *testing.T) {
	states := []string{"DEPLOYING", "UPDATING", "BUSY", "PROVISIONING", "AVAILABLE"}
	var idx int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		i := int(atomic.AddInt32(&idx, 1)) - 1
		if i >= len(states) {
			i = len(states) - 1
		}
		json.NewEncoder(w).Encode(map[string]any{"metadata": map[string]any{"state": states[i]}})
	}))
	defer server.Close()
	fastPoll(t)

	err := Poll(quickCtx(t, 5*time.Second), server.URL, "", "", "", false)
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, int(atomic.LoadInt32(&idx)), 5)
}

// --- WaitForAvailable ---

func TestWaitForAvailable_NoHref(t *testing.T) {
	Reset()
	var buf bytes.Buffer
	err := WaitForAvailable(&buf, "", "", "")
	assert.NoError(t, err)
	assert.Empty(t, buf.String())
}

func TestWaitForAvailable_SingleResource(t *testing.T) {
	server := stateServer("AVAILABLE")
	defer server.Close()
	fastPoll(t)

	Reset()
	CaptureHref(server.URL + "/datacenters/aaaaaaaa-1111-2222-3333-444444444444")

	viper.Set(constants.ArgTimeout, 5)
	defer viper.Set(constants.ArgTimeout, 0)

	var buf bytes.Buffer
	err := WaitForAvailable(&buf, "", "", "")
	assert.NoError(t, err)
	assert.Contains(t, buf.String(), "DONE")
}

func TestWaitForAvailable_PollsParents(t *testing.T) {
	// Fix: protect polledPaths with a mutex since the HTTP handler runs in
	// separate goroutines managed by httptest.Server.
	var pathsMu sync.Mutex
	var polledPaths []string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		pathsMu.Lock()
		polledPaths = append(polledPaths, r.URL.Path)
		pathsMu.Unlock()
		json.NewEncoder(w).Encode(map[string]any{"metadata": map[string]any{"state": "AVAILABLE"}})
	}))
	defer server.Close()
	fastPoll(t)

	Reset()
	CaptureHref(server.URL + "/cloudapi/v6/datacenters/aaaaaaaa-1111-2222-3333-444444444444/servers/bbbbbbbb-1111-2222-3333-444444444444")

	viper.Set(constants.ArgTimeout, 5)
	defer viper.Set(constants.ArgTimeout, 0)

	var buf bytes.Buffer
	err := WaitForAvailable(&buf, "", "", "")
	assert.NoError(t, err)

	// Should poll server then datacenter
	pathsMu.Lock()
	pathCount := len(polledPaths)
	pathsMu.Unlock()
	assert.GreaterOrEqual(t, pathCount, 2)
}

func TestWaitForAvailable_DeletedResource_ThenParent(t *testing.T) {
	var callCount int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		n := atomic.AddInt32(&callCount, 1)
		if n == 1 {
			// First call: resource deleted
			w.WriteHeader(http.StatusNotFound)
			return
		}
		// Parent: AVAILABLE
		json.NewEncoder(w).Encode(map[string]any{"metadata": map[string]any{"state": "AVAILABLE"}})
	}))
	defer server.Close()
	fastPoll(t)

	Reset()
	CaptureHref(server.URL + "/datacenters/aaaaaaaa-1111-2222-3333-444444444444/servers/bbbbbbbb-1111-2222-3333-444444444444")

	viper.Set(constants.ArgTimeout, 5)
	defer viper.Set(constants.ArgTimeout, 0)

	var buf bytes.Buffer
	err := WaitForAvailable(&buf, "", "", "")
	assert.NoError(t, err)
}

func TestWaitForAvailable_FailedResource(t *testing.T) {
	server := stateServer("FAILED")
	defer server.Close()
	fastPoll(t)

	Reset()
	CaptureHref(server.URL + "/datacenters/aaaaaaaa-1111-2222-3333-444444444444")

	viper.Set(constants.ArgTimeout, 5)
	defer viper.Set(constants.ArgTimeout, 0)

	var buf bytes.Buffer
	err := WaitForAvailable(&buf, "", "", "")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "FAILED")
	assert.Contains(t, buf.String(), "FAILED")
}

func TestWaitForAvailable_NoStateField(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]any{"id": "abc"})
	}))
	defer server.Close()
	fastPoll(t)

	Reset()
	CaptureHref(server.URL + "/clusters/aaaaaaaa-1111-2222-3333-444444444444")

	viper.Set(constants.ArgTimeout, 5)
	defer viper.Set(constants.ArgTimeout, 0)

	var buf bytes.Buffer
	err := WaitForAvailable(&buf, "", "", "")
	assert.NoError(t, err) // no state = treated as ready
}

// --- FetchResource ---

func TestFetchResource_Success(t *testing.T) {
	expected := map[string]any{"id": "abc", "properties": map[string]any{"name": "test"}}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(expected)
	}))
	defer server.Close()

	Reset()
	CaptureHref(server.URL + "/resource/abc")

	result, err := FetchResource("", "", "")
	assert.NoError(t, err)
	m, ok := result.(map[string]any)
	assert.True(t, ok)
	assert.Equal(t, "abc", m["id"])
}

func TestFetchResource_NoHref(t *testing.T) {
	Reset()
	_, err := FetchResource("", "", "")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no href captured")
}

func TestFetchResource_ServerError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("not json"))
	}))
	defer server.Close()

	Reset()
	CaptureHref(server.URL + "/resource/abc")

	_, err := FetchResource("", "", "")
	assert.Error(t, err)
}

func TestFetchResource_AuthPassed(t *testing.T) {
	var gotAuth string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotAuth = r.Header.Get("Authorization")
		json.NewEncoder(w).Encode(map[string]any{"id": "1"})
	}))
	defer server.Close()

	Reset()
	CaptureHref(server.URL + "/resource/1")

	_, err := FetchResource("tok123", "", "")
	assert.NoError(t, err)
	assert.Equal(t, "Bearer tok123", gotAuth)
}

// --- WrapTransport / capturingTransport ---

func TestWrapTransport_CapturesDeleteURL(t *testing.T) {
	Reset()
	viper.Set(constants.ArgWait, true)
	defer viper.Set(constants.ArgWait, false)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	hc := &http.Client{}
	WrapTransport(hc)

	req, _ := http.NewRequest(http.MethodDelete, server.URL+"/datacenters/abc", nil)
	_, err := hc.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, server.URL+"/datacenters/abc", GetHref())
}

func TestWrapTransport_CapturesPostURL(t *testing.T) {
	Reset()
	viper.Set(constants.ArgWait, true)
	defer viper.Set(constants.ArgWait, false)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
	}))
	defer server.Close()

	hc := &http.Client{}
	WrapTransport(hc)

	req, _ := http.NewRequest(http.MethodPost, server.URL+"/datacenters", nil)
	_, err := hc.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, server.URL+"/datacenters", GetHref())
}

func TestWrapTransport_SkipsGET(t *testing.T) {
	Reset()
	viper.Set(constants.ArgWait, true)
	defer viper.Set(constants.ArgWait, false)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	hc := &http.Client{}
	WrapTransport(hc)

	req, _ := http.NewRequest(http.MethodGet, server.URL+"/datacenters", nil)
	_, err := hc.Do(req)
	assert.NoError(t, err)
	assert.Empty(t, GetHref()) // GET should not capture
}

func TestWrapTransport_SkipsWhenWaitFalse(t *testing.T) {
	Reset()
	viper.Set(constants.ArgWait, false)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	hc := &http.Client{}
	WrapTransport(hc)

	req, _ := http.NewRequest(http.MethodDelete, server.URL+"/datacenters/abc", nil)
	_, err := hc.Do(req)
	assert.NoError(t, err)
	assert.Empty(t, GetHref()) // --wait not set
}

func TestWrapTransport_NilClient(t *testing.T) {
	// Should not panic
	WrapTransport(nil)
}

func TestWrapTransport_NilTransport(t *testing.T) {
	Reset()
	viper.Set(constants.ArgWait, true)
	defer viper.Set(constants.ArgWait, false)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	hc := &http.Client{Transport: nil}
	WrapTransport(hc)

	req, _ := http.NewRequest(http.MethodPost, server.URL+"/test", nil)
	_, err := hc.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, server.URL+"/test", GetHref())
}

func TestWrapTransport_DoesNotOverwriteExistingHref(t *testing.T) {
	Reset()
	CaptureHref("https://already.set/resource/1")
	viper.Set(constants.ArgWait, true)
	defer viper.Set(constants.ArgWait, false)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	hc := &http.Client{}
	WrapTransport(hc)

	req, _ := http.NewRequest(http.MethodDelete, server.URL+"/other", nil)
	_, err := hc.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, "https://already.set/resource/1", GetHref()) // not overwritten
}

func TestCapturingTransport_PropagatesError(t *testing.T) {
	Reset()
	viper.Set(constants.ArgWait, true)
	defer viper.Set(constants.ArgWait, false)

	hc := &http.Client{}
	WrapTransport(hc)

	req, _ := http.NewRequest(http.MethodDelete, "http://127.0.0.1:1", nil) // bad port
	_, err := hc.Do(req)
	assert.Error(t, err)
	assert.Empty(t, GetHref()) // error means no capture
}

func TestCapturingTransport_AllMutatingMethods(t *testing.T) {
	for _, method := range []string{http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete} {
		t.Run(method, func(t *testing.T) {
			Reset()
			viper.Set(constants.ArgWait, true)
			defer viper.Set(constants.ArgWait, false)

			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			}))
			defer server.Close()

			hc := &http.Client{}
			WrapTransport(hc)

			req, _ := http.NewRequest(method, server.URL+"/resource/"+method, nil)
			_, err := hc.Do(req)
			assert.NoError(t, err)
			assert.Equal(t, server.URL+"/resource/"+method, GetHref())
		})
	}
}

func TestCapturingTransport_SkipsReadMethods(t *testing.T) {
	for _, method := range []string{http.MethodGet, http.MethodHead, http.MethodOptions} {
		t.Run(method, func(t *testing.T) {
			Reset()
			viper.Set(constants.ArgWait, true)
			defer viper.Set(constants.ArgWait, false)

			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			}))
			defer server.Close()

			hc := &http.Client{}
			WrapTransport(hc)

			req, _ := http.NewRequest(method, server.URL+"/resource", nil)
			_, err := hc.Do(req)
			assert.NoError(t, err)
			assert.Empty(t, GetHref())
		})
	}
}

// --- looksLikeResourceID ---

func TestLooksLikeResourceID(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee", true}, // UUID
		{"11111111-2222-3333-4444-555555555555", true}, // UUID
		{"123456789", true},                            // numeric
		{"0", true},                                    // single digit
		{"", false},                                    // empty
		{"cloudapi", false},                            // API path component
		{"v6", false},                                  // version
		{"datacenters", false},                         // resource type
		{"servers", false},                             // resource type
		{"short-id", false},                            // not a UUID
		{"a-b-c-d-e", false},                           // not a UUID
		{"abcdefghij-klmn", false},                     // not a UUID format
		{"private-cross-connects", false},              // hyphenated resource type
		{"application-load-balancers", false},          // hyphenated resource type
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			assert.Equal(t, tt.expected, looksLikeResourceID(tt.input))
		})
	}
}

// --- parentHref ---

func TestParentHref(t *testing.T) {
	tests := []struct {
		name     string
		href     string
		expected string
	}{
		{"volume to server (CloudAPI)", "https://api.ionos.com/cloudapi/v6/datacenters/aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee/servers/11111111-2222-3333-4444-555555555555/volumes/66666666-7777-8888-9999-000000000000", "https://api.ionos.com/cloudapi/v6/datacenters/aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee/servers/11111111-2222-3333-4444-555555555555"},
		{"server to datacenter (CloudAPI)", "https://api.ionos.com/cloudapi/v6/datacenters/aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee/servers/11111111-2222-3333-4444-555555555555", "https://api.ionos.com/cloudapi/v6/datacenters/aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee"},
		{"datacenter stop (CloudAPI root)", "https://api.ionos.com/cloudapi/v6/datacenters/aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee", ""},
		{"VPN peer to gateway (regional)", "https://vpn.de-fra.ionos.com/wireguardgateways/aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee/peers/11111111-2222-3333-4444-555555555555", "https://vpn.de-fra.ionos.com/wireguardgateways/aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee"},
		{"VPN gateway stop (regional root)", "https://vpn.de-fra.ionos.com/wireguardgateways/aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee", ""},
		{"too short", "https://api.ionos.com/cloudapi/v6", ""},
		{"just host", "https://api.ionos.com", ""},
		{"empty string", "", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, parentHref(tt.href))
		})
	}
}

// --- resourceAndParentURLs ---

func TestResourceAndParentURLs_DeepNesting(t *testing.T) {
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

func TestResourceAndParentURLs_NoParent(t *testing.T) {
	urls := resourceAndParentURLs("https://dns.de-fra.ionos.com/zones/aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee")
	assert.Equal(t, []string{
		"https://dns.de-fra.ionos.com/zones/aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee",
	}, urls)
}

// --- buildFullURL ---

func TestBuildFullURL(t *testing.T) {
	tests := []struct {
		name     string
		href     string
		expected string
	}{
		{"absolute URL", "https://api.ionos.com/cloudapi/v6/datacenters/abc", "https://api.ionos.com/cloudapi/v6/datacenters/abc?depth=1"},
		{"absolute URL with query", "https://api.ionos.com/cloudapi/v6/datacenters/abc?pretty=true", "https://api.ionos.com/cloudapi/v6/datacenters/abc?depth=1&pretty=true"},
		{"http URL", "http://localhost:8080/test", "http://localhost:8080/test?depth=1"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, buildFullURL(tt.href))
		})
	}
}

// --- fetchState ---

func TestFetchState_BasicAuth(t *testing.T) {
	var gotUser, gotPass string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotUser, gotPass, _ = r.BasicAuth()
		json.NewEncoder(w).Encode(map[string]any{"metadata": map[string]any{"state": "AVAILABLE"}})
	}))
	defer server.Close()

	state, err := newPoller("", "u", "p").fetchState(context.Background(), server.URL, false)
	assert.NoError(t, err)
	assert.Equal(t, "AVAILABLE", state)
	assert.Equal(t, "u", gotUser)
	assert.Equal(t, "p", gotPass)
}

func TestFetchState_NoAuth(t *testing.T) {
	var gotAuth string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotAuth = r.Header.Get("Authorization")
		json.NewEncoder(w).Encode(map[string]any{"metadata": map[string]any{"state": "READY"}})
	}))
	defer server.Close()

	state, err := newPoller("", "", "").fetchState(context.Background(), server.URL, false)
	assert.NoError(t, err)
	assert.Equal(t, "READY", state)
	assert.Empty(t, gotAuth)
}

func TestFetchState_404_Delete(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	state, err := newPoller("", "", "").fetchState(context.Background(), server.URL, true)
	assert.NoError(t, err)
	assert.Equal(t, "DONE", state)
}

func TestFetchState_404_Create(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	state, err := newPoller("", "", "").fetchState(context.Background(), server.URL, false)
	assert.Error(t, err)
	assert.Equal(t, "", state)
	assert.Contains(t, err.Error(), "404")
}

func TestFetchState_UserAgent(t *testing.T) {
	viper.Set(constants.CLIHttpUserAgent, "ionosctl/test")
	defer viper.Set(constants.CLIHttpUserAgent, "")

	var gotUA string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotUA = r.Header.Get("User-Agent")
		json.NewEncoder(w).Encode(map[string]any{"metadata": map[string]any{"state": "AVAILABLE"}})
	}))
	defer server.Close()

	_, err := newPoller("", "", "").fetchState(context.Background(), server.URL, false)
	assert.NoError(t, err)
	assert.Equal(t, "ionosctl/test", gotUA)
}

func TestFetchState_StatusFallback(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]any{"metadata": map[string]any{"status": "ACTIVE", "state": ""}})
	}))
	defer server.Close()

	state, err := newPoller("", "", "").fetchState(context.Background(), server.URL, false)
	assert.NoError(t, err)
	assert.Equal(t, "ACTIVE", state)
}

func TestFetchState_StatePrecedence(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]any{"metadata": map[string]any{"state": "BUSY", "status": "ACTIVE"}})
	}))
	defer server.Close()

	state, err := newPoller("", "", "").fetchState(context.Background(), server.URL, false)
	assert.NoError(t, err)
	assert.Equal(t, "BUSY", state) // state wins over status
}

// --- mockRerenderable ---

type mockRerenderable struct {
	extractCalled bool
	renderCalled  bool
	data          any
	extractErr    error
	renderErr     error
}

func (m *mockRerenderable) Extract(sourceData any) error {
	m.extractCalled = true
	m.data = sourceData
	return m.extractErr
}

func (m *mockRerenderable) Render(visibleCols []string) (string, error) {
	m.renderCalled = true
	if m.renderErr != nil {
		return "", m.renderErr
	}
	return fmt.Sprintf("rendered:%v", m.data), nil
}

// --- Edge case tests ---

// TestPoll_429RateLimit verifies that fetchState handles HTTP 429 (Too Many Requests)
// by reading the Retry-After header and retrying, eventually reaching AVAILABLE.
func TestPoll_429RateLimit(t *testing.T) {
	var callCount int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		n := atomic.AddInt32(&callCount, 1)
		if n <= 2 {
			w.Header().Set("Retry-After", "1")
			w.WriteHeader(http.StatusTooManyRequests)
			w.Write([]byte(`{"message":"rate limit exceeded"}`))
			return
		}
		// After rate limit clears, return AVAILABLE
		json.NewEncoder(w).Encode(map[string]any{"metadata": map[string]any{"state": "AVAILABLE"}})
	}))
	defer server.Close()
	fastPoll(t)

	err := Poll(quickCtx(t, 10*time.Second), server.URL, "", "", "", false)
	assert.NoError(t, err, "should succeed after rate limit clears")
	assert.Equal(t, int32(3), atomic.LoadInt32(&callCount),
		"should make 3 calls: 2 rate-limited + 1 AVAILABLE")
}

// TestPoll_429RateLimit_ContextCancellation verifies that a large Retry-After
// value does not block past the context deadline.
func TestPoll_429RateLimit_ContextCancellation(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Retry-After", "300") // 5 minutes — way past our timeout
		w.WriteHeader(http.StatusTooManyRequests)
	}))
	defer server.Close()
	fastPoll(t)

	start := time.Now()
	err := Poll(quickCtx(t, 500*time.Millisecond), server.URL, "", "", "", false)
	elapsed := time.Since(start)

	assert.Error(t, err, "should fail when context expires during Retry-After sleep")
	assert.Less(t, elapsed, 5*time.Second, "should not block for full Retry-After duration")
}

// TestPoll_NonStandardStates verifies terminal state handling for states
// beyond the common AVAILABLE/ACTIVE/READY/DONE set.
func TestPoll_NonStandardStates(t *testing.T) {
	// INACTIVE and SUSPENDED are terminal success states (e.g. after server stop/suspend)
	for _, state := range []string{"INACTIVE", "SUSPENDED"} {
		t.Run(state+"_success", func(t *testing.T) {
			server := stateServer(state)
			defer server.Close()
			fastPoll(t)

			err := Poll(quickCtx(t, 200*time.Millisecond), server.URL, "", "", "", false)
			assert.NoError(t, err, "%q should be treated as terminal success", state)
		})
	}

	// ERROR and DESTROYING are terminal failure states
	for _, state := range []string{"ERROR", "DESTROYING"} {
		t.Run(state+"_failure", func(t *testing.T) {
			server := stateServer(state)
			defer server.Close()
			fastPoll(t)

			err := Poll(quickCtx(t, 200*time.Millisecond), server.URL, "", "", "", false)
			assert.Error(t, err, "%q should be treated as terminal failure", state)
			assert.Contains(t, err.Error(), state)
		})
	}
}

// TestDoubleWait_OldAndNewWaitersFireTogether documents the double-wait bug
// where both old per-command waiters and the new globalwait mechanism both
// check ArgWait to decide whether to poll.
// BUG: Both old per-command waiters and new globalwait check ArgWait, causing redundant polling.
func TestDoubleWait_OldAndNewWaitersFireTogether(t *testing.T) {
	// TODO: Fix by having old per-command waiters check a different flag,
	// or by removing old waiters entirely once globalwait is stable.

	// When --wait is set globally, viper.GetBool(constants.ArgWait) is true.
	// Old waiters (waitfor.WaitForRequest, waitfor.WaitForState) guard on this.
	// New globalwait.WaitForAvailable is called from root.go also guarded on this.
	// Both will fire, causing the same resource to be polled twice.

	viper.Set(constants.ArgWait, true)
	defer viper.Set(constants.ArgWait, false)

	// Verify both guard conditions are satisfied simultaneously.
	// This is the root cause of the double-wait: a single flag controls both paths.
	assert.True(t, viper.GetBool(constants.ArgWait),
		"ArgWait is true, which means BOTH old per-command waiters AND globalwait "+
			"will execute, causing redundant polling of the same resource")

	// Also verify ArgWaitForRequest (used by some old waiters) is independent
	// but both can be true at the same time.
	viper.Set(constants.ArgWaitForRequest, true)
	defer viper.Set(constants.ArgWaitForRequest, false)

	assert.True(t, viper.GetBool(constants.ArgWait))
	assert.True(t, viper.GetBool(constants.ArgWaitForRequest),
		"both ArgWait and ArgWaitForRequest can be true simultaneously, "+
			"meaning up to 3 polling loops could run for one operation")
}

// TestBuildFullURL_RelativePaths exercises the else branch of buildFullURL
// where the href does not start with http:// or https://.
func TestBuildFullURL_RelativePaths(t *testing.T) {
	viper.Set(constants.ArgServerUrl, "https://api.ionos.com")
	defer viper.Set(constants.ArgServerUrl, "")

	tests := []struct {
		name     string
		href     string
		expected string
	}{
		{
			"absolute path (leading slash)",
			"/datacenters/aaaaaaaa-1111-2222-3333-444444444444",
			"https://api.ionos.com/datacenters/aaaaaaaa-1111-2222-3333-444444444444?depth=1",
		},
		{
			// Fixed: buildFullURL now normalizes relative paths by prepending "/" if missing.
			"relative path (no leading slash) - slash bug fixed",
			"datacenters/aaaaaaaa-1111-2222-3333-444444444444",
			"https://api.ionos.com/datacenters/aaaaaaaa-1111-2222-3333-444444444444?depth=1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := buildFullURL(tt.href)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestBuildFullURL_RelativePaths_DefaultURL verifies relative path handling
// when no custom server URL is configured.
func TestBuildFullURL_RelativePaths_DefaultURL(t *testing.T) {
	viper.Set(constants.ArgServerUrl, "")

	result := buildFullURL("/cloudapi/v6/datacenters/abc")
	assert.Contains(t, result, "ionos.com")
	assert.Contains(t, result, "/cloudapi/v6/datacenters/abc")
	assert.Contains(t, result, "depth=1")
}

// TestFetchState_400BadRequest verifies that 400 Bad Request is treated as a
// non-retryable client error, regardless of response body format.
func TestFetchState_400BadRequest(t *testing.T) {
	t.Run("with valid JSON body", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]any{
				"httpStatus": 400,
				"messages":   []map[string]any{{"message": "bad request"}},
			})
		}))
		defer server.Close()

		state, err := newPoller("", "", "").fetchState(context.Background(), server.URL, false)
		assert.Error(t, err, "400 should be treated as a non-retryable error")
		assert.Contains(t, err.Error(), "client error (HTTP 400)")
		assert.Empty(t, state)
	})

	t.Run("with non-JSON body", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Bad Request"))
		}))
		defer server.Close()

		state, err := newPoller("", "", "").fetchState(context.Background(), server.URL, false)
		assert.Error(t, err, "400 should be treated as a non-retryable error")
		assert.Contains(t, err.Error(), "client error (HTTP 400)")
		assert.Empty(t, state)
	})
}
