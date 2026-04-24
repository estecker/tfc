/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/spf13/cobra"
	"tfc/internal/tfc"
)

// organizationCmd represents the organization command
var organizationCmd = &cobra.Command{
	Use:   "organization",
	Short: "Information about the organization",
	Long:  `Print membership of the organization. Use backup command to save output to files`,
	Run: func(cmd *cobra.Command, args []string) {
		tfc.OrganizationCmd(org)
	},
}

func init() {
	rootCmd.AddCommand(organizationCmd)
	organizationCmd.PersistentFlags().StringVar(&org, "org", "", "A path to write everything to")
	_ = organizationCmd.MarkPersistentFlagRequired("org")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// organizationCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// organizationCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
