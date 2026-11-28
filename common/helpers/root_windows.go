//go:build windows

package helpers

// RunningAsRoot reports whether the current process runs as root. Windows does
// not have the concept of a root user, so this always returns false.
func RunningAsRoot() bool {
	return false
}
