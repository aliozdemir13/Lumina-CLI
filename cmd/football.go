// Package cmd handles the command line interaction, this class handles football commands
package cmd

import (
	"time"

	"github.com/aliozdemir13/Lumina/internal"

	"github.com/spf13/cobra"
)

var footballCmd = &cobra.Command{
	Use:   "football [league]",
	Short: "Get football scores for a specific league",
	Args:  cobra.MinimumNArgs(1), // Ensures the user types a league name
	Run: func(cmd *cobra.Command, args []string) {
		league := args[0] // args[0] is the first word after 'football' to decide the slug
		var scoresService internal.Score

		// Map simple names to ESPN slugs
		mapping := map[string]string{
			"ger": "ger.1",
			"tur": "tur.1",
			"ita": "ita.1",
			"esp": "esp.1",
			"cl":  "uefa.champions",
			"eul": "uefa.europa",
		}

		slug, exists := mapping[league]
		if !exists {
			slug = league // Allow raw slugs too
		}
		filters := ""
		if slug != "uefa.champions" && slug != "uefa.europa" { // to distinct league level date structure (weekend vs non-weekend)
			weeks, _ := cmd.Flags().GetInt("weeks")

			// Get current time
			now := time.Now()

			// Calculate how many days we are past Monday
			// (now.Weekday() is 0 for Sunday, 1 for Monday... 6 for Saturday)
			daysSinceMonday := int(now.Weekday()) - 1
			if daysSinceMonday == -1 { // If today is Sunday, Monday was 6 days ago
				daysSinceMonday = 6
			}

			// Calculate the Monday and Sunday offsets for the requested week
			// weeks=0 is current week, weeks=1 is last week, etc.
			mondayOffset := -daysSinceMonday - (weeks * 7)
			sundayOffset := mondayOffset + 6

			// Get the dates from your internal helper
			start := internal.GetEspnDate(mondayOffset)
			end := internal.GetEspnDate(sundayOffset)

			filters = "?dates=" + start + "-" + end
		}
		results, err := scoresService.FetchResults("https://site.api.espn.com/apis/site/v2/sports/soccer/" + slug + "/scoreboard" + filters)
		internal.PrintTeamSportsScores(results, err)
	},
}

func init() {
	// This "plugs" the football command into the root command
	footballCmd.Flags().IntP("weeks", "w", 0, "How many weeks back to check scores")
	rootCmd.AddCommand(footballCmd)
}
