package cmd

import (
	"github.com/piyushdugawa/goforge/utils"

	"github.com/spf13/cobra"
)

var removeCmd = &cobra.Command{
	Use:   "remove",
	Short: "Remove the installed program from GOBIN",
	Long:  "Remove the installed binary from GOBIN.",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		utils.Remove()
	},
}

func init() {
	rootCmd.AddCommand(removeCmd)
}
