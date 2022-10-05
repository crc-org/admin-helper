package main

import (
	"fmt"
	"os"
	"runtime"

	"github.com/crc-org/admin-helper/pkg/constants"
	"github.com/spf13/cobra"
)

func main() {
	rootCmd := &cobra.Command{
		Use:          "admin-helper",
		Version:      constants.Version,
		SilenceUsage: true,
	}

	rootCmd.AddCommand(Add, Remove, Clean, Contains)
	if runtime.GOOS == "windows" {
		rootCmd.AddCommand(InstallDaemon, UninstallDaemon, Daemon)
	}

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
