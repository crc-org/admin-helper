package hosts

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/areYouLazy/libhosty"

	"github.com/stretchr/testify/assert"
)

const (
	hostsTemplate = `# Do not remove the following line, or various programs
# that require network functionality will fail.
127.0.0.1        localhost.localdomain localhost
::1              localhost6.localdomain6 localhost6
`
)

func TestAdd(t *testing.T) {
	hostsFile := filepath.Join(t.TempDir(), "hosts")
	assert.NoError(t, os.WriteFile(hostsFile, []byte(`127.0.0.1 entry1`), 0600))

	host := hosts(t, hostsFile)

	assert.NoError(t, host.Add("127.0.0.1", []string{"entry1", "entry2", "entry3"}))
	assert.NoError(t, host.Add("127.0.0.2", []string{"entry4"}))

	content, err := os.ReadFile(hostsFile)
	assert.NoError(t, err)
	assert.Equal(t, "127.0.0.1        entry1"+eol()+crcSection("127.0.0.1        entry1 entry2 entry3", "127.0.0.2        entry4")+eol(), string(content))
}

func TestAddMoreThen9Hosts(t *testing.T) {
	hostsFile := filepath.Join(t.TempDir(), "hosts")
	assert.NoError(t, os.WriteFile(hostsFile, []byte(hostsTemplate), 0600))

	host := hosts(t, hostsFile)

	assert.NoError(t, host.Add("127.0.0.1", []string{"entry1", "entry2", "entry3", "entry3", "entry4", "entry5", "entry6", "entry7", "entry8", "entry9", "entry10"}))

	content, err := os.ReadFile(hostsFile)
	assert.NoError(t, err)
	assert.Equal(t, hostsTemplate+eol()+crcSection("127.0.0.1        entry1 entry10 entry2 entry3 entry4 entry5 entry6 entry7 entry8", "127.0.0.1        entry9")+eol(), string(content))
}

func TestAddMoreThan18Hosts(t *testing.T) {
	hostsFile := filepath.Join(t.TempDir(), "hosts")
	assert.NoError(t, os.WriteFile(hostsFile, []byte(hostsTemplate), 0600))

	host := hosts(t, hostsFile)

	assert.NoError(t, host.Add("127.0.0.1", []string{"entry0"}))
	assert.NoError(t, host.Add("127.0.0.1", []string{"entry1", "entry2", "entry3", "entry3", "entry4", "entry5", "entry6", "entry7", "entry8", "entry9", "entry10", "entry11", "entry12", "entry13", "entry14", "entry15", "entry16", "entry17", "entry18", "entry19", "entry20"}))

	content, err := os.ReadFile(hostsFile)
	assert.NoError(t, err)
	assert.Equal(t, hostsTemplate+eol()+crcSection("127.0.0.1        entry0 entry1 entry10 entry11 entry12 entry13 entry14 entry15 entry16", "127.0.0.1        entry17 entry18 entry19 entry2 entry20 entry3 entry4 entry5 entry6", "127.0.0.1        entry7 entry8 entry9")+eol(), string(content))
}

func TestAddMoreThen9HostsInMultipleLines(t *testing.T) {
	hostsFile := filepath.Join(t.TempDir(), "hosts")
	assert.NoError(t, os.WriteFile(hostsFile, []byte(hostsTemplate+eol()+crcSection("127.0.0.1        entry1 entry10 entry2 entry3 entry4 entry5 entry6 entry7", "127.0.0.1        entry11 entry12 entry13 entry14 entry15 entry16 entry17 entry18")+eol()), 0600))

	host := hosts(t, hostsFile)

	assert.NoError(t, host.Add("127.0.0.1", []string{"entry8", "entry9", "entry10"}))

	content, err := os.ReadFile(hostsFile)
	assert.NoError(t, err)
	assert.Equal(t, hostsTemplate+eol()+crcSection("127.0.0.1        entry1 entry10 entry2 entry3 entry4 entry5 entry6 entry7 entry8", "127.0.0.1        entry11 entry12 entry13 entry14 entry15 entry16 entry17 entry18 entry9")+eol(), string(content))
}

func TestRemove(t *testing.T) {
	hostsFile := filepath.Join(t.TempDir(), "hosts")
	assert.NoError(t, os.WriteFile(hostsFile, []byte(hostsTemplate), 0600))

	host := hosts(t, hostsFile)
	assert.NoError(t, host.Add("127.0.0.1", []string{"entry1", "entry2"}))

	assert.NoError(t, host.Remove([]string{"entry2"}))

	content, err := os.ReadFile(hostsFile)
	assert.NoError(t, err)
	assert.Equal(t, hostsTemplate+eol()+crcSection("127.0.0.1        entry1")+eol(), string(content))
}

func TestClean(t *testing.T) {
	hostsFile := filepath.Join(t.TempDir(), "hosts")
	assert.NoError(t, os.WriteFile(hostsFile, []byte(hostsTemplate+eol()+crcSection("127.0.0.1 entry1.suffix1 entry2.suffix2")), 0600))

	host := hosts(t, hostsFile)

	assert.NoError(t, host.Clean())

	content, err := os.ReadFile(hostsFile)
	assert.NoError(t, err)
	assert.Equal(t, hostsTemplate, string(content))
}

func TestCleanWithoutCrcSection(t *testing.T) {
	hostsFile := filepath.Join(t.TempDir(), "hosts")
	assert.NoError(t, os.WriteFile(hostsFile, []byte(hostsTemplate), 0600))

	host := hosts(t, hostsFile)

	assert.NoError(t, host.Clean())

	content, err := os.ReadFile(hostsFile)
	assert.NoError(t, err)
	assert.Equal(t, hostsTemplate, string(content))
}

func TestContains(t *testing.T) {
	hostsFile := filepath.Join(t.TempDir(), "hosts")
	assert.NoError(t, os.WriteFile(hostsFile, []byte(`127.0.0.1 entry1.suffix1 entry2.suffix2`), 0600))

	host := hosts(t, hostsFile)

	assert.True(t, host.Contains("127.0.0.1", "entry1.suffix1"))
	assert.False(t, host.Contains("127.0.0.2", "entry1.suffix1"))
	assert.False(t, host.Contains("127.0.0.1", "entry1.suffix2"))
}

func TestSuffixFilter(t *testing.T) {
	hostsFile := filepath.Join(t.TempDir(), "hosts")
	assert.NoError(t, os.WriteFile(hostsFile, []byte(`127.0.0.1 localhost localhost.localdomain`), 0600))

	config, _ := libhosty.NewHostsFileConfig(hostsFile)
	file, err := libhosty.InitWithConfig(config)
	assert.NoError(t, err)
	host := Hosts{
		File:       file,
		HostFilter: defaultFilter,
	}

	assert.NoError(t, host.Add("127.0.0.1", []string{"entry1.crc.testing"}))
	assert.NoError(t, host.Add("127.0.0.1", []string{"entry2.nested.crc.testing"}))
	assert.Error(t, host.Add("127.0.0.1", []string{"evildomain #apps.crc.testing"}))
	assert.Error(t, host.Add("127.0.0.1", []string{"host.poison"}))
	assert.Error(t, host.Add("127.0.0.1", []string{"CAPITAL.crc.testing"}))
	assert.Error(t, host.Remove([]string{"localhost"}))
	assert.NoError(t, host.Clean())
}

func TestAddMoreThan9HostsWithFullLine(t *testing.T) {
	hostsFile := filepath.Join(t.TempDir(), "hosts")
	assert.NoError(t, os.WriteFile(hostsFile, []byte(hostsTemplate+eol()+crcSection("127.0.0.1        entry1  entry2 entry3 entry4 entry5 entry6 entry7 entry8 entry9")+eol()), 0600))

	host := hosts(t, hostsFile)

	assert.NoError(t, host.Add("127.0.0.1", []string{"entry10"}))

	content, err := os.ReadFile(hostsFile)
	assert.NoError(t, err)
	assert.Equal(t, hostsTemplate+eol()+crcSection("127.0.0.1        entry1 entry2 entry3 entry4 entry5 entry6 entry7 entry8 entry9", "127.0.0.1        entry10")+eol(), string(content))
}

func TestAddMoreThan9HostsWithOverfullLine(t *testing.T) {
	hostsFile := filepath.Join(t.TempDir(), "hosts")
	assert.NoError(t, os.WriteFile(hostsFile, []byte(hostsTemplate+eol()+crcSection("127.0.0.1        entry1  entry2 entry3 entry4 entry5 entry6 entry7 entry8 entry9 entry10")+eol()), 0600))

	host := hosts(t, hostsFile)

	assert.NoError(t, host.Add("127.0.0.1", []string{"entry11"}))

	content, err := os.ReadFile(hostsFile)
	assert.NoError(t, err)
	assert.Equal(t, hostsTemplate+eol()+crcSection("127.0.0.1        entry1 entry2 entry3 entry4 entry5 entry6 entry7 entry8 entry9", "127.0.0.1        entry10 entry11")+eol(), string(content))
}

func TestRemoveOnOldHostFile(t *testing.T) {
	hostsFile := filepath.Join(t.TempDir(), "hosts")
	assert.NoError(t, os.WriteFile(hostsFile, []byte(hostsTemplate+eol()+"127.0.0.1 entry1 entry2"), 0600))

	host := hosts(t, hostsFile)

	assert.NoError(t, host.Remove([]string{"entry1", "entry2"}))

	content, err := os.ReadFile(hostsFile)
	assert.NoError(t, err)
	assert.Equal(t, hostsTemplate, string(content))
}

func TestRemoveMultipleForwardSameLine(t *testing.T) {
	hostsFile := filepath.Join(t.TempDir(), "hosts")
	assert.NoError(t, os.WriteFile(hostsFile, []byte(hostsTemplate+eol()+crcSection("192.168.130.11   entry1 entry2")), 0600))
	host := hosts(t, hostsFile)

	assert.NoError(t, host.Remove([]string{"entry1", "entry2"}))

	content, err := os.ReadFile(hostsFile)
	assert.NoError(t, err)
	assert.Equal(t, hostsTemplate+eol()+crcSection(), string(content))
}

func TestRemoveMultipleBackwardsSameLine(t *testing.T) {
	hostsFile := filepath.Join(t.TempDir(), "hosts")
	assert.NoError(t, os.WriteFile(hostsFile, []byte(hostsTemplate+eol()+crcSection("192.168.130.11   entry1 entry2")), 0600))
	host := hosts(t, hostsFile)

	assert.NoError(t, host.Remove([]string{"entry2", "entry1"}))

	content, err := os.ReadFile(hostsFile)
	assert.NoError(t, err)
	assert.Equal(t, hostsTemplate+eol()+crcSection(), string(content))
}

func TestRemoveMultipleLines(t *testing.T) {
	hostsFile := filepath.Join(t.TempDir(), "hosts")
	assert.NoError(t, os.WriteFile(hostsFile, []byte(hostsTemplate+eol()+crcSection("192.168.130.11   api.crc.testing", "192.168.130.11    oauth-openshift.apps-crc.testing")), 0600))
	host := hosts(t, hostsFile)

	assert.NoError(t, host.Remove([]string{"api.crc.testing", "oauth-openshift.apps-crc.testing"}))

	content, err := os.ReadFile(hostsFile)
	assert.NoError(t, err)
	assert.Equal(t, hostsTemplate+eol()+crcSection(), string(content))
}

func TestRemoveMultipleNoCrcSection(t *testing.T) {
	hostsFile := filepath.Join(t.TempDir(), "hosts")
	assert.NoError(t, os.WriteFile(hostsFile, []byte("192.168.130.11   entry1 entry2"+eol()+"192.168.130.11   entry3 entry4"), 0600))
	host := hosts(t, hostsFile)

	assert.NoError(t, host.Remove([]string{"entry1", "entry2", "entry3", "entry4"}))

	content, err := os.ReadFile(hostsFile)
	assert.NoError(t, err)
	assert.Equal(t, "", string(content))
}

func hosts(t *testing.T, hostsFile string) Hosts {
	config, _ := libhosty.NewHostsFileConfig(hostsFile)
	file, err := libhosty.InitWithConfig(config)
	assert.NoError(t, err)
	return Hosts{
		File: file,
		HostFilter: func(_ string) bool {
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

func crcSection(lines ...string) string {
	var content = ""
	if len(lines) != 0 {
		content = strings.Join(lines, eol()) + eol()
	}
	return fmt.Sprintf("# Added by CRC"+eol()+"%s"+"# End of CRC section", content)
}
