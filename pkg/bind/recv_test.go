// +build !windows

package bind

import (
	"io/ioutil"
	"net"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRoundTrip(t *testing.T) {
	dir, err := ioutil.TempDir("", "test-bind")
	assert.NoError(t, err)
	defer os.RemoveAll(dir)

	sentLn, err := net.Listen("tcp", "127.0.0.1:0")
	assert.NoError(t, err)
	defer sentLn.Close()

	uds := filepath.Join(dir, "server.sock")
	viaLn, err := net.Listen("unix", uds)
	assert.NoError(t, err)
	defer viaLn.Close()
	go func() {
		viaconn, err := viaLn.Accept()
		assert.NoError(t, err)
		assert.NoError(t, Send(viaconn.(*net.UnixConn), sentLn.(*net.TCPListener)))
	}()

	clientConn, err := net.Dial("unix", uds)
	assert.NoError(t, err)

	receivedLn, err := Recv(clientConn.(*net.UnixConn), "127.0.0.1")
	assert.NoError(t, err)
	assert.NoError(t, receivedLn.Close())
}
