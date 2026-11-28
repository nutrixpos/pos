package services

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/chromedp/chromedp"
	"github.com/nutrixpos/pos/common/config"
	"github.com/nutrixpos/pos/common/helpers"
	"github.com/nutrixpos/pos/modules/core/models"
)

// resolveBrowserPath finds the Chromium-based executable used to render
// receipts. It uses the configured path when set and otherwise auto-detects a
// browser on the host.
func resolveBrowserPath(browserPath string) (string, error) {
	configured := strings.TrimSpace(browserPath)
	if configured != "" {
		path, err := lookupBrowser(configured)
		if err != nil {
			return "", fmt.Errorf("receipt browser %q is not executable: %w", configured, err)
		}
		return path, nil
	}

	for _, candidate := range browserCandidates() {
		if path, err := exec.LookPath(candidate); err == nil {
			return path, nil
		}
	}

	return "", fmt.Errorf("receipt rendering requires a Chromium-based browser; none found. Install Chromium/Chrome/Edge (or another Chromium-based browser), or select one in Settings -> Printer")
}

// detectBrowsers returns every Chromium-based executable found on the host.
// It is used to populate the browser dropdown in the settings screen.
func detectBrowsers() []string {
	seen := make(map[string]struct{})
	browsers := make([]string, 0)

	for _, candidate := range browserCandidates() {
		path, err := exec.LookPath(candidate)
		if err != nil {
			continue
		}
		if _, ok := seen[path]; ok {
			continue
		}
		seen[path] = struct{}{}
		browsers = append(browsers, path)
	}

	return browsers
}

func lookupBrowser(configured string) (string, error) {
	if filepath.IsAbs(configured) || strings.ContainsAny(configured, `/\`) {
		info, err := os.Stat(configured)
		if err != nil {
			return "", err
		}
		if info.IsDir() {
			return "", fmt.Errorf("is a directory")
		}
		return configured, nil
	}

	return exec.LookPath(configured)
}

func browserCandidates() []string {
	switch runtime.GOOS {
	case "windows":
		return []string{
			"msedge",
			"msedge.exe",
			"chrome",
			"chrome.exe",
			"brave",
			"brave.exe",
			"vivaldi",
			"vivaldi.exe",
			"opera",
			"opera.exe",
			"chromium",
			"chromium.exe",
			"headless_shell",
			"headless-shell",
			`C:\Program Files (x86)\Microsoft\Edge\Application\msedge.exe`,
			`C:\Program Files\Microsoft\Edge\Application\msedge.exe`,
			`C:\Program Files (x86)\Google\Chrome\Application\chrome.exe`,
			`C:\Program Files\Google\Chrome\Application\chrome.exe`,
			filepath.Join(os.Getenv("USERPROFILE"), `AppData\Local\Google\Chrome\Application\chrome.exe`),
			filepath.Join(os.Getenv("USERPROFILE"), `AppData\Local\Chromium\Application\chrome.exe`),
			`C:\Program Files\BraveSoftware\Brave-Browser\Application\brave.exe`,
			`C:\Program Files (x86)\BraveSoftware\Brave-Browser\Application\brave.exe`,
			filepath.Join(os.Getenv("USERPROFILE"), `AppData\Local\BraveSoftware\Brave-Browser\Application\brave.exe`),
			filepath.Join(os.Getenv("USERPROFILE"), `AppData\Local\Vivaldi\Application\vivaldi.exe`),
			`C:\Program Files\Vivaldi\Application\vivaldi.exe`,
			filepath.Join(os.Getenv("USERPROFILE"), `AppData\Local\Programs\Opera\opera.exe`),
			`C:\Program Files\Opera\opera.exe`,
		}
	case "darwin":
		return []string{
			"/Applications/Google Chrome.app/Contents/MacOS/Google Chrome",
			"/Applications/Chromium.app/Contents/MacOS/Chromium",
			"/Applications/Microsoft Edge.app/Contents/MacOS/Microsoft Edge",
			"/Applications/Brave Browser.app/Contents/MacOS/Brave Browser",
			"/Applications/Vivaldi.app/Contents/MacOS/Vivaldi",
			"/Applications/Opera.app/Contents/MacOS/Opera",
			"chromium",
			"google-chrome",
			"microsoft-edge",
			"brave-browser",
			"vivaldi",
			"opera",
			"headless_shell",
			"headless-shell",
		}
	default:
		return []string{
			"chromium",
			"chromium-browser",
			"google-chrome",
			"google-chrome-stable",
			"google-chrome-beta",
			"google-chrome-unstable",
			"microsoft-edge",
			"microsoft-edge-stable",
			"microsoft-edge-beta",
			"brave-browser",
			"brave-browser-stable",
			"brave",
			"vivaldi",
			"vivaldi-stable",
			"opera",
			"opera-stable",
			"headless_shell",
			"headless-shell",
			"chrome",
			"/usr/bin/google-chrome",
			"/usr/local/bin/chrome",
			"/snap/bin/chromium",
			"/snap/bin/brave",
			"/snap/bin/vivaldi",
			"/snap/bin/opera",
		}
	}
}

// newBrowserContext creates a chromedp context using the resolved browser
// executable. The returned cancel function releases both the browser context
// and its allocator.
func newBrowserContext(cfg config.Config, printing models.PrintingSettings) (context.Context, context.CancelFunc, error) {
	path, err := resolveBrowserPath(printing.BrowserPath)
	if err != nil {
		return nil, nil, err
	}

	opts := make([]chromedp.ExecAllocatorOption, len(chromedp.DefaultExecAllocatorOptions))
	copy(opts, chromedp.DefaultExecAllocatorOptions[:])
	opts = append(opts, chromedp.ExecPath(path))

	if cfg.Printing.NoSandbox || helpers.RunningAsRoot() {
		opts = append(opts, chromedp.NoSandbox)
	}

	allocCtx, allocCancel := chromedp.NewExecAllocator(context.Background(), opts...)
	ctx, cancel := chromedp.NewContext(allocCtx)

	return ctx, func() {
		cancel()
		allocCancel()
	}, nil
}
