package helpers

import (
	"encoding/hex"
	"math/rand"
	"os"
	"path/filepath"
	"regexp"
	"time"
)

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
