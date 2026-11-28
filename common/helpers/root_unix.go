//go:build !windows

package helpers

import "os"

// RunningAsRoot reports whether the current process runs as root.
func RunningAsRoot() bool {
	return os.Geteuid() == 0
}
