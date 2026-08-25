//go:build windows

package main

import (
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strings"

	"github.com/crc-org/admin-helper/pkg/logging"
	"golang.org/x/sys/windows"
)

const allowedProcessName = "crc.exe"

var allowedProcessPath = filepath.Join(os.Getenv("ProgramFiles"), "Red Hat OpenShift Local", allowedProcessName)

// fdProvider is an interface for extracting the underlying file descriptor
// from go-winio pipe connections (win32Pipe and win32MessageBytePipe both
// embed win32File which has Fd()).
type fdProvider interface {
	Fd() uintptr
}

// getClientExePath resolves the connecting client's PID to its executable path.
func getClientExePath(pipeHandle windows.Handle) (string, error) {
	var pid uint32
	if err := windows.GetNamedPipeClientProcessId(pipeHandle, &pid); err != nil {
		return "", fmt.Errorf("GetNamedPipeClientProcessId: %w", err)
	}

	proc, err := windows.OpenProcess(windows.PROCESS_QUERY_LIMITED_INFORMATION, false, pid)
	if err != nil {
		return "", fmt.Errorf("OpenProcess(%d): %w", pid, err)
	}
	defer func() { _ = windows.CloseHandle(proc) }()

	size := uint32(windows.MAX_PATH)
	buf := make([]uint16, size)
	if err := windows.QueryFullProcessImageName(proc, 0, &buf[0], &size); err != nil {
		return "", fmt.Errorf("QueryFullProcessImageName(%d): %w", pid, err)
	}

	return windows.UTF16ToString(buf[:size]), nil
}

type verifiedListener struct {
	net.Listener
	logger *logging.Logger
}

func newVerifiedListener(ln net.Listener) net.Listener {
	return &verifiedListener{
		ln, logging.GetLogger(),
	}
}

func (vl *verifiedListener) Accept() (net.Conn, error) {
	for {
		conn, err := vl.Listener.Accept()
		if err != nil {
			return nil, err
		}

		fd, ok := conn.(fdProvider)
		if !ok {
			vl.logger.Warn("peer check: connection does not expose file descriptor, rejecting")
			conn.Close()
			continue
		}

		exePath, err := getClientExePath(windows.Handle(fd.Fd()))
		if err != nil {
			vl.logger.Warn("peer check: failed to identify client process, rejecting",
				"error", err)
			conn.Close()
			continue
		}

		exeName := filepath.Base(exePath)
		if !strings.EqualFold(exeName, allowedProcessName) {
			vl.logger.Warn("peer check: rejecting connection from unauthorized process",
				"process", exePath,
				"expected", allowedProcessName)
			conn.Close()
			continue
		}

		if !strings.EqualFold(exePath, allowedProcessPath) {
			vl.logger.Warn("peer check: rejecting connection from process with unauthorized path",
				"process", exePath,
				"expected", allowedProcessPath)
			conn.Close()
			continue
		}

		return conn, nil
	}
}
