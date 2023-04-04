package confirm

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

// Ask uses os.Stdin to ask the user for confirmation.
// If any of the arguments passed as `overrides` is true, it returns true.
func Ask(s string, overrides ...bool) bool {
	return FAsk(s, os.Stdin, overrides...)
}

// FAsk asks the user for confirmation on any Reader
// If any of the arguments passed as `overrides` is true, it returns true.
func FAsk(s string, in io.Reader, overrides ...bool) bool {
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
