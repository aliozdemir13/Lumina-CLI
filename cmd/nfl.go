// Package cmd handles the command line interaction, this class handles nfl commands
package cmd

import (
	"github.com/aliozdemir13/Lumina/internal"

	"github.com/spf13/cobra"
)

var nflCmd = &cobra.Command{
	Use:   "nfl", // The word the user types: 'lumina nfl'
	Short: "Get NFL scores",
	Run: func(_ *cobra.Command, _ []string) {
		var scoresService internal.Score

		results, err := scoresService.FetchResults("https://site.api.espn.com/apis/site/v2/sports/football/nfl/scoreboard")
		internal.PrintTeamSportsScores(results, err)
	},
}

func init() {
	// This "plugs" the nfl command into the root command
	rootCmd.AddCommand(nflCmd)
}
