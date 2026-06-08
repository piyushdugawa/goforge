package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "goforge",
	Short: "Goforge is a minimal forge to build and manage your Go-based projects",
	Long: `Goforge - A minimal forge to build and manage your Go-based projects

For more information, visit: https://github.com/piyushdugawa/goforge`,
	Run: func(cmd *cobra.Command, args []string) {
		// Show help when no subcommand is specified
		cmd.Help()
	},
}

func Execute(version string) {
	rootCmd.Version = version
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
