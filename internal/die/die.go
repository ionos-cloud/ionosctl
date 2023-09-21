package die

import (
	"fmt"
	"os"
)

func Die(x string) {
	fmt.Fprintf(x)
	os.Exit(1)
}
