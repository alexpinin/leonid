package testutil

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
)

type pkgPathFinder struct{}

// FindEnvFile returns the env file of the service
//
// The expected directory structure is:
//
//	root
//	  service
//	    env
func FindEnvFile(service, fileName string) (string, error) {
	_, filename, _, _ := runtime.Caller(0)
	currDir := filepath.Dir(filename)
	currPkg := reflect.TypeOf(pkgPathFinder{}).PkgPath()

	index := strings.LastIndex(currDir, currPkg)
	if index == -1 {
		return "", fmt.Errorf("invalid directory structure: package name not found") // unlikely to happen
	}

	envFile := filepath.Join(currDir[:index], service, "env", fileName)
	if _, err := os.Stat(envFile); os.IsNotExist(err) {
		return "", fmt.Errorf("env file %s not found", filepath.Join(service, "env", fileName))
	}

	return envFile, nil
}
