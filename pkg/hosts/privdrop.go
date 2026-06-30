package hosts

import "errors"

// ErrPrivilegeDropUnsupported is returned when privilege dropping is not available.
var ErrPrivilegeDropUnsupported = errors.New("privilege drop not supported on this platform")
