package die

import (
	"fmt"
	"os"
)

func Die(x string) {
	_, _ = fmt.Fprintf(os.Stderr, x)
	os.Exit(1)
}
