package core

import (
	"fmt"
	"os"
)

// LogError writes an error message to stderr.
func LogError(msg string) {
	fmt.Fprintln(os.Stderr, "[ERROR]", msg)
}

// LogInfo writes an informational message to stdout.
func LogInfo(msg string) {
	fmt.Fprintln(os.Stdout, "[INFO]", msg)
}
