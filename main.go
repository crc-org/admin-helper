package main

import (
	"fmt"
	"os"

	"github.com/code-ready/admin-helper/cmd"
	"github.com/spf13/cobra"
)

var Version = "dev"

func main() {
	rootCmd := &cobra.Command{
		Use:     "admin-helper",
		Version: Version,
	}

	rootCmd.AddCommand(cmd.Commands()...)

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
