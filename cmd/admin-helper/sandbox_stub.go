//go:build !darwin || !cgo
// +build !darwin !cgo

package main

// applySandbox is a no-op on non-Darwin platforms or when CGO is disabled.
// Sandboxing is only available on macOS with CGO enabled.
func applySandbox() error {
	return nil
}
