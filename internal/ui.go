package internal

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/jedib0t/go-pretty/table"
)

func GetEspnDate(offsetDays int) string {
	// Get current time and add the offset
	t := time.Now().AddDate(0, 0, offsetDays)

	// Format as YYYYMMDD
	// In Go, we use the "Reference Date" 2006-01-02
	return t.Format("20060102")
}

func PrintTeamSportsScores(scores []*Score, err error) {
	if err != nil {
		fmt.Println(Dim("!!"), "Error:", err)
		return
	}
	fmt.Println("\n" + BgIndigo + " SCOREBOARD " + ColorReset)
	fmt.Println()
	for _, s := range scores {
		status := Dim("○")
		eventTimeDisplay := ""
		// design choice for deciding the status badges
		eventDate, err := FormatToLocal(s.Date)
		if err != nil {
			fmt.Print("Error on date localization ", err)
		}
		eventDateDisplay := eventDate
		if s.Completed {
			status = ColorIndigo + "●" + ColorReset
			eventDateDisplay = ""
		}
		t, err := time.Parse("2006-01-02T15:04Z", s.Date)
		if err != nil {
			fmt.Println("Error on date parse ", err)
		}
		if err == nil && !s.Completed && t.Before(time.Now()) {
			status = ColorGreen + "●" + ColorReset
			eventDateDisplay = ""
			eventTimeDisplay = "- Period: " + strconv.Itoa(s.Period) + " | " + s.DisplayClock
		}

		fmt.Printf(" %s %s %s @ %s %s %s %s\n\n", status, s.HomeTeam, s.HomeScore, s.AwayScore, s.AwayTeam, Dim(eventDateDisplay), eventTimeDisplay)

		// Initialize Table
		if s.Details != nil {
			tb := table.NewWriter()
			tb.SetOutputMirror(os.Stdout)
			tb.AppendHeader(table.Row{"EVENT", "PLAYER", "TEAM", "TIME"})
			for _, highlights := range s.Details {
				text := ""
				if strings.Contains(highlights.Text, "Yellow") {
					text = ColorYellow + highlights.Text + ColorReset
				} else if strings.Contains(highlights.Text, "Red") {
					text = ColorRed + highlights.Text + ColorReset
				} else {
					text = ColorLavender + highlights.Text + ColorReset
				}

				if highlights.Player != "" {
					tb.AppendRow(table.Row{
						Dim(text),
						highlights.Player,
						Dim(highlights.Team),
						Dim(highlights.Minute),
					})
				}
			}

			// Render Table
			style := table.StyleLight
			style.Options.DrawBorder = false
			style.Options.SeparateColumns = true
			style.Options.SeparateHeader = true
			style.Options.SeparateRows = true
			style.Box.PaddingRight = "  "
			tb.SetStyle(style)

			fmt.Println()
			tb.Render()
		}
		fmt.Println("\n─────────────────────────────────────────────────────────────")
		fmt.Println()
	}
}

func PrintRacingTable(scores []Results, err error) {
	// Initialize Table
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"STATUS", "SESSION", "PODIUM", "CIRCUIT"})

	if err != nil {
		fmt.Println(Dim("!!"), "Error:", err)
		return
	}
	fmt.Println("\n" + BgIndigo + " SESSION RESULTS " + ColorReset)
	fmt.Println()
	for _, s := range scores {
		status := Dim("○")
		if s.SessionComplete {
			status = ColorIndigo + "●" + ColorReset
		}
		session := "Race"
		if s.SessionType != "" {
			session = s.SessionType
		}
		t.AppendRow(table.Row{
			status,
			session,
			s.Podium,
			Dim(s.Location),
		})
	}

	// Render Table
	style := table.StyleLight
	style.Options.DrawBorder = true
	style.Options.SeparateColumns = true
	style.Options.SeparateHeader = true
	style.Options.SeparateRows = true
	style.Box.PaddingRight = "  "
	t.SetStyle(style)

	t.SetColumnConfigs([]table.ColumnConfig{
		{Name: "PODIUM", WidthMax: 60},
	})

	fmt.Println()
	t.Render()
}

func PrintHeader() {
	fmt.Print("\033[H\033[2J") // Clear screen

	fmt.Printf("\n %s", MegaLogo())
}
