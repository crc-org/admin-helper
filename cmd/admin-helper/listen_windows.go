package main

import (
	"fmt"
	"net"

	"github.com/Microsoft/go-winio"
)

func listen() (net.Listener, error) {
	// see https://github.com/moby/moby/blob/46cdcd206c56172b95ba5c77b827a722dab426c5/daemon/listeners/listeners_windows.go
	// allow Administrators and SYSTEM, plus whatever additional users or groups were specified
	sddl := "D:P(A;;GA;;;BA)(A;;GA;;;SY)"
	sid, err := winio.LookupSidByName("crc-users")
	if err != nil {
		return nil, err
	}
	sddl += fmt.Sprintf("(A;;GRGW;;;%s)", sid)

	return winio.ListenPipe(`\\.\pipe\crc-admin-helper`, &winio.PipeConfig{
		SecurityDescriptor: sddl,  // Administrators and system
		MessageMode:        true,  // Use message mode so that CloseWrite() is supported
		InputBufferSize:    65536, // Use 64KB buffers to improve performance
		OutputBufferSize:   65536,
	})
}
