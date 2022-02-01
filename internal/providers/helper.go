package providers

import (
	"fmt"
	"os"
)

func printEnv(s string, expandEnv bool) {
	if expandEnv {
		es := os.ExpandEnv(s)
		fmt.Print(es)
		return
	}
	fmt.Print(s)
}
