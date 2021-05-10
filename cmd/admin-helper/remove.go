package main

import (
	"github.com/code-ready/admin-helper/pkg/hosts"
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

	hosts, err := hosts.New()
	if err != nil {
		return err
	}
	return hosts.Remove(args)
}
