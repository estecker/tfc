package cmd

import (
    "github.com/estecker/tfc/internal/tfc"
    "github.com/spf13/cobra"
)

var backupCmd = &cobra.Command{
    Use:   "backup",
    Short: "Backup into a folder",
    Long:  "Download states and variables and write to a folder",
}

var backupStatesCmd = &cobra.Command{
    Use:   "states",
    Short: "Backup states into a folder",
    Long:  "Download states and variables and write to a folder",
    Run: func(cmd *cobra.Command, args []string) {
        if error := tfc.BackupStatesCmd(org, Folder); error != nil {
            panic(error)
        }
    },
}
var backupVariablesCmd = &cobra.Command{
    Use:   "variables",
    Short: "Backup variables into a folder",
    Long:  "Download states and variables and write to a folder",
    Run: func(cmd *cobra.Command, args []string) {
        if err := tfc.BackupVariablesCmd(org, Folder); err != nil {
            panic(err)
        }
    },
}
var backupWorkspaceCmd = &cobra.Command{
    Use:   "workspace",
    Short: "Backup workspace into a folder",
    Long:  "Download workspace and write to a folder",
    Run: func(cmd *cobra.Command, args []string) {
        if err := tfc.BackupWorkspacesCmd(org, Folder); err != nil {
            panic(err)
        }
    },
}

var Folder string

func init() {
    backupCmd.AddCommand(backupStatesCmd)
    backupCmd.AddCommand(backupVariablesCmd)
    backupCmd.AddCommand(backupWorkspaceCmd)

    backupCmd.PersistentFlags().StringVar(&Folder, "folder", "", "A path to write everything to")
    _ = backupCmd.MarkPersistentFlagRequired("folder")

    backupCmd.PersistentFlags().StringVar(&org, "org", "", "A path to write everything to")
    _ = backupCmd.MarkPersistentFlagRequired("org")

}
