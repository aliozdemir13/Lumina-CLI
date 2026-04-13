package cmd

import (
	"github.com/aliozdemir13/Lumina/internal"

	"github.com/spf13/cobra"
)

var nbaCmd = &cobra.Command{
	Use:   "nba", // The word the user types: 'lumina nba'
	Short: "Get NBA scores",
	Run: func(cmd *cobra.Command, args []string) {
		var scoresService internal.Score
		today := internal.GetEspnDate(0)
		yesterday := internal.GetEspnDate(-1)

		results, err := scoresService.FetchResults("https://site.api.espn.com/apis/site/v2/sports/basketball/nba/scoreboard?dates=" + yesterday + "-" + today) // only last 2 days logic for keeping list short
		internal.PrintTeamSportsScores(results, err)
	},
}

func init() {
	// This "plugs" the nba command into the root command
	rootCmd.AddCommand(nbaCmd)
}
