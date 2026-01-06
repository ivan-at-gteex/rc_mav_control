package tools

import (
	"errors"
	"path/filepath"
	"runtime"
)

func GetCurrentFilePath() (string, error) {
	_, filename, _, ok := runtime.Caller(1)
	if !ok {
		return "", errors.New("failed to get caller information")
	}

	absPath := filepath.Dir(filename)
	return absPath, nil
}
