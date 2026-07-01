//go:build !darwin && !linux
// +build !darwin,!linux

package hosts

import "os"

// OpenHostsFileAndDropPrivileges is not supported on this platform.
func OpenHostsFileAndDropPrivileges(path string) (*os.File, error) {
	return nil, ErrPrivilegeDropUnsupported
}
