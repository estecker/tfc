package cmd

import (
	"fmt"
	"tfc/internal/tfc"

	"github.com/spf13/cobra"
)

// teamsCmd represents the teams command
var teamsCmd = &cobra.Command{
	Use:   "teams",
	Short: "List teams in an organization",
	Long:  `Not really used. So implementation is light `,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("teams called")
		tfc.TeamsCmd(org)
	},
}

func init() {
	rootCmd.AddCommand(teamsCmd)

}
