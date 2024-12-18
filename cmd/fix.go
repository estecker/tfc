package cmd

import (
	"github.com/estecker/tfc/internal/tfc"
	"github.com/spf13/cobra"
)

var Problem string
var Option string

var fixCmd = &cobra.Command{
	Use:   "fix",
	Short: "Fix common issues with Terraform Cloud",
	Long:  "Rather than having to for each multiple commands, run fix once",
	Run: func(cmd *cobra.Command, args []string) {
		tfc.FixCmd()
	},
}

func init() {
	rootCmd.AddCommand(fixCmd)
}
