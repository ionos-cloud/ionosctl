package client

import (
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// mockRoundTripper is a simple mock that returns a preconfigured response.
type mockRoundTripper struct {
	resp *http.Response
	err  error
}

func (m *mockRoundTripper) RoundTrip(_ *http.Request) (*http.Response, error) {
	return m.resp, m.err
}

// newXMLResponse builds an *http.Response with the given XML body and content type.
func newXMLResponse(body string) *http.Response {
	return &http.Response{
		StatusCode:    200,
		Header:        http.Header{"Content-Type": []string{"application/xml"}},
		Body:          io.NopCloser(strings.NewReader(body)),
		ContentLength: int64(len(body)),
	}
}

func TestOwnerIDFixTransport_RewritesAnonymousOwner(t *testing.T) {
	input := `<ListBucketResult><Contents><Owner><ID>anonymous</ID></Owner></Contents></ListBucketResult>`
	expected := `<ListBucketResult><Contents><Owner><ID>0</ID></Owner></Contents></ListBucketResult>`

	transport := &ownerIDFixTransport{
		base: &mockRoundTripper{resp: newXMLResponse(input)},
	}

	resp, err := transport.RoundTrip(&http.Request{})
	if !assert.NoError(t, err) {
		return
	}

	body, err := io.ReadAll(resp.Body)
	if !assert.NoError(t, err) {
		return
	}
	assert.Equal(t, expected, string(body))
}

func TestOwnerIDFixTransport_PreservesNumericOwner(t *testing.T) {
	input := `<ListBucketResult><Contents><Owner><ID>12345</ID></Owner></Contents></ListBucketResult>`

	transport := &ownerIDFixTransport{
		base: &mockRoundTripper{resp: newXMLResponse(input)},
	}

	resp, err := transport.RoundTrip(&http.Request{})
	if !assert.NoError(t, err) {
		return
	}

	body, err := io.ReadAll(resp.Body)
	if !assert.NoError(t, err) {
		return
	}

	// Pure numeric IDs contain no \D match, so the regex should NOT match.
	assert.Equal(t, input, string(body), "numeric Owner ID should not be rewritten")
}

func TestOwnerIDFixTransport_SkipsNonXMLResponse(t *testing.T) {
	jsonBody := `{"owner":{"id":"anonymous"}}`
	resp := &http.Response{
		StatusCode:    200,
		Header:        http.Header{"Content-Type": []string{"application/json"}},
		Body:          io.NopCloser(strings.NewReader(jsonBody)),
		ContentLength: int64(len(jsonBody)),
	}

	transport := &ownerIDFixTransport{
		base: &mockRoundTripper{resp: resp},
	}

	got, err := transport.RoundTrip(&http.Request{})
	if !assert.NoError(t, err) {
		return
	}

	body, err := io.ReadAll(got.Body)
	if !assert.NoError(t, err) {
		return
	}
	assert.Equal(t, jsonBody, string(body), "non-XML response body should be unchanged")
}

func TestOwnerIDFixTransport_HandlesNilBody(t *testing.T) {
	resp := &http.Response{
		StatusCode: 204,
		Header:     http.Header{"Content-Type": []string{"application/xml"}},
		Body:       nil,
	}

	transport := &ownerIDFixTransport{
		base: &mockRoundTripper{resp: resp},
	}

	got, err := transport.RoundTrip(&http.Request{})
	if !assert.NoError(t, err) {
		return
	}
	assert.Nil(t, got.Body, "nil body should remain nil")
}

func TestOwnerIDFixTransport_MultipleOwners(t *testing.T) {
	input := `<ListBucketResult>` +
		`<Contents><Owner><ID>anonymous</ID></Owner></Contents>` +
		`<Contents><Owner><ID>user-abc</ID></Owner></Contents>` +
		`</ListBucketResult>`
	expected := `<ListBucketResult>` +
		`<Contents><Owner><ID>0</ID></Owner></Contents>` +
		`<Contents><Owner><ID>0</ID></Owner></Contents>` +
		`</ListBucketResult>`

	transport := &ownerIDFixTransport{
		base: &mockRoundTripper{resp: newXMLResponse(input)},
	}

	resp, err := transport.RoundTrip(&http.Request{})
	if !assert.NoError(t, err) {
		return
	}

	body, err := io.ReadAll(resp.Body)
	if !assert.NoError(t, err) {
		return
	}
	assert.Equal(t, expected, string(body))
}

func TestOwnerIDFixTransport_UpdatesContentLength(t *testing.T) {
	// "anonymous" (9 chars) gets replaced with "0" (1 char), so the body shrinks.
	input := `<Owner><ID>anonymous</ID></Owner>`
	expected := `<Owner><ID>0</ID></Owner>`

	transport := &ownerIDFixTransport{
		base: &mockRoundTripper{resp: newXMLResponse(input)},
	}

	resp, err := transport.RoundTrip(&http.Request{})
	if !assert.NoError(t, err) {
		return
	}

	body, err := io.ReadAll(resp.Body)
	if !assert.NoError(t, err) {
		return
	}
	assert.Equal(t, int64(len(expected)), resp.ContentLength,
		"ContentLength should match actual rewritten body length")
	assert.Equal(t, expected, string(body))
}

func TestOwnerIDFixTransport_UsesDefaultTransportWhenBaseNil(t *testing.T) {
	// With a nil base, RoundTrip should use http.DefaultTransport under the hood.
	// We verify this by temporarily replacing http.DefaultTransport with our mock,
	// then restoring it after the test.
	originalInput := `<Owner><ID>anonymous</ID></Owner>`
	expectedOutput := `<Owner><ID>0</ID></Owner>`

	mockTransport := &mockRoundTripper{resp: newXMLResponse(originalInput)}

	originalDefault := http.DefaultTransport
	http.DefaultTransport = mockTransport
	defer func() { http.DefaultTransport = originalDefault }()

	transport := &ownerIDFixTransport{base: nil}

	resp, err := transport.RoundTrip(&http.Request{})
	if !assert.NoError(t, err) {
		return
	}

	body, err := io.ReadAll(resp.Body)
	if !assert.NoError(t, err) {
		return
	}
	assert.Equal(t, expectedOutput, string(body),
		"nil base should fall back to http.DefaultTransport and still rewrite")
}

func TestOwnerIDFixTransport_PreservesNonOwnerXML(t *testing.T) {
	// The regex should only match <Owner><ID>...</ID> patterns, not other elements.
	input := `<ListBucketResult>` +
		`<Name>my-bucket</Name>` +
		`<ID>some-request-id</ID>` +
		`<Contents><Key>file.txt</Key><Owner><ID>anonymous</ID></Owner></Contents>` +
		`</ListBucketResult>`
	expected := `<ListBucketResult>` +
		`<Name>my-bucket</Name>` +
		`<ID>some-request-id</ID>` +
		`<Contents><Key>file.txt</Key><Owner><ID>0</ID></Owner></Contents>` +
		`</ListBucketResult>`

	transport := &ownerIDFixTransport{
		base: &mockRoundTripper{resp: newXMLResponse(input)},
	}

	resp, err := transport.RoundTrip(&http.Request{})
	if !assert.NoError(t, err) {
		return
	}

	body, err := io.ReadAll(resp.Body)
	if !assert.NoError(t, err) {
		return
	}
	assert.Equal(t, expected, string(body),
		"only <Owner><ID>...</ID> should be rewritten; other <ID> elements must be preserved")
}

func TestOwnerIDFixTransport_OwnerWithWhitespace(t *testing.T) {
	// The regex allows \s* between <Owner> and <ID>, test that whitespace is handled.
	input := "<Owner>\n  <ID>anonymous</ID>\n</Owner>"
	expected := "<Owner>\n  <ID>0</ID>\n</Owner>"

	transport := &ownerIDFixTransport{
		base: &mockRoundTripper{resp: newXMLResponse(input)},
	}

	resp, err := transport.RoundTrip(&http.Request{})
	if !assert.NoError(t, err) {
		return
	}

	body, err := io.ReadAll(resp.Body)
	if !assert.NoError(t, err) {
		return
	}
	assert.Equal(t, expected, string(body))
}

func TestOwnerIDRe_MatchesNonNumeric(t *testing.T) {
	// Direct regex tests to verify matching behavior precisely.
	tests := []struct {
		name    string
		input   string
		matches bool
		output  string
	}{
		{"pure alpha", "<Owner><ID>anonymous</ID>", true, "<Owner><ID>0</ID>"},
		{"alphanumeric", "<Owner><ID>user123</ID>", true, "<Owner><ID>0</ID>"},
		{"uuid-like", "<Owner><ID>550e8400-e29b-41d4-a716-446655440000</ID>", true, "<Owner><ID>0</ID>"},
		{"pure numeric", "<Owner><ID>12345</ID>", false, "<Owner><ID>12345</ID>"},
		{"empty id", "<Owner><ID></ID>", false, "<Owner><ID></ID>"},
		{"with whitespace", "<Owner>\n<ID>anon</ID>", true, "<Owner>\n<ID>0</ID>"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ownerIDRe.ReplaceAll([]byte(tt.input), []byte("${1}0${3}"))
			if tt.matches {
				assert.NotEqual(t, tt.input, string(result), "expected regex to match")
			} else {
				assert.Equal(t, tt.input, string(result), "expected regex NOT to match")
			}
			assert.Equal(t, tt.output, string(result))
		})
	}
}
