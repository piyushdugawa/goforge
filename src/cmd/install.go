package cmd

import (
	"GoForge/utils"

	"github.com/spf13/cobra"
)

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Install project as a program in GOBIN",
	Long:  "Build the project for the host OS and install the binary to $GOBIN.",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		utils.Build()
		utils.Install()
	},
}

func init() {
	rootCmd.AddCommand(installCmd)
}
