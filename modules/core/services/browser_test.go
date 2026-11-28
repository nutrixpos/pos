package services

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/nutrixpos/pos/common/config"
	"github.com/nutrixpos/pos/modules/core/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestResolveBrowserPath_ConfiguredAbsolutePath(t *testing.T) {
	path := filepath.Join(t.TempDir(), "browser")
	require.NoError(t, os.WriteFile(path, nil, 0o755))

	resolved, err := resolveBrowserPath(path)
	require.NoError(t, err)
	assert.Equal(t, path, resolved)
}

func TestResolveBrowserPath_ConfiguredName(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("PATH lookup semantics differ on windows")
	}

	dir := t.TempDir()
	name := filepath.Join(dir, "mybrowser")
	require.NoError(t, os.WriteFile(name, nil, 0o755))
	t.Setenv("PATH", dir)

	resolved, err := resolveBrowserPath("mybrowser")
	require.NoError(t, err)
	assert.Equal(t, name, resolved)
}

func TestResolveBrowserPath_MissingConfiguredPath(t *testing.T) {
	path := filepath.Join(t.TempDir(), "does-not-exist")

	_, err := resolveBrowserPath(path)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "is not executable")
}

func TestResolveBrowserPath_AutoDetect(t *testing.T) {
	if runtime.GOOS != "linux" {
		t.Skip("auto-detection test targets the linux candidate list")
	}

	dir := t.TempDir()
	name := filepath.Join(dir, "chromium")
	require.NoError(t, os.WriteFile(name, nil, 0o755))
	t.Setenv("PATH", dir)

	resolved, err := resolveBrowserPath("")
	require.NoError(t, err)
	assert.Equal(t, name, resolved)
}

func TestResolveBrowserPath_NoBrowser(t *testing.T) {
	if runtime.GOOS != "linux" {
		t.Skip("auto-detection test targets the linux candidate list")
	}

	t.Setenv("PATH", t.TempDir())

	// Absolute candidates may still exist on the machine; only assert the
	// error is actionable when none are found.
	_, err := resolveBrowserPath("")
	if err != nil {
		assert.Contains(t, err.Error(), "Settings -> Printer")
	}
}

func TestDetectBrowsers(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("PATH lookup semantics differ on windows")
	}

	dir := t.TempDir()
	chromium := filepath.Join(dir, "chromium")
	brave := filepath.Join(dir, "brave-browser")
	require.NoError(t, os.WriteFile(chromium, nil, 0o755))
	require.NoError(t, os.WriteFile(brave, nil, 0o755))
	t.Setenv("PATH", dir)

	browsers := detectBrowsers()
	assert.Contains(t, browsers, chromium)
	assert.Contains(t, browsers, brave)
}

func TestBrowserCandidates(t *testing.T) {
	assert.NotEmpty(t, browserCandidates())
}

func TestNewBrowserContext(t *testing.T) {
	path := filepath.Join(t.TempDir(), "browser")
	require.NoError(t, os.WriteFile(path, nil, 0o755))

	ctx, cancel, err := newBrowserContext(
		config.Config{},
		models.PrintingSettings{BrowserPath: path},
	)
	require.NoError(t, err)
	require.NotNil(t, ctx)
	cancel()

	_, _, err = newBrowserContext(
		config.Config{},
		models.PrintingSettings{BrowserPath: filepath.Join(t.TempDir(), "missing")},
	)
	assert.Error(t, err)
}
