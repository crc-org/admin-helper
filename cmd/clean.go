package cmd

import (
	"strings"

	"github.com/urfave/cli/v2"
)

func Clean() *cli.Command {
	return &cli.Command{
		Name:   "clean",
		Usage:  "clean all entries added with a particular suffix",
		Action: clean,
	}
}
func clean(c *cli.Context) error {
	args := c.Args()

	if args.Len() < 1 {
		return cli.NewExitError("clean requires at least one domain suffix", 1)
	}

	var suffixes []string
	for _, suffix := range args.Slice() {
		if !strings.HasPrefix(suffix, ".") {
			return cli.NewExitError("suffix should start with a dot", 1)
		}
		suffixes = append(suffixes, suffix)
	}

	hostsFile, err := loadHostsFile()
	if err != nil {
		return err
	}

	var toDelete []string
	for _, line := range hostsFile.Lines {
		for _, host := range line.Hosts {
			for _, suffix := range suffixes {
				if strings.HasSuffix(host, suffix) {
					toDelete = append(toDelete, host)
					break
				}
			}
		}
	}

	for _, host := range toDelete {
		if err := hostsFile.RemoveByHostname(host); err != nil {
			return err
		}
	}
	return hostsFile.Flush()
}
