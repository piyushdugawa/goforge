package cmd

import (
	"GoForge/utils"

	"github.com/spf13/cobra"
)

var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Build the Go project",
	Long:  "Build the Go project in the current directory and output the executable.",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		utils.Build()
	},
}

func init() {
	rootCmd.AddCommand(buildCmd)
}
