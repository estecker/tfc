package cmd

import (
	"github.com/estecker/tfc/internal/tfc"
	"github.com/spf13/cobra"
)

var workspacesCmd = &cobra.Command{
	Use:   "workspaces",
	Short: "Search for a workspace",
	Long:  "Example search strings here",
	Run: func(cmd *cobra.Command, args []string) {
		_ = tfc.WorkspaceListCommand(org, verbose, Search, Tags, ExcludeTags, WildcardName, Status)
	},
}
var workspacesPlanningCmd = &cobra.Command{
	Use:   "planning",
	Short: "List the workspaces that are currently planning",
	Long:  "Should catch folks running plans from CLI",
	Run: func(cmd *cobra.Command, args []string) {
		_ = tfc.WorkspacePlanning(org)
	},
}
var workspacesLockCmd = &cobra.Command{
	Use:   "lock",
	Short: "Lock a workspace by its ID.",
	Run: func(cmd *cobra.Command, args []string) {
		_ = tfc.WorkspaceLock(ID)
	},
}
var workspacesUnlockCmd = &cobra.Command{
	Use:   "unlock",
	Short: "Unlock a workspace by its ID.",
	Run: func(cmd *cobra.Command, args []string) {
		_ = tfc.WorkspaceUnlock(ID)
	},
}
var Tags string
var Search string
var ExcludeTags string
var WildcardName string
var Status string

var ID string

func init() {

	workspacesCmd.PersistentFlags().StringVarP(&org, "org", "o", "", "Terraform Cloud Organization (required)")

	workspacesCmd.Flags().StringVar(&Search, "search", "", "A search string (partial workspace name) used to filter the results.")
	workspacesCmd.Flags().StringVar(&Tags, "tags", "", "A search string (comma-separated tag names) used to filter the results.")
	workspacesCmd.Flags().StringVar(&ExcludeTags, "exclude-tags", "", "A search string (comma-separated tag names to exclude) used to filter the results.")
	workspacesCmd.Flags().StringVar(&WildcardName, "wildcardname", "", "A search on substring matching to filter the results.")
	workspacesCmd.Flags().StringVar(&Status, "status", "", "Also filter by workspace CurrentRun RunStatus options https://pkg.go.dev/github.com/hashicorp/go-tfe#RunStatus")
	//https://github.com/spf13/pflag/issues/236

	workspacesCmd.AddCommand(workspacesPlanningCmd)
	workspacesCmd.AddCommand(workspacesLockCmd)
	workspacesCmd.AddCommand(workspacesUnlockCmd)

	workspacesLockCmd.PersistentFlags().StringVar(&ID, "id", "", "Workspace ID")
	_ = workspacesLockCmd.MarkFlagRequired("id")

	workspacesUnlockCmd.PersistentFlags().StringVar(&ID, "id", "", "Workspace ID")
	_ = workspacesUnlockCmd.MarkFlagRequired("id")

	err := workspacesCmd.MarkPersistentFlagRequired("org")
	if err != nil {
		return
	}
}
