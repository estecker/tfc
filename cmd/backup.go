package cmd

import (
	"github.com/spf13/cobra"
	"tfc/internal/tfc"
)

var backupCmd = &cobra.Command{
	Use:   "backup",
	Short: "Backup into a folder",
	Long:  "Download states and variables and write to a folder",
}

var backupStatesCmd = &cobra.Command{
	Use:   "states",
	Short: "Backup states",
	Long:  "Download states and variables and write to a folder",
	Run: func(cmd *cobra.Command, args []string) {
		if error := tfc.BackupStatesCmd(org, Folder); error != nil {
			panic(error)
		}
	},
}
var backupVariablesCmd = &cobra.Command{
	Use:   "variables",
	Short: "Backup variables",
	Long:  "Download states and variables and write to a folder",
	Run: func(cmd *cobra.Command, args []string) {
		if err := tfc.BackupVariablesCmd(org, Folder); err != nil {
			panic(err)
		}
	},
}
var backupWorkspaceCmd = &cobra.Command{
	Use:   "workspace",
	Short: "Backup workspace",
	Long:  "Download workspace and write to a folder",
	Run: func(cmd *cobra.Command, args []string) {
		if err := tfc.BackupWorkspacesCmd(org, Folder); err != nil {
			panic(err)
		}
	},
}
var backupOrgMembershipCmd = &cobra.Command{
	Use:   "org-membership",
	Short: "Backup membership of an organization",
	Long:  "Backup membership of an organization into a folder",
	Run: func(cmd *cobra.Command, args []string) {
		if err := tfc.BackupOrgMembershipCmd(org, Folder); err != nil {
			panic(err)
		}
	},
}
var Folder string

func init() {
	rootCmd.AddCommand(backupCmd)
	backupCmd.AddCommand(backupStatesCmd)
	backupCmd.AddCommand(backupVariablesCmd)
	backupCmd.AddCommand(backupWorkspaceCmd)
	backupCmd.AddCommand(backupOrgMembershipCmd)

	backupCmd.PersistentFlags().StringVar(&Folder, "folder", "", "A path to write everything to")
	_ = backupCmd.MarkPersistentFlagRequired("folder")

	backupCmd.PersistentFlags().StringVar(&org, "org", "", "A path to write everything to")
	_ = backupCmd.MarkPersistentFlagRequired("org")

}
