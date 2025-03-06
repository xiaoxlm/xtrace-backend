package util

import (
	"fmt"
	"os"
	"path/filepath"
)

func GetProjectPath() (string, error) {
	currentDir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	// Walk up the directory tree to find project root (where go.mod is located)
	projectRoot := currentDir
	for {
		if _, err := os.Stat(filepath.Join(projectRoot, "go.mod")); err == nil {
			break
		}
		parent := filepath.Dir(projectRoot)
		if parent == projectRoot {
			return "", fmt.Errorf("Could not find project root - no go.mod file found")
		}
		projectRoot = parent
	}

	return projectRoot, nil
}
