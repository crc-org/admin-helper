package main

import (
	"fmt"
	"os"
	"runtime"

	"github.com/crc-org/admin-helper/pkg/constants"
	"github.com/spf13/cobra"
)

func main() {
	// Apply sandbox immediately after startup (macOS only, no-op on other platforms)
	// This restricts the process to only access /etc/hosts and denies network/exec
	if err := applySandbox(); err != nil {
		// Non-fatal: warn but continue for compatibility
		fmt.Fprintf(os.Stderr, "Warning: failed to apply sandbox: %v\n", err)
	}

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
