// +build !windows

package bind

import (
	"net"
	"syscall"
)

// Send a file descriptor (*net.TCPListener here) to an unix domain socket
// https://github.com/ftrvxmtrx/fd/blob/master/fd.go#L85
func Send(via *net.UnixConn, ln *net.TCPListener) error {
	file, err := ln.File()
	if err != nil {
		return err
	}
	fd := int(file.Fd())
	viaf, err := via.File()
	if err != nil {
		return err
	}
	socket := int(viaf.Fd())
	defer viaf.Close()

	rights := syscall.UnixRights(fd)
	return syscall.Sendmsg(socket, nil, rights, nil, 0)
}
