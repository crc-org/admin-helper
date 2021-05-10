package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var Version = "dev"

func main() {
	rootCmd := &cobra.Command{
		Use:          "admin-helper",
		Version:      Version,
		SilenceUsage: true,
	}

	rootCmd.AddCommand(Add, Remove, Clean, Contains)

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
