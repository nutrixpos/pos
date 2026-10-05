package helpers

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRandStringBytesMaskImprSrc_Length(t *testing.T) {
	for _, n := range []int{0, 1, 2, 5, 10, 31, 100} {
		s := RandStringBytesMaskImprSrc(n)
		assert.Len(t, s, n, "length mismatch for n=%d", n)
	}
}

func TestRandStringBytesMaskImprSrc_Unique(t *testing.T) {
	seen := make(map[string]bool)
	for i := 0; i < 1000; i++ {
		s := RandStringBytesMaskImprSrc(16)
		require.False(t, seen[s], "duplicate random string %q", s)
		seen[s] = true
	}
}

func TestResolveOsEnvPath_EnvExpansion(t *testing.T) {
	t.Setenv("NUTRIX_TEST_DIR", "/tmp/nutrix/test")

	assert.Equal(t, "/tmp/nutrix/test/sub", ResolveOsEnvPath("$NUTRIX_TEST_DIR/sub"))
	assert.Equal(t, "/tmp/nutrix/test/sub", ResolveOsEnvPath("${NUTRIX_TEST_DIR}/sub"))
}

func TestResolveOsEnvPath_CleansSlashes(t *testing.T) {
	assert.Equal(t, "/a/b/c", ResolveOsEnvPath("/a//b/./c"))
}

func TestResolveOsEnvPath_UnsetVarExpandsToEmpty(t *testing.T) {
	// Unset variables expand to the empty string; cleaning an empty path
	// yields ".".
	assert.Equal(t, ".", ResolveOsEnvPath("$UNSET_NUTRIX_VAR"))
}
