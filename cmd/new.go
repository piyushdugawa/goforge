package cmd

import (
	"github.com/piyushdugawa/girocco/utils"

	"github.com/spf13/cobra"
)

var newCmd = &cobra.Command{
	Use:   "new [pkg-name]",
	Short: "Initialize a new girocco project",
	Long:  "Create a new Go project in the current directory and initialize go.mod.",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		pkgName := ""
		if len(args) > 0 {
			pkgName = args[0]
		}
		utils.New(pkgName)
	},
}

func init() {
	rootCmd.AddCommand(newCmd)
}
