// Package cmd handles the command line interaction, this class is the root
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// This is the base command when running 'lumina'
var rootCmd = &cobra.Command{
	Use:   "lumina",
	Short: "Lumina is a CLI for live sports scores",
	Long:  `A high-performance CLI tool to track NBA, F1, and Soccer scores directly in your terminal.`,
	Run: func(_ *cobra.Command, _ []string) {
		// If the user didn't type a subcommand (like 'nba'),
		// interactive logic to be added.
		// RunInteractiveMenu()
	},
}

// Execute is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
