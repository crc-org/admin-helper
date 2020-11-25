package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var Add = &cobra.Command{
	Use:     "add",
	Aliases: []string{"a"},
	Short:   "Add an entry to the hostsfile",
	RunE: func(cmd *cobra.Command, args []string) error {
		return add(args)
	},
}

func add(args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("adding to hosts file requires an ip and a hostname")
	}

	hostsFile, err := loadHostsFile()
	if err != nil {
		return err
	}

	ip := args[0]
	uniqueHosts := map[string]bool{}
	var hostEntries []string

	for i := 1; i < len(args); i++ {
		uniqueHosts[args[i]] = true
	}

	for key := range uniqueHosts {
		hostEntries = append(hostEntries, key)
	}

	if err := hostsFile.Add(ip, hostEntries...); err != nil {
		return err
	}
	return hostsFile.Flush()
}
