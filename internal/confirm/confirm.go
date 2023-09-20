package confirm

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"strings"
)

// FAsk asks the user for confirmation on any Reader
// If any of the arguments passed as `overrides` is true, it returns true.
func FAsk(in io.Reader, s string, overrides ...bool) bool {
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

	if resp == "y" || resp == "yes" {
		return true
	}
	return false
}
