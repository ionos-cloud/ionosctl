package utils

import (
	"fmt"
	"io/ioutil"
	"net"
	"strings"
	"time"

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

// ParseDate parses a date string in RFC3339 format
// and returns a time.Time object
// It also removes the [UTC] suffix if present
func ParseDate(strDate string) (time.Time, error) {
	strDate = strings.TrimSuffix(strDate, "[UTC]")
	date, err := time.Parse(time.RFC3339, strDate)
	if err != nil {
		return time.Time{}, fmt.Errorf("date is not a valid RFC3339: %w", err)
	}
	return date, nil
}
