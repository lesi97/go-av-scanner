package utils

import (
	"os"
	"strings"
)

var (
	ModuleName = "clamav-service"
	ModuleVersion    = "v1.0.0"
)

/*
Function to read the version number and module name from the `go.mod` file
*/
func readModuleInfo() (name, version string) {
	content, err := os.ReadFile("./go.mod")
	if err != nil {
		return ModuleName, ModuleVersion
	}

	lines := strings.Split(string(content), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "module ") {
			name = strings.TrimSpace(strings.TrimPrefix(line, "module"))
		}
		if strings.HasPrefix(line, "// Version:") {
			version = strings.TrimSpace(strings.TrimPrefix(line, "// Version:"))
		}
	}

	if version == "" {
		version = ModuleVersion
	}
	return name, version
}