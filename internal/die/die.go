package die

import (
	"fmt"
	"os"
)

func Die(x string) {
	fmt.Fprintf(os.Stderr, x)
	os.Exit(1)
}
