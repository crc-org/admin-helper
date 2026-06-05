package sandbox

// Profile is the macOS sandbox profile applied by admin-helper.
// It restricts the process to read/write /etc/hosts only, denies network
// and process execution. System paths come from the bsd.sb import.
const Profile = `(version 1)
(deny default)
(import "bsd.sb")

; Allow read/write to hosts file only
(allow file-read* file-write*
  (literal "/etc/hosts")
  (literal "/private/etc/hosts"))

; Deny network completely
(deny network*)

; Deny process spawning
(deny process-exec process-fork)
`
