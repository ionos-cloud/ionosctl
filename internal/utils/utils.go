package utils

import (
	"fmt"
	"io/ioutil"
	"net"
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

func ValidateIPv6CidrBlockAgainstParentCidrBlock(cidr string, expectedMask int, parentCidr string) error {
	ip, ipNet, err := net.ParseCIDR(cidr)
	if err != nil {
		return err
	}

	if ip.To4() != nil {
		return fmt.Errorf("this is not an IPv6 Cidr Block")
	}

	if ones, _ := ipNet.Mask.Size(); ones != expectedMask {
		return fmt.Errorf("network mask is not the expected size: %d should be %d", ones, expectedMask)
	}

	if !ip.Equal(ip.Mask(ipNet.Mask)) {
		return fmt.Errorf("network mask does not cover all IP bits set")
	}

	_, parentIPNet, err := net.ParseCIDR(parentCidr)
	if err != nil {
		return err
	}

	if !parentIPNet.Contains(ip) {
		return fmt.Errorf("child Cidr Block (%s) is not inside parent Cidr Block range (%s)", cidr, parentCidr)
	}

	return nil
}
