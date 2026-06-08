package cmd

import (
	"GoForge/utils"

	"github.com/spf13/cobra"
)

var cleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "Removes all builds and temporary files",
	Long:  "Removes the entire output build directory recursively.",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		utils.Clean()
	},
}

func init() {
	rootCmd.AddCommand(cleanCmd)
}
