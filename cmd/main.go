package cmd

import (
	"github.com/goodhosts/hostsfile"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

func Commands() []*cli.Command {
	return []*cli.Command{
		Add(),
		Remove(),
		Clean(),
	}
}

func loadHostsfile() (hostsfile.Hosts, error) {
	logrus.Debugf("loading default hosts file: %s\n", hostsfile.HostsFilePath)
	hfile, err := hostsfile.NewHosts()

	if err != nil {
		return hfile, cli.NewExitError(err, 1)
	}

	if !hfile.IsWritable() {
		return hfile, cli.NewExitError("Host file not writable. Try running with elevated privileges.", 1)
	}

	return hfile, nil
}
