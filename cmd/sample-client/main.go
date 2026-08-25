//go:build windows
// +build windows

package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/Microsoft/go-winio"
)

func main() {
	client := &http.Client{
		Transport: &http.Transport{
			DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
				return winio.DialPipeContext(ctx, `\\.\pipe\crc-admin-helper`)
			},
		},
		Timeout: 5 * time.Second,
	}

	res, err := client.Get("http://unix/version")
	if err != nil {
		log.Fatal(err)
	}
	bin, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	fmt.Println(string(bin))
}
