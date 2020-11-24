package cmd

import (
	"fmt"

	"github.com/goodhosts/hostsfile"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func Commands() []*cobra.Command {
	return []*cobra.Command{
		Add,
		Remove,
		Clean,
	}
}

func loadHostsFile() (hostsfile.Hosts, error) {
	logrus.Debugf("loading default hosts file: %s", hostsfile.HostsFilePath)
	file, err := hostsfile.NewHosts()
	if err != nil {
		return file, err
	}
	if !file.IsWritable() {
		return file, fmt.Errorf("host file not writable, try running with elevated privileges")
	}
	return file, nil
}
