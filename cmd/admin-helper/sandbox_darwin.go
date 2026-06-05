//go:build darwin && cgo
// +build darwin,cgo

package main

/*
#cgo LDFLAGS: -framework Security
#include <sandbox.h>
#include <stdlib.h>
*/
import "C"
import (
	"fmt"
	"unsafe"

	"github.com/crc-org/admin-helper/pkg/sandbox"
)

func applySandbox() error {
	cProfile := C.CString(sandbox.Profile)
	defer C.free(unsafe.Pointer(cProfile))

	var errBuf *C.char
	// https://github.com/go-critic/go-critic/issues/897: dupSubExpr false positive from cgo
	ret := C.sandbox_init(cProfile, 0, &errBuf) //nolint:gocritic

	if ret != 0 {
		if errBuf != nil {
			err := C.GoString(errBuf)
			C.sandbox_free_error(errBuf)
			return fmt.Errorf("sandbox init failed: %s", err)
		}
		return fmt.Errorf("sandbox init failed with code %d", ret)
	}

	return nil
}
