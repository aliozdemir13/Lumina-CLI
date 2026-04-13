package cmd

import (
	"strconv"
	"time"

	"github.com/aliozdemir13/Lumina/internal"

	"github.com/spf13/cobra"
)

var weeksOffsetRacing int // for being able to check results weekly basis
var showAll bool          // for being able to get all results from year to date

var racingCmd = &cobra.Command{
	Use:   "racing [motorsportName]",
	Short: "Get racing results for the selected motorsport",
	Args:  cobra.MinimumNArgs(1), // Ensures the user types a motorsport name
	Run: func(cmd *cobra.Command, args []string) {
		league := args[0] // args[0] is the first word after 'racing'
		var racingService internal.Results

		// Map simple names to ESPN slugs
		mapping := map[string]string{
			"f1":     "f1",
			"indy":   "irl",
			"nascar": "nascar-premier",
		}

		slug, exists := mapping[league]
		if !exists {
			slug = league // Allow raw slugs too
		}
		today := internal.GetEspnDate(0 + (weeksOffsetRacing * -7))
		yesterday := internal.GetEspnDate(-1 + (weeksOffsetRacing * -7))
		if showAll { // dynamic year to date logic
			year, _, _ := time.Now().Date()
			yesterday = strconv.Itoa(year) + "0101"
		}
		results, err := racingService.FetchResults("https://site.api.espn.com/apis/site/v2/sports/racing/" + slug + "/scoreboard?dates=" + yesterday + "-" + today)
		internal.PrintRacingTable(results, err)
	},
}

func init() {
	// This "plugs" the racing command into the root command
	racingCmd.Flags().IntVarP(&weeksOffsetRacing, "weeks", "w", 0, "How many weeks back to check results")
	racingCmd.Flags().BoolVarP(&showAll, "all", "a", false, "Show all results for the current year")
	rootCmd.AddCommand(racingCmd)
}
