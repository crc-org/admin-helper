package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

var Clean = &cobra.Command{
	Use:   "clean",
	Short: "Clean all entries added with a particular suffix",
	RunE: func(cmd *cobra.Command, args []string) error {
		return clean(args)
	},
}

func clean(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("clean requires at least one domain suffix")
	}

	var suffixes []string
	for _, suffix := range args {
		if !strings.HasPrefix(suffix, ".") {
			return fmt.Errorf("suffix should start with a dot")
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
