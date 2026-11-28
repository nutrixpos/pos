package services

import (
	"fmt"
	"io"
	"net"
	"runtime"

	"github.com/nutrixpos/pos/modules/core/models"
)

// printerTransport is a write-only connection to a receipt printer.
type printerTransport interface {
	io.Writer
	Close() error
}

// PrinterDiscovery lists the printers available on the host running the
// backend. It is used by the settings screen to populate the USB printer
// dropdown.
type PrinterDiscovery struct {
	OS       string   `json:"os"`
	Printers []string `json:"printers"`
	Browsers []string `json:"browsers"`
}

// openPrinter returns a transport for the configured printer. Empty
// connection types default to network for backward compatibility with
// settings saved before USB support existed.
func openPrinter(printer models.PrinterSettings) (printerTransport, error) {
	if printer.IsUSB() {
		if printer.Name == "" {
			return nil, fmt.Errorf("usb printer name is not configured")
		}
		return openUSBPrinter(printer.Name)
	}

	if printer.Host == "" {
		return nil, fmt.Errorf("printer host is not configured")
	}

	return net.Dial("tcp", fmt.Sprintf("%s:9100", printer.Host))
}

// ListAvailablePrinters returns the OS printer queues and USB devices that can
// be selected for USB printing.
func ListAvailablePrinters() (PrinterDiscovery, error) {
	printers, err := listUSBPrinters()
	if err != nil {
		return PrinterDiscovery{}, err
	}
	if printers == nil {
		printers = []string{}
	}

	return PrinterDiscovery{
		OS:       runtime.GOOS,
		Printers: printers,
		Browsers: detectBrowsers(),
	}, nil
}
