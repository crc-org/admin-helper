//go:build ignore
// +build ignore

// This is a test program to verify the sandbox is actually working.
// It should be able to access /etc/hosts but fail to access other files.
package main

import (
	"fmt"
	"os"
	"runtime"
	"unsafe"

	"github.com/crc-org/admin-helper/pkg/sandbox"
)

/*
#cgo LDFLAGS: -framework Security
#include <sandbox.h>
#include <stdlib.h>
*/
import "C"

func applySandbox() error {
	if runtime.GOOS != "darwin" {
		return nil
	}

	cProfile := C.CString(sandbox.Profile(os.Getenv("HOME")))
	defer C.free(unsafe.Pointer(cProfile))

	var errBuf *C.char
	ret := C.sandbox_init(cProfile, 0, &errBuf)

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

func testFileAccess(path string, shouldSucceed bool) bool {
	_, err := os.ReadFile(path)
	if shouldSucceed {
		if err != nil {
			fmt.Printf("✗ FAIL: Should be able to read %s but got error: %v\n", path, err)
			return false
		}
		fmt.Printf("✓ PASS: Successfully read %s\n", path)
		return true
	}
	if err != nil {
		fmt.Printf("✓ PASS: Correctly blocked access to %s (error: %v)\n", path, err)
		return true
	}
	fmt.Printf("✗ FAIL: Should NOT be able to read %s but succeeded!\n", path)
	return false
}

func main() {
	fmt.Println("=== Testing Sandbox Restrictions ===")
	fmt.Println()

	if err := applySandbox(); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to apply sandbox: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Sandbox applied successfully!")
	fmt.Println()

	fmt.Println("Testing file access (should succeed):")
	failed := 0
	if !testFileAccess("/etc/hosts", true) {
		failed++
	}
	if !testFileAccess("/private/etc/hosts", true) {
		failed++
	}
	if !testFileAccess("/etc/ssh/ssh_config", false) {
		failed++
	}
	fmt.Println()

	fmt.Println("Testing file access (should be blocked by sandbox):")
	if !testFileAccess("/Users", false) {
		failed++
	}
	fmt.Println()

	if failed > 0 {
		fmt.Fprintf(os.Stderr, "%d sandbox test(s) failed\n", failed)
		os.Exit(1)
	}

	fmt.Println("All sandbox tests passed!")
	fmt.Println("To see sandbox denials from the CLI (use /usr/bin/log — zsh shadows 'log' with a builtin):")
	fmt.Println(`/usr/bin/log show --last 5m --style compact --predicate 'eventMessage CONTAINS "verify_sandbox" AND eventMessage CONTAINS "deny"'`)
}
