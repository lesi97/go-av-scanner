package utils

import (
	"os"
)

/*
A startup function that should be called on the main thread

that provides a developer with useful information
*/
func Startup(port string) {
	env := os.Getenv("GO_ENV")
	const protocol = "http://"
	if env == "" {
		env = "development"
	}
	name, version := readModuleInfo()
	ip := getLocalIp()

	PrintColour("brightWhite", "\n  > %s %s\n", name, version) 		// Package name
	PrintColour("brightBlack", "\tEnvironment: %s", env)			// Current environment
	PrintColour("brightMagenta", "\n\t- Local:")					// Localhost address label
	PrintColour("cyan", "\t  %v%v%v", protocol, "localhost", port)	// Localhost address value
	if ip != "" {
		PrintColour("brightMagenta", "\n\t- Network:")				// Network address label
		PrintColour("cyan", "\t  %v%v%v\n", protocol,ip, port)		// Network address value
	}

	PrintColour("green", "\n  ✓ Server Ready\n\n")					// Server ready
}