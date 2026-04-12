package cmd

import (
	"github.com/aliozdemir13/Lumina/internal"

	"github.com/spf13/cobra"
)

var nflCmd = &cobra.Command{
	Use:   "nfl", // The word the user types: 'lumina nfl'
	Short: "Get NFL scores",
	Run: func(cmd *cobra.Command, args []string) {
		scoresService := internal.Score{URL: "https://site.api.espn.com/apis/site/v2/sports"}

		results, err := scoresService.FetchResults("/football/nfl/scoreboard")
		internal.PrintTeamSportsScores(results, err)
	},
}

func init() {
	// This "plugs" the nfl command into the root command
	rootCmd.AddCommand(nflCmd)
}
