package utils

import (
	"log/slog"
	"os"
	"strconv"
)

// This file provides utilities for storing boolean values in a file

// Read file with "true" / "false" as content. In case of issues, return defaultValue
func ReadBoolFile(path string, defaultValue bool) bool {
	data, err := os.ReadFile(path)
	if err != nil {
		slog.Error("Error in readBoolFile", "path", path, "error", err)
		return defaultValue
	}

	// Not using strconv.ParseBool() as it allows other values like : "TRUE"
	switch string(data) {
	case TRUE_STRING:
		return true
	case FALSE_STRING:
		return false
	default:
		return defaultValue
	}
}

func WriteBoolFile(path string, value bool) error {
	return os.WriteFile(path, []byte(strconv.FormatBool(value)), 0644)

}
