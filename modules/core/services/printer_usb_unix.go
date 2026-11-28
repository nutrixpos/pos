//go:build !windows

package services

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"sync"
)

// cupsPrinter sends raw ESC/POS bytes to a CUPS queue. Data is streamed to the
// lp process and the job is waited on when the transport is closed.
type cupsPrinter struct {
	cmd    *exec.Cmd
	stdin  io.WriteCloser
	stderr bytes.Buffer
	once   sync.Once
	err    error
}

func openUSBPrinter(name string) (printerTransport, error) {
	// An absolute path is treated as a device node such as /dev/usb/lp0.
	if filepath.IsAbs(name) {
		return os.OpenFile(name, os.O_RDWR, 0)
	}

	lp, err := exec.LookPath("lp")
	if err != nil {
		return nil, fmt.Errorf("cups lp command not found for usb printer %q: %w", name, err)
	}

	cmd := exec.Command(lp, "-d", name, "-o", "raw")
	cmd.Stdout = io.Discard

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, err
	}

	printer := &cupsPrinter{cmd: cmd, stdin: stdin}
	cmd.Stderr = &printer.stderr

	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("failed to start cups print job for %q: %w", name, err)
	}

	return printer, nil
}

func (c *cupsPrinter) Write(p []byte) (int, error) {
	return c.stdin.Write(p)
}

func (c *cupsPrinter) Close() error {
	c.once.Do(func() {
		if err := c.stdin.Close(); err != nil {
			_ = c.cmd.Wait()
			c.err = err
			return
		}

		if err := c.cmd.Wait(); err != nil {
			if msg := strings.TrimSpace(c.stderr.String()); msg != "" {
				c.err = fmt.Errorf("cups print job failed: %w: %s", err, msg)
				return
			}
			c.err = fmt.Errorf("cups print job failed: %w", err)
		}
	})

	return c.err
}

func listUSBPrinters() ([]string, error) {
	names := make(map[string]struct{})

	if out, err := exec.Command("lpstat", "-a").Output(); err == nil {
		for _, line := range strings.Split(string(out), "\n") {
			fields := strings.Fields(line)
			if len(fields) > 0 {
				names[fields[0]] = struct{}{}
			}
		}
	}

	if devices, err := filepath.Glob("/dev/usb/lp*"); err == nil {
		for _, device := range devices {
			names[device] = struct{}{}
		}
	}

	printers := make([]string, 0, len(names))
	for name := range names {
		printers = append(printers, name)
	}
	sort.Strings(printers)

	return printers, nil
}
