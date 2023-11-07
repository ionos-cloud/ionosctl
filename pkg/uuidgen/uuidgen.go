package uuidgen

import (
	"fmt"

	"github.com/gofrs/uuid/v5"
	"github.com/ionos-cloud/ionosctl/v6/pkg/die"
)

const ns = "github.com/ionos-cloud/ionosctl"

// Must returns an UUIDv5 namespaced to ionosctl repo, or fatally dies.
// If given names as parameters, it will iterate through each of the names, using the previously generated IDv5 as a namespace
func Must(names ...string) string {
	// Convert the namespace string to a UUID
	namespaceUUID := uuid.NewV5(uuid.NamespaceDNS, ns)

	if len(names) == 0 {
		v4, err := uuid.NewV4()
		if err != nil {
			die.Die(fmt.Errorf("failed generating a random name UUID: %w", err).Error())
		}
		names = append(names, v4.String())
	}

	for _, name := range names {
		namespaceUUID = uuid.NewV5(namespaceUUID, name)
	}
	return namespaceUUID.String()
}

// MustSingle generates an UUIDv5 namespaced to github.com/ionos-cloud/ionosctl while guaranteeing a single name is used
// this func simply uses Must() under the hood, but guarantees no other names are used, and that no UUIDv4 is used as name
func MustSingle(name string) string {
	return Must(name)
}
