package hosts

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/goodhosts/hostsfile"
	"github.com/stretchr/testify/assert"
)

func TestAdd(t *testing.T) {
	dir, err := ioutil.TempDir("", "hosts")
	assert.NoError(t, err)
	defer os.RemoveAll(dir)

	hostsFile := filepath.Join(dir, "hosts")
	assert.NoError(t, ioutil.WriteFile(hostsFile, []byte(`127.0.0.1 entry1`), 0600))

	host := hosts(t, hostsFile)

	assert.NoError(t, host.Add("127.0.0.1", []string{"entry1", "entry2", "entry3"}))
	assert.NoError(t, host.Add("127.0.0.2", []string{"entry4"}))

	content, err := ioutil.ReadFile(hostsFile)
	assert.NoError(t, err)
	assert.Equal(t, "127.0.0.1 entry1 entry2 entry3"+eol()+"127.0.0.2 entry4"+eol(), string(content))
}

func TestRemove(t *testing.T) {
	dir, err := ioutil.TempDir("", "hosts")
	assert.NoError(t, err)
	defer os.RemoveAll(dir)

	hostsFile := filepath.Join(dir, "hosts")
	assert.NoError(t, ioutil.WriteFile(hostsFile, []byte(`127.0.0.1 entry1 entry2`), 0600))

	host := hosts(t, hostsFile)

	assert.NoError(t, host.Remove([]string{"entry2"}))

	content, err := ioutil.ReadFile(hostsFile)
	assert.NoError(t, err)
	assert.Equal(t, "127.0.0.1 entry1"+eol(), string(content))
}

func TestClean(t *testing.T) {
	dir, err := ioutil.TempDir("", "hosts")
	assert.NoError(t, err)
	defer os.RemoveAll(dir)

	hostsFile := filepath.Join(dir, "hosts")
	assert.NoError(t, ioutil.WriteFile(hostsFile, []byte(`127.0.0.1 entry1.suffix1 entry2.suffix2`), 0600))

	host := hosts(t, hostsFile)

	assert.NoError(t, host.Clean([]string{".suffix1"}))

	content, err := ioutil.ReadFile(hostsFile)
	assert.NoError(t, err)
	assert.Equal(t, "127.0.0.1 entry2.suffix2"+eol(), string(content))
}

func TestContains(t *testing.T) {
	dir, err := ioutil.TempDir("", "hosts")
	assert.NoError(t, err)
	defer os.RemoveAll(dir)

	hostsFile := filepath.Join(dir, "hosts")
	assert.NoError(t, ioutil.WriteFile(hostsFile, []byte(`127.0.0.1 entry1.suffix1 entry2.suffix2`), 0600))

	host := hosts(t, hostsFile)

	assert.True(t, host.Contains("127.0.0.1", "entry1.suffix1"))
	assert.False(t, host.Contains("127.0.0.2", "entry1.suffix1"))
	assert.False(t, host.Contains("127.0.0.1", "entry1.suffix2"))
}

func TestSuffixFilter(t *testing.T) {
	dir, err := ioutil.TempDir("", "hosts")
	assert.NoError(t, err)
	defer os.RemoveAll(dir)

	hostsFile := filepath.Join(dir, "hosts")
	assert.NoError(t, ioutil.WriteFile(hostsFile, []byte(`127.0.0.1 localhost`), 0600))

	file, err := hostsfile.NewCustomHosts(hostsFile)
	assert.NoError(t, err)
	host := Hosts{
		File:       &file,
		HostFilter: defaultFilter,
	}

	assert.NoError(t, host.Add("127.0.0.1", []string{"entry1.crc.testing"}))
	assert.NoError(t, host.Add("127.0.0.1", []string{"entry2.nested.crc.testing"}))
	assert.Error(t, host.Add("127.0.0.1", []string{"evildomain #apps.crc.testing"}))
	assert.Error(t, host.Add("127.0.0.1", []string{"host.poison"}))
	assert.Error(t, host.Add("127.0.0.1", []string{"CAPITAL.crc.testing"}))
	assert.Error(t, host.Remove([]string{"localhost"}))
}

func hosts(t *testing.T, hostsFile string) Hosts {
	file, err := hostsfile.NewCustomHosts(hostsFile)
	assert.NoError(t, err)
	return Hosts{
		File: &file,
		HostFilter: func(s string) bool {
			return true
		},
	}
}

func eol() string {
	if runtime.GOOS == "windows" {
		return "\r\n"
	}
	return "\n"
}
