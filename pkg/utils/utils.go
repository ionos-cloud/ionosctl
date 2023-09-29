package utils

import (
	"fmt"
	"io/ioutil"
	"regexp"

	"github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6/resources"
	"golang.org/x/crypto/ssh"
)

// ReadPublicKey from a specific path
func ReadPublicKey(path string) (key string, err error) {
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}
	pubKey, _, _, _, err := ssh.ParseAuthorizedKey(bytes)
	if err != nil {
		return "", err
	}
	return string(ssh.MarshalAuthorizedKey(pubKey)[:]), nil
}

// StringSlicesEqual returns true if 2 slices of
// type string are equal.
func StringSlicesEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}

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
