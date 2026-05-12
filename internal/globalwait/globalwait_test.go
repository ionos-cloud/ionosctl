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

func fastpollURL(t *testing.T) {
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

// --- extractHref ---

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
			assert.Equal(t, tt.expected, extractHref(tt.input))
		})
	}
}

// --- State management ---

func TestCaptureAndGetHref(t *testing.T) {
	Reset()
	assert.Empty(t, getHref())

	captureHref("https://api.ionos.com/test")
	assert.Equal(t, "https://api.ionos.com/test", getHref())

	// Overwrite
	captureHref("https://api.ionos.com/other")
	assert.Equal(t, "https://api.ionos.com/other", getHref())
}

func TestCaptureAndGetRerenderable(t *testing.T) {
	Reset()
	r, cols := getRerenderable()
	assert.Nil(t, r)
	assert.Nil(t, cols)

	mock := &mockRerenderable{}
	captureRerenderable(mock, []string{"Col1", "Col2"})

	r, cols = getRerenderable()
	assert.Equal(t, mock, r)
	assert.Equal(t, []string{"Col1", "Col2"}, cols)
}

func TestCaptureRerenderable_NilCols(t *testing.T) {
	Reset()
	mock := &mockRerenderable{}
	captureRerenderable(mock, nil)
	r, cols := getRerenderable()
	assert.Equal(t, mock, r)
	assert.Nil(t, cols)
}

func TestIsRerendering(t *testing.T) {
	Reset()
	assert.False(t, isRerendering())
	setRerendering(true)
	assert.True(t, isRerendering())
	setRerendering(false)
	assert.False(t, isRerendering())
}

func TestReset(t *testing.T) {
	captureHref("test-href")
	captureRerenderable(&mockRerenderable{}, []string{"col"})
	setRerendering(true)

	Reset()

	assert.Empty(t, getHref())
	r, cols := getRerenderable()
	assert.Nil(t, r)
	assert.Nil(t, cols)
	assert.False(t, isRerendering())
}

func TestCaptureRequestURL(t *testing.T) {
	t.Run("sets href when empty", func(t *testing.T) {
		Reset()
		captureRequestURL(http.MethodPost, "https://api.ionos.com/cloudapi/v6/datacenters/aaaaaaaa-1111-2222-3333-444444444444", "")
		assert.Equal(t, "https://api.ionos.com/cloudapi/v6/datacenters/aaaaaaaa-1111-2222-3333-444444444444", getHref())
	})

	t.Run("does not overwrite existing href", func(t *testing.T) {
		Reset()
		captureHref("https://api.ionos.com/first")
		captureRequestURL(http.MethodPost, "https://api.ionos.com/second", "")
		assert.Equal(t, "https://api.ionos.com/first", getHref())
	})

	t.Run("empty URL does nothing", func(t *testing.T) {
		Reset()
		captureRequestURL(http.MethodPost, "", "")
		assert.Empty(t, getHref())
	})
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
	fastpollURL(t)

	err := pollURL(quickCtx(t, 5*time.Second), server.URL, "", "", "", false)
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, int(atomic.LoadInt32(&callCount)), 2)
}

func TestPoll_ImmediateAvailable(t *testing.T) {
	server := stateServer("AVAILABLE")
	defer server.Close()
	fastpollURL(t)

	err := pollURL(quickCtx(t, 5*time.Second), server.URL, "", "", "", false)
	assert.NoError(t, err)
}

func TestPoll_AllTerminalStates(t *testing.T) {
	for _, state := range []string{"AVAILABLE", "ACTIVE", "READY", "DONE", "available", "Active"} {
		t.Run(state, func(t *testing.T) {
			server := stateServer(state)
			defer server.Close()
			fastpollURL(t)

			err := pollURL(quickCtx(t, 5*time.Second), server.URL, "", "", "", false)
			assert.NoError(t, err)
		})
	}
}

func TestPoll_Failed(t *testing.T) {
	server := stateServer("FAILED")
	defer server.Close()
	fastpollURL(t)

	err := pollURL(quickCtx(t, 5*time.Second), server.URL, "", "", "", false)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "FAILED")
}

func TestPoll_FailedCaseInsensitive(t *testing.T) {
	server := stateServer("failed")
	defer server.Close()
	fastpollURL(t)

	err := pollURL(quickCtx(t, 5*time.Second), server.URL, "", "", "", false)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed")
}

func TestPoll_Timeout(t *testing.T) {
	server := stateServer("BUSY")
	defer server.Close()
	fastpollURL(t)

	err := pollURL(quickCtx(t, 200*time.Millisecond), server.URL, "", "", "", false)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "timeout")
}

func TestPoll_StatusField(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]any{"metadata": map[string]any{"status": "ACTIVE"}})
	}))
	defer server.Close()
	fastpollURL(t)

	err := pollURL(quickCtx(t, 5*time.Second), server.URL, "", "", "", false)
	assert.NoError(t, err)
}

func TestPoll_404_Delete(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"message":"not found"}`))
	}))
	defer server.Close()
	fastpollURL(t)

	err := pollURL(quickCtx(t, 5*time.Second), server.URL, "", "", "", true)
	assert.NoError(t, err)
}

func TestPoll_404_Create_Transient(t *testing.T) {
	// 404 during create is transient (resource provisioning). Should retry until timeout.
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()
	fastpollURL(t)

	err := pollURL(quickCtx(t, 200*time.Millisecond), server.URL, "", "", "", false)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "timeout")
}

func TestPoll_NoMetadataState_TreatedAsReady(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]any{"id": "abc", "properties": map[string]any{}})
	}))
	defer server.Close()
	fastpollURL(t)

	err := pollURL(quickCtx(t, 5*time.Second), server.URL, "", "", "", false)
	assert.NoError(t, err) // no state = resource doesn't track state, treat as ready
}

func TestPoll_NilMetadata_TreatedAsReady(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]any{"metadata": nil})
	}))
	defer server.Close()
	fastpollURL(t)

	err := pollURL(quickCtx(t, 5*time.Second), server.URL, "", "", "", false)
	assert.NoError(t, err)
}

func TestPoll_EmptyStateFields_TreatedAsReady(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]any{"metadata": map[string]any{"createdDate": "2024-01-01"}})
	}))
	defer server.Close()
	fastpollURL(t)

	err := pollURL(quickCtx(t, 5*time.Second), server.URL, "", "", "", false)
	assert.NoError(t, err)
}

func TestPoll_ServerError_FailsImmediately(t *testing.T) {
	var callCount int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&callCount, 1)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("server error"))
	}))
	defer server.Close()
	fastpollURL(t)

	err := pollURL(quickCtx(t, 5*time.Second), server.URL, "", "", "", false)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "server error (HTTP 500)")
	assert.Equal(t, int32(1), atomic.LoadInt32(&callCount), "should not retry 5xx")
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
	fastpollURL(t)

	err := pollURL(quickCtx(t, 5*time.Second), server.URL, "", "", "", false)
	assert.NoError(t, err)
}

func TestPoll_AuthHeaders_BearerToken(t *testing.T) {
	var gotAuth string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotAuth = r.Header.Get("Authorization")
		json.NewEncoder(w).Encode(map[string]any{"metadata": map[string]any{"state": "AVAILABLE"}})
	}))
	defer server.Close()
	fastpollURL(t)

	err := pollURL(quickCtx(t, 5*time.Second), server.URL, "my-token", "", "", false)
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
	fastpollURL(t)

	err := pollURL(quickCtx(t, 5*time.Second), server.URL, "", "user", "pass", false)
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
	fastpollURL(t)

	err := pollURL(quickCtx(t, 5*time.Second), server.URL, "tok", "user", "pass", false)
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
	fastpollURL(t)

	err := pollURL(quickCtx(t, 5*time.Second), server.URL, "", "", "", false)
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
	fastpollURL(t)

	Reset()
	captureHref(server.URL + "/datacenters/aaaaaaaa-1111-2222-3333-444444444444")

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
	fastpollURL(t)

	Reset()
	captureHref(server.URL + "/cloudapi/v6/datacenters/aaaaaaaa-1111-2222-3333-444444444444/servers/bbbbbbbb-1111-2222-3333-444444444444")

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
	fastpollURL(t)

	Reset()
	captureHref(server.URL + "/datacenters/aaaaaaaa-1111-2222-3333-444444444444/servers/bbbbbbbb-1111-2222-3333-444444444444")

	viper.Set(constants.ArgTimeout, 5)
	defer viper.Set(constants.ArgTimeout, 0)

	var buf bytes.Buffer
	err := WaitForAvailable(&buf, "", "", "")
	assert.NoError(t, err)
}

func TestWaitForAvailable_FailedResource(t *testing.T) {
	server := stateServer("FAILED")
	defer server.Close()
	fastpollURL(t)

	Reset()
	captureHref(server.URL + "/datacenters/aaaaaaaa-1111-2222-3333-444444444444")

	viper.Set(constants.ArgTimeout, 5)
	defer viper.Set(constants.ArgTimeout, 0)

	var buf bytes.Buffer
	err := WaitForAvailable(&buf, "", "", "")
	assert.NoError(t, err, "FAILED state should not cause error exit : command succeeded, resource just ended in bad state")
	assert.Contains(t, buf.String(), "FAILED", "should print warning about FAILED state")
}

func TestWaitForAvailable_NoStateField(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]any{"id": "abc"})
	}))
	defer server.Close()
	fastpollURL(t)

	Reset()
	captureHref(server.URL + "/clusters/aaaaaaaa-1111-2222-3333-444444444444")

	viper.Set(constants.ArgTimeout, 5)
	defer viper.Set(constants.ArgTimeout, 0)

	var buf bytes.Buffer
	err := WaitForAvailable(&buf, "", "", "")
	assert.NoError(t, err) // no state = treated as ready
}

// --- fetchResource ---

func TestFetchResource_Success(t *testing.T) {
	expected := map[string]any{"id": "abc", "properties": map[string]any{"name": "test"}}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(expected)
	}))
	defer server.Close()

	Reset()
	captureHref(server.URL + "/resource/abc")

	result, err := fetchResource("", "", "")
	assert.NoError(t, err)
	m, ok := result.(map[string]any)
	assert.True(t, ok)
	assert.Equal(t, "abc", m["id"])
}

func TestFetchResource_NoHref(t *testing.T) {
	Reset()
	_, err := fetchResource("", "", "")
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
	captureHref(server.URL + "/resource/abc")

	_, err := fetchResource("", "", "")
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
	captureHref(server.URL + "/resource/1")

	_, err := fetchResource("tok123", "", "")
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
	assert.Equal(t, server.URL+"/datacenters/abc", getHref())
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
	assert.Equal(t, server.URL+"/datacenters", getHref())
}

func TestWrapTransport_CapturesGET(t *testing.T) {
	Reset()
	viper.Set(constants.ArgWait, true)
	defer viper.Set(constants.ArgWait, false)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	hc := &http.Client{}
	WrapTransport(hc)

	url := server.URL + "/datacenters/dc-id/servers/srv-id"
	req, _ := http.NewRequest(http.MethodGet, url, nil)
	_, err := hc.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, url, getHref())
	assert.True(t, isGetOperation())
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
	assert.Empty(t, getHref()) // --wait not set
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
	assert.Equal(t, server.URL+"/test", getHref())
}

func TestWrapTransport_DoesNotOverwriteExistingHref(t *testing.T) {
	Reset()
	captureHref("https://already.set/resource/1")
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
	assert.Equal(t, "https://already.set/resource/1", getHref()) // not overwritten
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
	assert.Empty(t, getHref()) // error means no capture
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
			assert.Equal(t, server.URL+"/resource/"+method, getHref())
		})
	}
}

func TestCapturingTransport_SkipsReadMethods(t *testing.T) {
	for _, method := range []string{http.MethodHead, http.MethodOptions} {
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
			assert.Empty(t, getHref())
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
	fastpollURL(t)

	err := pollURL(quickCtx(t, 10*time.Second), server.URL, "", "", "", false)
	assert.NoError(t, err, "should succeed after rate limit clears")
	assert.Equal(t, int32(3), atomic.LoadInt32(&callCount),
		"should make 3 calls: 2 rate-limited + 1 AVAILABLE")
}

// TestPoll_429RateLimit_ContextCancellation verifies that a large Retry-After
// value does not block past the context deadline.
func TestPoll_429RateLimit_ContextCancellation(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Retry-After", "300") // 5 minutes, way past our timeout
		w.WriteHeader(http.StatusTooManyRequests)
	}))
	defer server.Close()
	fastpollURL(t)

	start := time.Now()
	err := pollURL(quickCtx(t, 500*time.Millisecond), server.URL, "", "", "", false)
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
			fastpollURL(t)

			err := pollURL(quickCtx(t, 200*time.Millisecond), server.URL, "", "", "", false)
			assert.NoError(t, err, "%q should be treated as terminal success", state)
		})
	}

	// ERROR is always a terminal failure state
	t.Run("ERROR_failure", func(t *testing.T) {
		server := stateServer("ERROR")
		defer server.Close()
		fastpollURL(t)

		err := pollURL(quickCtx(t, 200*time.Millisecond), server.URL, "", "", "", false)
		assert.Error(t, err, "ERROR should be treated as terminal failure")
		assert.Contains(t, err.Error(), "ERROR")
	})

	// DESTROYING is a failure state for non-delete operations
	t.Run("DESTROYING_failure_non_delete", func(t *testing.T) {
		server := stateServer("DESTROYING")
		defer server.Close()
		fastpollURL(t)

		err := pollURL(quickCtx(t, 200*time.Millisecond), server.URL, "", "", "", false)
		assert.Error(t, err, "DESTROYING should be treated as terminal failure for non-delete")
		assert.Contains(t, err.Error(), "DESTROYING")
	})
}

// TestPoll_Destroying_Delete_ContinuesPolling verifies that DESTROYING is treated
// as transient during delete operations - the poller keeps going until 404.
func TestPoll_Destroying_Delete_ContinuesPolling(t *testing.T) {
	callCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callCount++
		if callCount <= 2 {
			// First two polls return DESTROYING
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]any{
				"metadata": map[string]any{"state": "DESTROYING"},
			})
			return
		}
		// Third poll returns 404 - resource gone
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()
	fastpollURL(t)

	err := pollURL(quickCtx(t, 2*time.Second), server.URL, "", "", "", true)
	assert.NoError(t, err, "DESTROYING during delete should be transient, 404 should succeed")
	assert.GreaterOrEqual(t, callCount, 3)
}

// TestPoll_TransientError_ThenNoState_ContinuesPolling verifies that after a
// transient error, an empty-state response does NOT cause early exit.
func TestPoll_TransientError_ThenNoState_ContinuesPolling(t *testing.T) {
	callCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callCount++
		switch callCount {
		case 1:
			// First poll: malformed JSON (transient error)
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("not json"))
		case 2:
			// Second poll: valid response with no state field
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]any{"id": "test"})
		default:
			// Third poll: AVAILABLE
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]any{
				"metadata": map[string]any{"state": "AVAILABLE"},
			})
		}
	}))
	defer server.Close()
	fastpollURL(t)

	err := pollURL(quickCtx(t, 2*time.Second), server.URL, "", "", "", false)
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, callCount, 3, "should not early-exit after transient error + empty state")
}

// TestPoll_FirstSuccess_NoState_ExitsEarly verifies that if the first poll
// succeeds with no state field, the poller exits immediately (resource has no state tracking).
func TestPoll_FirstSuccess_NoState_ExitsEarly(t *testing.T) {
	callCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callCount++
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]any{"id": "test"})
	}))
	defer server.Close()
	fastpollURL(t)

	err := pollURL(quickCtx(t, 2*time.Second), server.URL, "", "", "", false)
	assert.NoError(t, err)
	assert.Equal(t, 1, callCount, "should exit after first successful poll with no state")
}

// TestWaitForAvailable_NoTargets_Warning verifies a warning is emitted when
// --wait is active but no resource URL could be determined for polling.
// This happens when an action endpoint is captured but no Location header is returned.
func TestWaitForAvailable_NoTargets_Warning(t *testing.T) {
	Reset()
	fastpollURL(t)
	viper.Set(constants.ArgWait, true)
	defer viper.Set(constants.ArgWait, false)
	viper.Set(constants.ArgTimeout, 10)
	defer viper.Set(constants.ArgTimeout, 0)

	// Simulate action endpoint captured with no Location header
	captureRequestURL(http.MethodPost, "https://api.ionos.com/cloudapi/v6/datacenters/dc-id/servers/srv-id/start", "")

	var buf bytes.Buffer
	err := WaitForAvailable(&buf, "", "", "")
	assert.NoError(t, err)
	assert.Contains(t, buf.String(), "Warning: --wait active but no resource URL could be determined")
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

func TestIsActionEndpoint(t *testing.T) {
	tests := []struct {
		href   string
		expect bool
	}{
		{"https://api.ionos.com/cloudapi/v6/datacenters/dc/servers/srv/start", true},
		{"https://api.ionos.com/cloudapi/v6/datacenters/dc/servers/srv/stop", true},
		{"https://api.ionos.com/cloudapi/v6/datacenters/dc/servers/srv/reboot", true},
		{"https://api.ionos.com/cloudapi/v6/datacenters/dc/servers/srv/suspend", true},
		{"https://api.ionos.com/cloudapi/v6/datacenters/dc/servers/srv/resume", true},
		{"https://api.ionos.com/cloudapi/v6/clusters/c1/restore", true},
		{"https://api.ionos.com/dns/v1/zones/z1/transfer", true},
		// Query params must be stripped before checking
		{"https://api.ionos.com/cloudapi/v6/datacenters/dc/servers/srv/start?depth=1&limit=50", true},
		// Trailing slash
		{"https://api.ionos.com/cloudapi/v6/datacenters/dc/servers/srv/stop/", true},
		// Non-action endpoints
		{"https://api.ionos.com/cloudapi/v6/datacenters/dc/servers/srv", false},
		{"https://api.ionos.com/cloudapi/v6/datacenters/dc", false},
		{"https://api.ionos.com/cloudapi/v6/datacenters/dc/servers/srv/create", false},
		{"", false},
		{"http://[invalid", false},
	}
	for _, tt := range tests {
		t.Run(tt.href, func(t *testing.T) {
			assert.Equal(t, tt.expect, isActionEndpoint(tt.href))
		})
	}
}

func TestExtractID(t *testing.T) {
	tests := []struct {
		name     string
		input    any
		expected string
	}{
		{"single resource with id", map[string]any{"id": "abc-123", "type": "server"}, "abc-123"},
		{"list response with items", map[string]any{"id": "col", "items": []any{}}, ""},
		{"no id field", map[string]any{"href": "/resource"}, ""},
		{"nil input", nil, ""},
		{"non-map input (string)", "not a map", ""},
		{"non-map input (int)", 42, ""},
		{"id is non-string", map[string]any{"id": 123}, ""},
		{"nested struct with id", struct {
			ID   string `json:"id"`
			Name string `json:"name"`
		}{ID: "struct-id", Name: "test"}, "struct-id"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, extractID(tt.input))
		})
	}
}

func TestIsDeleteOperation(t *testing.T) {
	t.Run("DELETE method", func(t *testing.T) {
		Reset()
		captureRequestURL(http.MethodDelete, "https://api/resource", "")
		assert.True(t, isDeleteOperation())
	})
	t.Run("POST method", func(t *testing.T) {
		Reset()
		captureRequestURL(http.MethodPost, "https://api/resource", "")
		assert.False(t, isDeleteOperation())
	})
	t.Run("no capture", func(t *testing.T) {
		Reset()
		assert.False(t, isDeleteOperation())
	})
}

func TestGetRequestStatusURL(t *testing.T) {
	t.Run("with Location header", func(t *testing.T) {
		Reset()
		captureRequestURL(http.MethodPost, "https://api/resource", "https://api/requests/123/status")
		assert.Equal(t, "https://api/requests/123/status", getRequestStatusURL())
	})
	t.Run("no Location header", func(t *testing.T) {
		Reset()
		captureRequestURL(http.MethodPost, "https://api/resource", "")
		assert.Empty(t, getRequestStatusURL())
	})
	t.Run("no capture", func(t *testing.T) {
		Reset()
		assert.Empty(t, getRequestStatusURL())
	})
}

func TestIsStructuredOutput(t *testing.T) {
	for _, tt := range []struct {
		output string
		expect bool
	}{
		{"json", true},
		{"api-json", true},
		{"text", false},
		{"", false},
	} {
		t.Run(tt.output, func(t *testing.T) {
			viper.Set(constants.ArgOutput, tt.output)
			defer viper.Set(constants.ArgOutput, "")
			assert.Equal(t, tt.expect, isStructuredOutput())
		})
	}
}

func TestReset_AllFields(t *testing.T) {
	// Set all state
	captureHref("https://api/resource")
	captureRequestURL(http.MethodDelete, "https://api/resource", "https://api/requests/1/status")
	setRerendering(true)

	// Set sdkTransport via WrapTransport
	hc := &http.Client{Transport: http.DefaultTransport}
	viper.Set(constants.ArgWait, true)
	WrapTransport(hc)
	viper.Set(constants.ArgWait, false)

	// Verify state is set
	assert.NotEmpty(t, getHref())
	assert.NotEmpty(t, getRequestStatusURL())
	assert.True(t, isDeleteOperation())
	assert.True(t, isRerendering())

	mu.Lock()
	assert.NotNil(t, sdkTransport)
	mu.Unlock()

	// Reset and verify all cleared
	Reset()

	assert.Empty(t, getHref())
	assert.Empty(t, getRequestStatusURL())
	assert.False(t, isDeleteOperation())
	assert.False(t, isGetOperation())
	assert.False(t, isRerendering())
	r, cols := getRerenderable()
	assert.Nil(t, r)
	assert.Nil(t, cols)

	mu.Lock()
	assert.Nil(t, sdkTransport)
	assert.False(t, lastHrefFromGet)
	mu.Unlock()
}

func TestCaptureRequestURL_LocationHeader(t *testing.T) {
	Reset()
	captureRequestURL(http.MethodPost, "https://api/resource", "https://api/requests/123/status")
	assert.Equal(t, "https://api/requests/123/status", getRequestStatusURL())
}

func TestCaptureRequestURL_MethodTracking(t *testing.T) {
	t.Run("DELETE tracked", func(t *testing.T) {
		Reset()
		captureRequestURL(http.MethodDelete, "https://api/resource", "")
		assert.True(t, isDeleteOperation())
	})
	t.Run("PATCH not delete", func(t *testing.T) {
		Reset()
		captureRequestURL(http.MethodPatch, "https://api/resource", "")
		assert.False(t, isDeleteOperation())
	})
}

func TestCaptureRequestURL_MultipleCalls(t *testing.T) {
	Reset()
	captureRequestURL(http.MethodDelete, "https://api/resource/1", "https://api/requests/1/status")
	captureRequestURL(http.MethodDelete, "https://api/resource/2", "https://api/requests/2/status")
	captureRequestURL(http.MethodDelete, "https://api/resource/3", "https://api/requests/3/status")

	// First href wins (href only captured when empty)
	assert.Equal(t, "https://api/resource/1", getHref())
	// Last Location header wins (always overwritten)
	assert.Equal(t, "https://api/requests/3/status", getRequestStatusURL())
	assert.True(t, isDeleteOperation())
}

func TestWaitForAvailable_PollsRequestStatusFirst(t *testing.T) {
	Reset()
	fastpollURL(t)

	var callOrder []string
	var mu2 sync.Mutex

	reqStatusServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mu2.Lock()
		callOrder = append(callOrder, "requestStatus")
		mu2.Unlock()
		json.NewEncoder(w).Encode(map[string]any{"metadata": map[string]any{"status": "DONE"}})
	}))
	defer reqStatusServer.Close()

	resourceServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mu2.Lock()
		callOrder = append(callOrder, "resource")
		mu2.Unlock()
		json.NewEncoder(w).Encode(map[string]any{"metadata": map[string]any{"state": "AVAILABLE"}})
	}))
	defer resourceServer.Close()

	captureHref(resourceServer.URL)
	captureRequestURL(http.MethodPost, resourceServer.URL, reqStatusServer.URL)

	viper.Set(constants.ArgWait, true)
	viper.Set(constants.ArgTimeout, 5)
	defer viper.Set(constants.ArgWait, false)
	defer viper.Set(constants.ArgTimeout, 0)

	var buf bytes.Buffer
	err := WaitForAvailable(&buf, "", "", "")
	assert.NoError(t, err)

	mu2.Lock()
	defer mu2.Unlock()
	assert.GreaterOrEqual(t, len(callOrder), 2)
	assert.Equal(t, "requestStatus", callOrder[0], "request status should be polled first")
}

func TestWaitForAvailable_ActionEndpoint_SkipsResourcepollURL(t *testing.T) {
	Reset()
	fastpollURL(t)

	var resourceCalls int32
	resourceServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&resourceCalls, 1)
		json.NewEncoder(w).Encode(map[string]any{"metadata": map[string]any{"state": "AVAILABLE"}})
	}))
	defer resourceServer.Close()

	reqStatusServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]any{"metadata": map[string]any{"status": "DONE"}})
	}))
	defer reqStatusServer.Close()

	// href is an action endpoint URL (ends with /start)
	actionURL := "https://api.ionos.com/cloudapi/v6/datacenters/dc/servers/srv/start"
	captureHref(actionURL)
	captureRequestURL(http.MethodPost, actionURL, reqStatusServer.URL)

	viper.Set(constants.ArgWait, true)
	viper.Set(constants.ArgTimeout, 5)
	defer viper.Set(constants.ArgWait, false)
	defer viper.Set(constants.ArgTimeout, 0)

	var buf bytes.Buffer
	err := WaitForAvailable(&buf, "", "", "")
	assert.NoError(t, err)

	// Resource server should NOT be called - action endpoints skip resource polling
	assert.Equal(t, int32(0), atomic.LoadInt32(&resourceCalls),
		"action endpoint should only poll request status, not resource")
}

func TestWaitForAvailable_OnlyRequestStatusURL_NoHref(t *testing.T) {
	Reset()
	// Only Location header captured, no resource href
	mu.Lock()
	lastRequestURL = "https://api/requests/123/status"
	mu.Unlock()

	// getHref() is empty so WaitForAvailable returns nil immediately
	var buf bytes.Buffer
	err := WaitForAvailable(&buf, "", "", "")
	assert.NoError(t, err, "no href means early return, no polling")
}

func TestPoll_Unauthorized_NoRetry(t *testing.T) {
	var callCount int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&callCount, 1)
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{"message":"unauthorized"}`))
	}))
	defer server.Close()
	fastpollURL(t)

	err := pollURL(quickCtx(t, 5*time.Second), server.URL, "", "", "", false)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "authentication failed")
	assert.Equal(t, int32(1), atomic.LoadInt32(&callCount), "should not retry 401")
}

func TestPoll_Forbidden_NoRetry(t *testing.T) {
	var callCount int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&callCount, 1)
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte(`{"message":"forbidden"}`))
	}))
	defer server.Close()
	fastpollURL(t)

	err := pollURL(quickCtx(t, 5*time.Second), server.URL, "", "", "", false)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "authentication failed")
	assert.Equal(t, int32(1), atomic.LoadInt32(&callCount), "should not retry 403")
}

func TestPoll_ClientErrors_NoRetry(t *testing.T) {
	for _, code := range []int{405, 409, 410, 422} {
		t.Run(fmt.Sprintf("HTTP_%d", code), func(t *testing.T) {
			var callCount int32
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				atomic.AddInt32(&callCount, 1)
				w.WriteHeader(code)
				w.Write([]byte(`{"message":"error"}`))
			}))
			defer server.Close()
			fastpollURL(t)

			err := pollURL(quickCtx(t, 5*time.Second), server.URL, "", "", "", false)
			assert.Error(t, err)
			assert.Contains(t, err.Error(), "client error")
			assert.Equal(t, int32(1), atomic.LoadInt32(&callCount), "should not retry %d", code)
		})
	}
}

func TestPoll_ContextCancellation(t *testing.T) {
	server := stateServer("BUSY")
	defer server.Close()
	fastpollURL(t)

	start := time.Now()
	err := pollURL(quickCtx(t, 200*time.Millisecond), server.URL, "", "", "", false)
	elapsed := time.Since(start)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "timeout")
	assert.Less(t, elapsed, 2*time.Second, "should return promptly after context cancellation")
}

func TestSetHeaders_ContractNumber(t *testing.T) {
	t.Run("with env var", func(t *testing.T) {
		t.Setenv("IONOS_CONTRACT_NUMBER", "12345")
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "12345", r.Header.Get("X-Contract-Number"))
			json.NewEncoder(w).Encode(map[string]any{"metadata": map[string]any{"state": "AVAILABLE"}})
		}))
		defer server.Close()
		fastpollURL(t)

		err := pollURL(quickCtx(t, 5*time.Second), server.URL, "", "", "", false)
		assert.NoError(t, err)
	})

	t.Run("without env var", func(t *testing.T) {
		t.Setenv("IONOS_CONTRACT_NUMBER", "")
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Empty(t, r.Header.Get("X-Contract-Number"))
			json.NewEncoder(w).Encode(map[string]any{"metadata": map[string]any{"state": "AVAILABLE"}})
		}))
		defer server.Close()
		fastpollURL(t)

		err := pollURL(quickCtx(t, 5*time.Second), server.URL, "", "", "", false)
		assert.NoError(t, err)
	})
}

func TestAppendDepthParam(t *testing.T) {
	tests := []struct {
		name     string
		depth    int
		url      string
		contains string
	}{
		{"depth=0 defaults to 1", 0, "https://api.ionos.com/test", "depth=1"},
		{"depth=1", 1, "https://api.ionos.com/test", "depth=1"},
		{"depth=3", 3, "https://api.ionos.com/test", "depth=3"},
		{"negative defaults to 1", -1, "https://api.ionos.com/test", "depth=1"},
		{"overwrites existing depth", 2, "https://api.ionos.com/test?depth=5", "depth=2"},
		{"preserves other params", 1, "https://api.ionos.com/test?pretty=true", "pretty=true"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			viper.Set(constants.FlagDepth, tt.depth)
			defer viper.Set(constants.FlagDepth, 0)

			result := appendDepthParam(tt.url)
			assert.Contains(t, result, tt.contains)
		})
	}

	// Verify overwrite doesn't keep old value
	t.Run("no stale depth=5", func(t *testing.T) {
		viper.Set(constants.FlagDepth, 2)
		defer viper.Set(constants.FlagDepth, 0)
		result := appendDepthParam("https://api.ionos.com/test?depth=5")
		assert.NotContains(t, result, "depth=5")
	})
}

func TestWaitForAvailable_ProgressOutput_Done(t *testing.T) {
	Reset()
	fastpollURL(t)
	server := stateServer("AVAILABLE")
	defer server.Close()

	captureHref(server.URL)
	viper.Set(constants.ArgWait, true)
	viper.Set(constants.ArgTimeout, 5)
	viper.Set(constants.ArgOutput, "text")
	defer func() {
		viper.Set(constants.ArgWait, false)
		viper.Set(constants.ArgTimeout, 0)
		viper.Set(constants.ArgOutput, "")
	}()

	var buf bytes.Buffer
	err := WaitForAvailable(&buf, "", "", "")
	assert.NoError(t, err)
	assert.Contains(t, buf.String(), "DONE")
}

func TestWaitForAvailable_ProgressOutput_Failed(t *testing.T) {
	Reset()
	fastpollURL(t)
	server := stateServer("FAILED")
	defer server.Close()

	captureHref(server.URL)
	viper.Set(constants.ArgWait, true)
	viper.Set(constants.ArgTimeout, 5)
	viper.Set(constants.ArgOutput, "text")
	defer func() {
		viper.Set(constants.ArgWait, false)
		viper.Set(constants.ArgTimeout, 0)
		viper.Set(constants.ArgOutput, "")
	}()

	var buf bytes.Buffer
	err := WaitForAvailable(&buf, "", "", "")
	assert.NoError(t, err, "FAILED state is a warning, not an error")
	assert.Contains(t, buf.String(), "FAILED")
}

func TestWaitForAvailable_JsonOutput_Silent(t *testing.T) {
	Reset()
	fastpollURL(t)
	server := stateServer("AVAILABLE")
	defer server.Close()

	captureHref(server.URL)
	viper.Set(constants.ArgWait, true)
	viper.Set(constants.ArgTimeout, 5)
	viper.Set(constants.ArgOutput, "json")
	defer func() {
		viper.Set(constants.ArgWait, false)
		viper.Set(constants.ArgTimeout, 0)
		viper.Set(constants.ArgOutput, "")
	}()

	var buf bytes.Buffer
	err := WaitForAvailable(&buf, "", "", "")
	assert.NoError(t, err)
	assert.Empty(t, buf.String(), "JSON mode should produce no progress output")
}

func TestWaitForAvailable_TimeoutZeroWarning(t *testing.T) {
	Reset()
	fastpollURL(t)
	server := stateServer("AVAILABLE")
	defer server.Close()

	captureHref(server.URL)
	viper.Set(constants.ArgWait, true)
	viper.Set(constants.ArgTimeout, 0)
	viper.Set(constants.ArgOutput, "text")
	defer func() {
		viper.Set(constants.ArgWait, false)
		viper.Set(constants.ArgTimeout, 0)
		viper.Set(constants.ArgOutput, "")
	}()

	var buf bytes.Buffer
	err := WaitForAvailable(&buf, "", "", "")
	assert.NoError(t, err)
	assert.Contains(t, buf.String(), "Warning: --timeout 0")
}

func TestWrapTransport_CapturesSdkTransport(t *testing.T) {
	Reset()
	customTransport := &http.Transport{}
	hc := &http.Client{Transport: customTransport}

	WrapTransport(hc)

	mu.Lock()
	captured := sdkTransport
	mu.Unlock()

	assert.Equal(t, customTransport, captured,
		"sdkTransport should be the original transport, not the wrapper")

	// Verify the client's transport IS the wrapper
	_, isCapturing := hc.Transport.(*capturingTransport)
	assert.True(t, isCapturing, "client transport should be wrapped")
}

func TestNewPoller_UsesSdkTransport(t *testing.T) {
	Reset()
	customTransport := &http.Transport{}

	mu.Lock()
	sdkTransport = customTransport
	mu.Unlock()

	p := newPoller("", "", "")
	assert.Equal(t, customTransport, p.client.Transport,
		"poller should reuse sdkTransport")
}

func TestNewPoller_FallsBackToDefault(t *testing.T) {
	Reset() // clears sdkTransport

	p := newPoller("", "", "")
	assert.Equal(t, http.DefaultTransport, p.client.Transport,
		"poller should fall back to http.DefaultTransport when sdkTransport is nil")
}

func TestFetchResource_MalformedJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("{invalid json"))
	}))
	defer server.Close()

	Reset()
	captureHref(server.URL)

	_, err := fetchResource("", "", "")
	assert.Error(t, err)
}

func TestFetchResource_HTTPError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	Reset()
	captureHref(server.URL)

	_, err := fetchResource("", "", "")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "HTTP 500")
}

// --- Re-render flow tests ---

func TestRerender_FetchResource_Success(t *testing.T) {
	expected := map[string]any{"id": "test-id", "metadata": map[string]any{"state": "AVAILABLE"}}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(expected)
	}))
	defer server.Close()

	Reset()
	captureHref(server.URL)

	data, err := fetchResource("", "", "")
	assert.NoError(t, err)
	assert.NotNil(t, data)
	// Verify it's a parsed JSON map
	m, ok := data.(map[string]any)
	assert.True(t, ok)
	assert.Equal(t, "test-id", m["id"])
}

func TestRerender_CaptureAndRetrieve(t *testing.T) {
	Reset()
	mock := &mockRerenderable{}
	cols := []string{"Id", "Name", "State"}

	captureRerenderable(mock, cols)

	r, gotCols := getRerenderable()
	assert.Equal(t, mock, r)
	assert.Equal(t, cols, gotCols)
}

func TestRerender_SetIsRerendering(t *testing.T) {
	Reset()
	assert.False(t, isRerendering())

	setRerendering(true)
	assert.True(t, isRerendering())

	setRerendering(false)
	assert.False(t, isRerendering())
}

func TestRerender_MockExtractAndRender(t *testing.T) {
	mock := &mockRerenderable{}
	data := map[string]any{"id": "abc", "metadata": map[string]any{"state": "AVAILABLE"}}

	err := mock.Extract(data)
	assert.NoError(t, err)
	assert.True(t, mock.extractCalled)
	assert.Equal(t, data, mock.data)

	out, err := mock.Render([]string{"Id", "State"})
	assert.NoError(t, err)
	assert.True(t, mock.renderCalled)
	assert.Contains(t, out, "rendered:")
}

func TestRerender_ExtractError(t *testing.T) {
	mock := &mockRerenderable{extractErr: fmt.Errorf("extract failed")}
	err := mock.Extract(map[string]any{})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "extract failed")
}

func TestRerender_RenderError(t *testing.T) {
	mock := &mockRerenderable{renderErr: fmt.Errorf("render failed")}
	out, err := mock.Render(nil)
	assert.Error(t, err)
	assert.Empty(t, out)
}

// --- Non-compute URL parsing tests ---

func TestResourceAndParentURLs_DBaaS_Postgres(t *testing.T) {
	// Top-level DBaaS cluster: no parent (databases/postgresql are API prefix, not resources)
	urls := resourceAndParentURLs("https://api.ionos.com/databases/postgresql/clusters/aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee")
	assert.Equal(t, []string{
		"https://api.ionos.com/databases/postgresql/clusters/aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee",
	}, urls)
}

func TestResourceAndParentURLs_DBaaS_Mongo(t *testing.T) {
	urls := resourceAndParentURLs("https://api.ionos.com/databases/mongodb/clusters/aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee")
	assert.Equal(t, []string{
		"https://api.ionos.com/databases/mongodb/clusters/aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee",
	}, urls)
}

func TestResourceAndParentURLs_DBaaS_Nested(t *testing.T) {
	// Nested DBaaS resource: cluster is the parent
	urls := resourceAndParentURLs("https://api.ionos.com/databases/postgresql/clusters/aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee/users/11111111-2222-3333-4444-555555555555")
	assert.Equal(t, []string{
		"https://api.ionos.com/databases/postgresql/clusters/aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee/users/11111111-2222-3333-4444-555555555555",
		"https://api.ionos.com/databases/postgresql/clusters/aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee",
	}, urls)
}

func TestResourceAndParentURLs_DNS_Record(t *testing.T) {
	urls := resourceAndParentURLs("https://dns.de-fra.ionos.com/zones/aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee/records/11111111-2222-3333-4444-555555555555")
	assert.Equal(t, []string{
		"https://dns.de-fra.ionos.com/zones/aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee/records/11111111-2222-3333-4444-555555555555",
		"https://dns.de-fra.ionos.com/zones/aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee",
	}, urls)
}

// --- Bulk delete capture test ---

func TestCaptureRequestURL_BulkDelete_LastRequestWins(t *testing.T) {
	Reset()

	captureRequestURL(http.MethodDelete, "https://api/resource/1", "https://api/requests/loc1/status")
	captureRequestURL(http.MethodDelete, "https://api/resource/2", "https://api/requests/loc2/status")
	captureRequestURL(http.MethodDelete, "https://api/resource/3", "https://api/requests/loc3/status")

	// First href wins (captured from first call)
	assert.Equal(t, "https://api/resource/1", getHref())
	// Last Location header wins (overwritten each call)
	assert.Equal(t, "https://api/requests/loc3/status", getRequestStatusURL())
	// Method is still DELETE
	assert.True(t, isDeleteOperation())
}

// --- GET capture priority tests ---

func TestIsGetOperation(t *testing.T) {
	t.Run("GET method", func(t *testing.T) {
		Reset()
		captureGetURL("https://api/resource")
		assert.True(t, isGetOperation())
	})
	t.Run("POST method", func(t *testing.T) {
		Reset()
		captureRequestURL(http.MethodPost, "https://api/resource", "")
		assert.False(t, isGetOperation())
	})
	t.Run("no capture", func(t *testing.T) {
		Reset()
		assert.False(t, isGetOperation())
	})
}

func TestCaptureGetURL_DoesNotOverwriteMutating(t *testing.T) {
	Reset()
	// POST captures first
	captureRequestURL(http.MethodPost, "https://api/servers", "https://api/requests/1/status")
	// GET fires after (e.g. completer or validator)
	captureGetURL("https://api/datacenters/dc-id")

	// POST href and method should be preserved
	assert.Equal(t, "https://api/servers", getHref())
	assert.False(t, isGetOperation())
	assert.Equal(t, "https://api/requests/1/status", getRequestStatusURL())
}

func TestCaptureRequestURL_OverwritesGetHref(t *testing.T) {
	Reset()
	// GET captures first (e.g. PreCmdRun validator)
	captureGetURL("https://api/datacenters/dc-id")
	assert.Equal(t, "https://api/datacenters/dc-id", getHref())
	assert.True(t, isGetOperation())

	// POST fires for actual command
	captureRequestURL(http.MethodPost, "https://api/servers", "https://api/requests/1/status")

	// POST should overwrite GET-captured href
	assert.Equal(t, "https://api/servers", getHref())
	assert.False(t, isGetOperation())
}

func TestCaptureGetURL_SetByGetThenCaptureHrefOverwrites(t *testing.T) {
	Reset()
	captureGetURL("https://api/resource/1")
	assert.True(t, isGetOperation())

	// BeforeRender extracts href from response body
	captureHref("https://api/resource/1?depth=1")

	assert.Equal(t, "https://api/resource/1?depth=1", getHref())
	// Method still GET since captureHref doesn't change method
	assert.True(t, isGetOperation())
}
