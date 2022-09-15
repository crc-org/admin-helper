//go:build !windows
// +build !windows

package main

import "net"

func listen() (net.Listener, error) {
	return net.Listen("unix", "daemon.sock")
}
