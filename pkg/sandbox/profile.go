package sandbox

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/crc-org/admin-helper/pkg/constants"
)

// Profile returns the macOS sandbox profile for the given home directory.
// It restricts the process to read/write /etc/hosts, allows SSH config paths,
// denies network and process execution. System paths come from the bsd.sb import.
func Profile(home string) string {
	crcAdminHelperLogFile := filepath.Join(home, ".crc", os.Getenv(constants.LogFileEnvVar))
	return fmt.Sprintf(`(version 1)
(deny default)
(import "bsd.sb")

; Allow read/write to hosts file only
(allow file-read* file-write*
  (literal "/etc/hosts")
  (literal "/private/etc/hosts")
  (literal %q))

; Deny network completely
(deny network*)

; Deny process spawning
(deny process-exec process-fork)
`, crcAdminHelperLogFile)
}
