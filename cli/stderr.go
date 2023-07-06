package cli

import (
	"fmt"
	"os"
)

func PrintStderr(msg string) {
	fmt.Fprintf(os.Stderr, msg + "\n")
}

func FatalStderr(msg string, code int) {
	PrintStderr("Fatal Error: " + msg)
	os.Exit(code)
}
