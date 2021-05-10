package main

import (
	"fmt"

	"github.com/code-ready/admin-helper/pkg/hosts"
	"github.com/spf13/cobra"
)

var Contains = &cobra.Command{
	Use:   "contains",
	Short: "Check if an ip and host are present in hosts file",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		return contains(args)
	},
}

func contains(args []string) error {
	hosts, err := hosts.New()
	if err != nil {
		return err
	}
	if hosts.Contains(args[0], args[1]) {
		fmt.Println("true")
	} else {
		fmt.Println("false")
	}
	return nil
}
