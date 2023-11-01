package request

import (
	"fmt"
	"regexp"

	"github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6/resources"
)

// This package operates on SDK's Response type

var requestPathRegex = regexp.MustCompile(`https?://[a-zA-Z0-9./-]+/requests/([a-z0-9-]+)/status`)

func GetRequestId(path string) (string, error) {
	if !requestPathRegex.MatchString(path) {
		return "", fmt.Errorf("%s does not contain requestId", path)
	}
	return requestPathRegex.FindStringSubmatch(path)[1], nil
}

func GetId(r *resources.Response) string {
	if id, err := GetRequestId(GetRequestPath(r)); err == nil {
		return id
	}
	return ""
}

func GetRequestPath(r *resources.Response) string {
	if r != nil && r.Header != nil && len(r.Header) > 0 {
		return r.Header.Get("location")
	}
	return ""
}
