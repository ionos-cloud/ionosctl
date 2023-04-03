/*
 * Certificate Manager Service API
 *
 * Using the Certificate Manager Service, you can conveniently provision and manage SSL certificates with IONOS services and your internal connected resources. For the [Application Load Balancer](https://api.ionos.com/docs/cloud/v6/#Application-Load-Balancers-get-datacenters-datacenterId-applicationloadbalancers), you usually need a certificate to encrypt your HTTPS traffic.  The service provides the basic functions of uploading and deleting your certificates for this purpose.
 *
 * API version: 1.0
 */

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package cert

import (
	"github.com/ionos-cloud/sdk-go-bundle/shared"
	"strings"
	"time"
)

type IonosTime struct {
	time.Time
}

func (t *IonosTime) UnmarshalJSON(data []byte) error {
	str := string(data)
	if shared.Strlen(str) == 0 {
		t = nil
		return nil
	}
	if str[0] == '"' {
		str = str[1:]
	}
	if str[len(str)-1] == '"' {
		str = str[:len(str)-1]
	}
	if !strings.Contains(str, "Z") {
		/* forcefully adding timezone suffix to be able to parse the
		 * string using RFC3339 */
		str += "Z"
	}
	tt, err := time.Parse(time.RFC3339, str)
	if err != nil {
		return err
	}
	*t = IonosTime{tt}
	return nil
}
