package main

import (
	"fmt"

	"github.com/crc-org/admin-helper/pkg/hosts"
	"github.com/spf13/cobra"
)

var Clean = &cobra.Command{
	Use:   "clean",
	Short: "Clean all entries added with a particular suffix",
	RunE: func(_ *cobra.Command, args []string) error {
		return clean(args)
	},
}

func clean(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("clean requires at least one domain suffix")
	}

	hosts, err := hosts.New()
	if err != nil {
		return err
	}
	return hosts.Clean()
}
