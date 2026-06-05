package main

import (
	"fmt"

	"github.com/crc-org/admin-helper/pkg/hosts"
	"github.com/crc-org/admin-helper/pkg/logging"
	"github.com/spf13/cobra"
)

var Add = &cobra.Command{
	Use:     "add",
	Aliases: []string{"a"},
	Short:   "Add an entry to the hostsfile",
	RunE: func(_ *cobra.Command, args []string) error {
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
	err = hosts.Add(args[0], args[1:])
	logger := logging.GetLogger()
	logger.LogModification(logging.Modification{
		Operation: "add",
		IP:        args[0],
		Hosts:     args[1:],
		Caller:    cliCaller(),
		Error:     err,
	})
	return err
}
