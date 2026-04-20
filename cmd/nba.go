// Package cmd handles the command line interaction, this class handles nba commands
package cmd

import (
	"github.com/aliozdemir13/Lumina/internal"

	"github.com/spf13/cobra"
)

var nbaCmd = &cobra.Command{
	Use:   "nba", // The word the user types: 'lumina nba'
	Short: "Get NBA scores",
	Run: func(cmd *cobra.Command, _ []string) {
		var scoresService internal.Score
		// Get the 'days' flag using cmd
		days, _ := cmd.Flags().GetInt("days")
		start, _ := cmd.Flags().GetInt("start")

		// Calculate range based on the flag
		// If days is 1, it shows yesterday to today.
		// If days is 5, it shows 5 days ago to today.
		today := internal.GetEspnDate(start)
		startDate := internal.GetEspnDate(start - days)

		results, err := scoresService.FetchResults("https://site.api.espn.com/apis/site/v2/sports/basketball/nba/scoreboard?dates=" + startDate + "-" + today) // only last 2 days logic for keeping list short
		internal.PrintTeamSportsScores(results, err)
	},
}

func init() {
	// This "plugs" the nba command into the root command
	nbaCmd.Flags().IntP("days", "d", 1, "How many days back to check scores")
	nbaCmd.Flags().IntP("start", "s", 0, "How many days back from today to check scores (regular numbers for future, minus numbers for past)")

	rootCmd.AddCommand(nbaCmd)
}
