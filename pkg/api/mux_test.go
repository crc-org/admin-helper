package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/code-ready/admin-helper/pkg/client"
	"github.com/code-ready/admin-helper/pkg/constants"
	"github.com/stretchr/testify/assert"
)

func TestMux(t *testing.T) {
	ts := httptest.NewServer(Mux(nil))
	defer ts.Close()

	client := client.New(http.DefaultClient, ts.URL)
	version, err := client.Version()
	assert.NoError(t, err)
	assert.Equal(t, constants.Version, version)
}
