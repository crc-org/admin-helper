// +build !windows

package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"strconv"

	"github.com/code-ready/admin-helper/pkg/bind"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	ip := "127.0.0.1"
	port := 8080

	// #nosec G204
	cmd := exec.Command("./admin-helper", "bind", ip, strconv.Itoa(port))
	cmd.Stderr = os.Stderr
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	if err := cmd.Start(); err != nil {
		return err
	}

	// admin-helper sends either an error or the unix domain socket to connect to
	bin, err := ioutil.ReadAll(stdout)
	if err != nil {
		return err
	}
	var result bind.Result
	if err := json.Unmarshal(bin, &result); err != nil {
		return err
	}

	if result.Error != "" {
		return errors.New(result.Error)
	}

	// connect the returned socket and grab the listener
	viaconn, err := net.Dial("unix", result.UnixDomainSocket)
	if err != nil {
		return err
	}

	ln, err := bind.Recv(viaconn.(*net.UnixConn), ip)
	if err != nil {
		return err
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		_, _ = writer.Write([]byte(`Hello world`))
	})
	return http.Serve(ln, mux)
}
