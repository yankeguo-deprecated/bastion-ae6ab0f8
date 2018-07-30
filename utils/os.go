package utils

import "os"

var (
	exitCode int
)

// WillExit set a exit code
func WillExit(code int) {
	exitCode = code
}

// DoExit execute a previous set exit code with WillExit()
func DoExit() {
	os.Exit(exitCode)
}
