//go:build windows

package services

import (
	"fmt"
	"sort"
	"sync"
	"unsafe"

	"golang.org/x/sys/windows"
)

var (
	winspool             = windows.NewLazySystemDLL("winspool.drv")
	procOpenPrinterW     = winspool.NewProc("OpenPrinterW")
	procClosePrinter     = winspool.NewProc("ClosePrinter")
	procStartDocPrinterW = winspool.NewProc("StartDocPrinterW")
	procEndDocPrinter    = winspool.NewProc("EndDocPrinter")
	procStartPagePrinter = winspool.NewProc("StartPagePrinter")
	procEndPagePrinter   = winspool.NewProc("EndPagePrinter")
	procWritePrinter     = winspool.NewProc("WritePrinter")
	procEnumPrintersW    = winspool.NewProc("EnumPrintersW")
)

const (
	printerEnumLocal       = 0x00000002
	printerEnumConnections = 0x00000004
)

type docInfo1 struct {
	DocName    *uint16
	OutputFile *uint16
	DataType   *uint16
}

type printerInfo4 struct {
	PrinterName *uint16
	ServerName  *uint16
	Attributes  uint32
}

type windowsPrinter struct {
	handle windows.Handle
	once   sync.Once
	err    error
}

func openUSBPrinter(name string) (printerTransport, error) {
	namePtr, err := windows.UTF16PtrFromString(name)
	if err != nil {
		return nil, err
	}

	var handle windows.Handle
	if r, _, callErr := procOpenPrinterW.Call(uintptr(unsafe.Pointer(namePtr)), uintptr(unsafe.Pointer(&handle)), 0); r == 0 {
		return nil, fmt.Errorf("failed to open printer %q: %w", name, callErr)
	}

	docName, _ := windows.UTF16PtrFromString("NutrixPOS Receipt")
	dataType, _ := windows.UTF16PtrFromString("RAW")
	docInfo := docInfo1{DocName: docName, DataType: dataType}

	if r, _, callErr := procStartDocPrinterW.Call(uintptr(handle), 1, uintptr(unsafe.Pointer(&docInfo))); r == 0 {
		procClosePrinter.Call(uintptr(handle))
		return nil, fmt.Errorf("failed to start print job for %q: %w", name, callErr)
	}

	if r, _, callErr := procStartPagePrinter.Call(uintptr(handle)); r == 0 {
		procEndDocPrinter.Call(uintptr(handle))
		procClosePrinter.Call(uintptr(handle))
		return nil, fmt.Errorf("failed to start print page for %q: %w", name, callErr)
	}

	return &windowsPrinter{handle: handle}, nil
}

func (p *windowsPrinter) Write(data []byte) (int, error) {
	if len(data) == 0 {
		return 0, nil
	}

	var written uint32
	r, _, callErr := procWritePrinter.Call(
		uintptr(p.handle),
		uintptr(unsafe.Pointer(&data[0])),
		uintptr(len(data)),
		uintptr(unsafe.Pointer(&written)),
	)
	if r == 0 {
		return int(written), fmt.Errorf("failed to write to printer: %w", callErr)
	}

	return int(written), nil
}

func (p *windowsPrinter) Close() error {
	p.once.Do(func() {
		var firstErr error
		if r, _, callErr := procEndPagePrinter.Call(uintptr(p.handle)); r == 0 {
			firstErr = callErr
		}
		if r, _, callErr := procEndDocPrinter.Call(uintptr(p.handle)); r == 0 && firstErr == nil {
			firstErr = callErr
		}
		if r, _, callErr := procClosePrinter.Call(uintptr(p.handle)); r == 0 && firstErr == nil {
			firstErr = callErr
		}
		if firstErr != nil {
			p.err = fmt.Errorf("failed to close printer: %w", firstErr)
		}
	})

	return p.err
}

func listUSBPrinters() ([]string, error) {
	flags := uintptr(printerEnumLocal | printerEnumConnections)

	var needed, returned uint32
	_, _, _ = procEnumPrintersW.Call(
		flags,
		0,
		4,
		0,
		0,
		uintptr(unsafe.Pointer(&needed)),
		uintptr(unsafe.Pointer(&returned)),
	)
	if needed == 0 {
		return []string{}, nil
	}

	buf := make([]byte, needed)
	r, _, callErr := procEnumPrintersW.Call(
		flags,
		0,
		4,
		uintptr(unsafe.Pointer(&buf[0])),
		uintptr(needed),
		uintptr(unsafe.Pointer(&needed)),
		uintptr(unsafe.Pointer(&returned)),
	)
	if r == 0 {
		return nil, fmt.Errorf("failed to enumerate printers: %w", callErr)
	}

	infos := unsafe.Slice((*printerInfo4)(unsafe.Pointer(&buf[0])), int(returned))
	printers := make([]string, 0, len(infos))
	for _, info := range infos {
		if info.PrinterName != nil {
			printers = append(printers, windows.UTF16PtrToString(info.PrinterName))
		}
	}
	sort.Strings(printers)

	return printers, nil
}
