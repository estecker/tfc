package cmd

import (
	"github.com/estecker/tfc/internal/tfc"
	"github.com/spf13/cobra"
)

var Workspaces []string
var Workspace string
var Force bool
var CurrentRun bool
var RunID string

var runsCmd = &cobra.Command{
	Use:   "runs",
	Short: "Workspace runs interface",
	Long:  "https://pkg.go.dev/github.com/hashicorp/go-tfe#Runs",
	Run: func(cmd *cobra.Command, args []string) {
		tfc.RunsCmd(Workspaces, Force, CurrentRun)
	},
}

var runsDiscardCmd = &cobra.Command{
	Use:   "discard",
	Short: "Discard the current pending run in a workspace",
	Long:  "https://pkg.go.dev/github.com/hashicorp/go-tfe#Runs",
	Run: func(cmd *cobra.Command, args []string) {
		tfc.RunsDiscardCmd(Workspace)
	},
}
var runsCancelCmd = &cobra.Command{
	Use:   "cancel",
	Short: "Cancel a run by its ID.",
	Long:  "https://pkg.go.dev/github.com/hashicorp/go-tfe#Runs",
	Run: func(cmd *cobra.Command, args []string) {
		tfc.RunCancelCmd(RunID)
	},
}

var runsLogCmd = &cobra.Command{
	Use:   "log",
	Short: "Show log for run-id",
	Long:  "https://pkg.go.dev/github.com/hashicorp/go-tfe#Runs",
	Run: func(cmd *cobra.Command, args []string) {
		tfc.RunsLogCmd(RunID)
	},
}

func init() {
	runsCmd.Flags().StringSliceVarP(&Workspaces, "workspace-ids", "w", []string{}, "A comma seperated list of workspace ID's to run")
	_ = runsCmd.MarkFlagRequired("workspace-ids")
	runsCmd.Flags().BoolVar(&Force, "force", false, "https://developer.hashicorp.com/terraform/cloud-docs/api-docs/run#forcefully-execute-a-run")
	runsCmd.Flags().BoolVar(&CurrentRun, "current-run", false, "Look at current run first, before starting a new run")
	_ = runsCmd.MarkFlagRequired("workspace-id")

	runsCmd.AddCommand(runsDiscardCmd)
	runsDiscardCmd.Flags().StringVarP(&Workspace, "workspace-id", "w", "", "Workspace ID of which to look for a pending run to discard")
	_ = runsDiscardCmd.MarkFlagRequired("workspace-id")

	runsCmd.AddCommand(runsCancelCmd)
	runsCancelCmd.Flags().StringVarP(&RunID, "run-id", "r", "", "Run")
	_ = runsCancelCmd.MarkFlagRequired("run-id")

	runsCmd.AddCommand(runsLogCmd)
	runsLogCmd.Flags().StringVarP(&RunID, "run-id", "i", "", "Run ID to show logs of")
	_ = runsLogCmd.MarkFlagRequired("run-id")
}
