package main

import (
	"fmt"

	"github.com/code-ready/admin-helper/pkg/hosts"
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
	hosts, err := hosts.New()
	if err != nil {
		return err
	}
	return hosts.Add(args[0], args[1:])
}
