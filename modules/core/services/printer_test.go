//go:build !windows

package services

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/nutrixpos/pos/modules/core/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOpenPrinter_Validation(t *testing.T) {
	_, err := openPrinter(models.PrinterSettings{})
	assert.Error(t, err)

	_, err = openPrinter(models.PrinterSettings{ConnectionType: models.PrinterConnectionTypeUSB})
	assert.Error(t, err)

	_, err = openPrinter(models.PrinterSettings{Host: "127.0.0.1:1"})
	assert.Error(t, err)
}

func TestOpenPrinter_USBDevicePath(t *testing.T) {
	path := filepath.Join(t.TempDir(), "lp0")
	require.NoError(t, os.WriteFile(path, nil, 0o600))

	conn, err := openPrinter(models.PrinterSettings{
		ConnectionType: models.PrinterConnectionTypeUSB,
		Name:           path,
	})
	require.NoError(t, err)

	_, err = conn.Write([]byte("hello"))
	require.NoError(t, err)
	require.NoError(t, conn.Close())

	data, err := os.ReadFile(path)
	require.NoError(t, err)
	assert.Equal(t, "hello", string(data))
}

func TestOpenPrinter_USBCUPS(t *testing.T) {
	dir := t.TempDir()
	output := filepath.Join(dir, "out.bin")
	lp := filepath.Join(dir, "lp")

	script := "#!/bin/sh\ncat > \"" + output + "\"\n"
	require.NoError(t, os.WriteFile(lp, []byte(script), 0o755))
	t.Setenv("PATH", dir+string(os.PathListSeparator)+os.Getenv("PATH"))

	conn, err := openPrinter(models.PrinterSettings{
		ConnectionType: models.PrinterConnectionTypeUSB,
		Name:           "TestQueue",
	})
	require.NoError(t, err)

	_, err = conn.Write([]byte("receipt"))
	require.NoError(t, err)
	require.NoError(t, conn.Close())

	data, err := os.ReadFile(output)
	require.NoError(t, err)
	assert.Equal(t, "receipt", string(data))
}

func TestListAvailablePrinters(t *testing.T) {
	result, err := ListAvailablePrinters()
	require.NoError(t, err)
	assert.NotEmpty(t, result.OS)
	assert.NotNil(t, result.Printers)
	assert.NotNil(t, result.Browsers)
}
