package confirm

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"strings"
)

const UserDenied = "user denied confirmation"

// Confirmer is an interface for different confirmation strategies.
type Confirmer interface {
	Ask(in io.Reader, s string, overrides ...bool) bool
}

// defaultConfirmer is the default behavior used by FAsk.
type defaultConfirmer struct{}

func (d defaultConfirmer) Ask(in io.Reader, s string, overrides ...bool) bool {
	for _, o := range overrides {
		if o {
			return true
		}
	}

	rr := bufio.NewReader(in)
	fmt.Printf("%s? [y/n]: ", s)

	resp, err := rr.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}

	resp = strings.ToLower(strings.TrimSpace(resp))

	return resp == "y" || resp == "yes"
}

// currentStrategy holds the current confirmation strategy.
// If nil, default behavior is used.
var currentStrategy Confirmer

// SetStrategy sets the current confirmation strategy.
func SetStrategy(strategy Confirmer) {
	currentStrategy = strategy
}

// FAsk asks the user for confirmation using the current strategy
// or the default behavior if no strategy is set.
func FAsk(in io.Reader, s string, overrides ...bool) bool {
	if currentStrategy != nil {
		return currentStrategy.Ask(in, s, overrides...)
	}
	return defaultConfirmer{}.Ask(in, s, overrides...)
}
