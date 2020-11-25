package cmd

import (
	"github.com/spf13/cobra"
)

var Remove = &cobra.Command{
	Use:     "remove",
	Aliases: []string{"rm", "r"},
	Short:   "Remove host(s) if exists",
	RunE: func(cmd *cobra.Command, args []string) error {
		return remove(args)
	},
}

func remove(args []string) error {
	if len(args) == 0 {
		return nil
	}

	hostsFile, err := loadHostsFile()
	if err != nil {
		return err
	}

	uniqueHosts := map[string]bool{}
	var hostEntries []string

	for i := 0; i < len(args); i++ {
		uniqueHosts[args[i]] = true
	}

	for key := range uniqueHosts {
		hostEntries = append(hostEntries, key)
	}

	for _, host := range hostEntries {
		if err := hostsFile.RemoveByHostname(host); err != nil {
			return err
		}
	}

	return hostsFile.Flush()
}
