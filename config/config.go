// Package config resolves the runtime configuration directory for Godo.
// It follows the XDG Base Directory Specification:
//   - If $XDG_CONFIG_HOME is set and non-empty, use $XDG_CONFIG_HOME/godo
//   - Otherwise fall back to $HOME/.config/godo
//   - If $HOME is also unset, fall back to the current working directory
//
// No other package should hardcode a path — always go through Dir().
package config

import (
	"os"
	"path/filepath"
)

const appName = "godo"

// Dir returns the absolute path to Godo's configuration directory.
// The directory is not created by this function — that is the repository's
// responsibility via os.MkdirAll in NewJSONRepository.
func Dir() string {
	// XDG_CONFIG_HOME takes priority.
	if xdg := os.Getenv("XDG_CONFIG_HOME"); xdg != "" {
		return filepath.Join(xdg, appName)
	}

	// Fall back to $HOME/.config/godo.
	if home := os.Getenv("HOME"); home != "" {
		return filepath.Join(home, ".config", appName)
	}

	// Last resort: store config alongside the binary in the working directory.
	// This should only happen in unusual environments (e.g. containers with no $HOME).
	wd, err := os.Getwd()
	if err != nil {
		return appName // relative path of last resort
	}
	return filepath.Join(wd, appName)
}
