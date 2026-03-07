package helpers

import (
	"encoding/hex"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"time"
)

func OpenURL(url string) error {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "windows":
		cmd = "cmd"
		args = []string{"/c", "start", url}
	case "darwin":
		cmd = "open"
		args = []string{url}
	default: // "linux", "freebsd", "openbsd", "netbsd"
		cmd = "xdg-open"
		args = []string{url}
	}

	// For Windows, to prevent issues with URLs starting with a quote,
	// an empty string is added as the first argument to 'start'.
	if runtime.GOOS == "windows" && len(args) > 1 {
		args = append(args[:1], append([]string{""}, args[1:]...)...)
	}

	return exec.Command(cmd, args...).Start()
}

func RandStringBytesMaskImprSrc(n int) string {

	var src = rand.New(rand.NewSource(time.Now().UnixNano()))

	b := make([]byte, (n+1)/2) // can be simplified to n/2 if n is always even

	if _, err := src.Read(b); err != nil {
		panic(err)
	}

	return hex.EncodeToString(b)[:n]
}

func ResolveOsEnvPath(inputPath string) string {
	// 1. Convert Windows %Variable% to $Variable
	// This regex finds content between two % and adds a $ prefix
	re := regexp.MustCompile(`%([^%]+)%`)
	intermediate := re.ReplaceAllString(inputPath, "$$1")

	// 2. Expand all $Variables using the system environment
	expanded := os.ExpandEnv(intermediate)

	// 3. Fix slashes: converts / to \ on Windows, and vice versa on Linux
	return filepath.FromSlash(filepath.Clean(expanded))
}
