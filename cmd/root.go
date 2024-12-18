package cmd

import (
	"github.com/spf13/cobra"
	"os"
)

var org string
var verbose bool
var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "tfc",
	Short: "Terraform Cloud CLI",
	Long:  `A tool to easily make bulk changes to Terraform Cloud via the API.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "config file (default is $HOME/.tfc.yaml)")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")

	rootCmd.AddCommand(workspacesCmd)
	rootCmd.AddCommand(backupCmd)
	rootCmd.AddCommand(runsCmd)
}
