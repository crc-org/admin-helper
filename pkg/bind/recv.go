// +build !windows

package bind

import (
	"fmt"
	"net"
	"syscall"
)

type Result struct {
	Error            string `json:"error,omitempty"`
	UnixDomainSocket string `json:"socket,omitempty"`
}

// Receive a file descriptor (*net.TCPListener here) from an unix domain socket
// https://github.com/moby/vpnkit/blob/master/go/pkg/vpnkit/forward/vmnet_darwin.go#L16
// https://github.com/ftrvxmtrx/fd/blob/master/fd.go
func Recv(via *net.UnixConn, localIP string) (*net.TCPListener, error) {
	viaf, err := via.File()
	if err != nil {
		return nil, err
	}
	socket := int(viaf.Fd())
	defer viaf.Close()

	buf := make([]byte, syscall.CmsgSpace(4))
	_, _, _, _, err = syscall.Recvmsg(socket, nil, buf, 0)
	if err != nil {
		return nil, err
	}

	var msgs []syscall.SocketControlMessage
	msgs, err = syscall.ParseSocketControlMessage(buf)
	if err != nil {
		return nil, err
	}
	if len(msgs) != 1 {
		return nil, fmt.Errorf("unexpected number of messages (got %d)", len(msgs))
	}
	fds, err := syscall.ParseUnixRights(&msgs[0])
	if err != nil {
		return nil, err
	}
	if len(fds) != 1 {
		return nil, fmt.Errorf("unexpected number of fd (got %d)", len(fds))
	}
	return fdToListener(localIP, uintptr(fds[0]))
}

func fdToListener(ip string, newFD uintptr) (*net.TCPListener, error) {
	ln, err := net.ListenTCP("tcp", &net.TCPAddr{
		IP:   net.ParseIP(ip),
		Port: 0,
	})
	if err != nil {
		return nil, err
	}
	raw, err := ln.SyscallConn()
	if err != nil {
		_ = ln.Close()
		return nil, err
	}
	if err := switchFDs(raw, newFD); err != nil {
		_ = ln.Close()
		_ = syscall.Close(int(newFD))
		return nil, err
	}
	_ = syscall.Close(int(newFD))
	return ln, nil
}

func switchFDs(raw syscall.RawConn, newFD uintptr) error {
	var dupErr error
	err := raw.Control(func(fd uintptr) {
		dupErr = syscall.Dup2(int(newFD), int(fd))
	})
	if dupErr != nil {
		return dupErr
	}
	if err != nil {
		return err
	}
	return nil
}
