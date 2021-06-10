package utils

import "fmt"

func ColsMessage(cols []string) string {
	return fmt.Sprintf("Set of columns to be printed on output \nAvailable columns: %v", cols)
}
