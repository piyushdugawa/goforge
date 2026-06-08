package cmd

import (
	"GoForge/utils"

	"github.com/spf13/cobra"
)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run the current project",
	Long:  "Run the compiled project binary. Builds the binary first if not already present.",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		utils.Run()
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
}
