package utils

import (
	"os"
)

/*
A startup function that should be called on the main thread

that provides a developer with useful information
*/
func Startup(l *Logger,port string) {
	env := os.Getenv("GO_ENV")
	const protocol = "http://"
	if env == "" {
		env = "development"
	}
	name, version := readModuleInfo()
	ip := getLocalIp()

	l.PrintColour(false, "brightWhite", "\n  > %s %s\n", name, version) 	// Package name
	l.PrintColour(false, "brightBlack", "\tEnvironment: %s", env)			// Current environment
	l.PrintColour(false, "brightMagenta", "\n\t- Local:")					// Localhost address label
	l.PrintColour(false, "cyan", "\t  %v%v%v", protocol, "localhost", port)	// Localhost address value
	if ip != "" {
		l.PrintColour(false, "brightMagenta", "\n\t- Network:")				// Network address label
		l.PrintColour(false, "cyan", "\t  %v%v%v\n", protocol,ip, port)		// Network address value
	}

	l.PrintColour(false, "green", "\n  ✓ Server Ready\n\n")					// Server ready
}