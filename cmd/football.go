package cmd

import (
	"github.com/aliozdemir13/Lumina/internal"

	"github.com/spf13/cobra"
)

var weeksOffsetFootball int

var footballCmd = &cobra.Command{
	Use:   "football [league]",
	Short: "Get football scores for a specific league",
	Args:  cobra.MinimumNArgs(1), // Ensures the user types a league name
	Run: func(cmd *cobra.Command, args []string) {
		league := args[0] // args[0] is the first word after 'football' to decide the slug
		scoresService := internal.Score{URL: "https://site.api.espn.com/apis/site/v2/sports"}

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
		if slug != "uefa.champions" && slug != "uefa.europa" { // to distinct league level date structure (weekend vs non-weekend)
			today := internal.GetEspnDate(0 + (weeksOffsetFootball * -7))
			yesterday := internal.GetEspnDate(-1 + (weeksOffsetFootball * -7))
			results, err := scoresService.FetchResults("/soccer/" + slug + "/scoreboard?dates=" + yesterday + "-" + today)
			internal.PrintTeamSportsScores(results, err)
		} else {
			results, err := scoresService.FetchResults("/soccer/" + slug + "/scoreboard")
			internal.PrintTeamSportsScores(results, err)
		}

	},
}

func init() {
	// This "plugs" the football command into the root command
	footballCmd.Flags().IntVarP(&weeksOffsetFootball, "weeks", "w", 0, "How many weeks back to check scores")
	rootCmd.AddCommand(footballCmd)
}
