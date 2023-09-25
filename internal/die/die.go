package die

import (
	"fmt"
	"os"
)

var ErrAction = func() {
	os.Exit(1)
}

func Die(x string) {
	_, _ = fmt.Fprintf(os.Stderr, x)
	os.Exit(1)
}
