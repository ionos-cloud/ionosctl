package uuidgen

import (
	"fmt"

	"github.com/gofrs/uuid/v5"
	"github.com/ionos-cloud/ionosctl/v6/internal/die"
)

// Must returns an UUIDv5 namespaced to ionostl repo, or fatally dies.
// If given names as parameters, it will iterate through each of the names, using the previously generated IDv5 as a namespace
func Must(names ...string) string {
	if len(names) == 0 {
		v4, err := uuid.NewV4()
		if err != nil {
			die.Die(fmt.Errorf("failed generating a random name UUID: %w", err).Error())
		}
		names[0] = v4.String()
	}

	ns := uuid.NewV5(uuid.NamespaceURL, "github.com/ionos-cloud/ionosctl")

	for _, name := range names {
		ns = uuid.NewV5(ns, name)
	}
	return ns.String()
}
