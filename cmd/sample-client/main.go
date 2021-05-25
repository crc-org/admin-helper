// +build windows

package main

import (
	"context"
	"fmt"
	"io/ioutil"
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
	bin, err := ioutil.ReadAll(res.Body)
	defer res.Body.Close()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(bin))
}
