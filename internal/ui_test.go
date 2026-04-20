package internal

import (
	"fmt"
	"regexp"
	"testing"
	"time"
)

// --- generateHighlights tests (pure function, no network) ---

func TestGetEspnDateAndPrintHeader(t *testing.T) {
	datePattern := regexp.MustCompile(`^\d{8}$`)

	t.Run("zero offset returns today", func(t *testing.T) {
		got := GetEspnDate(0)
		want := time.Now().Format("20060102")

		if got != want {
			t.Errorf("GetEspnDate(0) = %q, want %q", got, want)
		}
	})

	t.Run("negative offset returns past date", func(t *testing.T) {
		got := GetEspnDate(-1)
		want := time.Now().AddDate(0, 0, -1).Format("20060102")

		if got != want {
			t.Errorf("GetEspnDate(-1) = %q, want %q", got, want)
		}
	})

	t.Run("positive offset returns future date", func(t *testing.T) {
		got := GetEspnDate(7)
		want := time.Now().AddDate(0, 0, 7).Format("20060102")

		if got != want {
			t.Errorf("GetEspnDate(7) = %q, want %q", got, want)
		}
	})

	t.Run("output format is always YYYYMMDD", func(t *testing.T) {
		got := GetEspnDate(-30)

		if !datePattern.MatchString(got) {
			t.Errorf("GetEspnDate(-30) = %q, does not match YYYYMMDD pattern", got)
		}
	})

	t.Run("print header", func(t *testing.T) {
		PrintHeader()
	})
}

func TestPrintTeamSportsScores(t *testing.T) {
	const liveGameJSON = `{
		"events": [{
			"id": "601",
			"date": "2023-01-01T10:00Z", 
			"name": "Live Game",
			"competitions": [{
				"competitors": [
					{"homeAway": "home", "score": "1", "team": {"displayName": "Team A"}},
					{"homeAway": "away", "score": "1", "team": {"displayName": "Team B"}}
				],
				"details": [
					{"type": {"text": "Red Card"}, "athletesInvolved": [{"displayName": "Player X"}], "clock": {"displayValue": "75'"}},
					{"type": {"text": "No Player Event"}, "athletesInvolved": [], "clock": {"displayValue": "80'"}}
				]
			}],
			"status": {"type": {"completed": false}, "displayClock": "75:00", "period": 2}
		}]
	}`

	const invalidDateGameJSON = `{
		"events": [{
			"date": "this-is-not-a-date",
			"competitions": [{"competitors": []}]
		}]
	}`

	t.Run("error case", func(t *testing.T) {
		err := fmt.Errorf("Something went wrong")
		PrintTeamSportsScores(nil, err)

	})

	t.Run("print results", func(t *testing.T) {
		server := newFakeServer(completedGameJSON)
		defer server.Close()

		var svc Score
		scores, _ := svc.FetchResults(server.URL + "/basketball/nba/scoreboard")

		PrintTeamSportsScores(scores, nil)
	})

	t.Run("print results with highlight", func(t *testing.T) {

		server := newFakeServer(gameWithHighlightsJSON)
		defer server.Close()

		var svc Score
		scores, _ := svc.FetchResults(server.URL + "/soccer/eng.1/scoreboard")

		PrintTeamSportsScores(scores, nil)
	})

	t.Run("print results with multiple games", func(t *testing.T) {

		server := newFakeServer(multiGameJSON)
		defer server.Close()

		var svc Score
		scores, _ := svc.FetchResults(server.URL + "/basketball/nba/scoreboard")

		PrintTeamSportsScores(scores, nil)
	})

	t.Run("print results with finished games", func(t *testing.T) {

		server := newFakeServer(completedGameJSON)
		defer server.Close()

		var svc Score
		scores, _ := svc.FetchResults(server.URL + "/basketball/nba/scoreboard")

		PrintTeamSportsScores(scores, nil)
	})

	t.Run("live game with red card and empty player", func(t *testing.T) {
		server := newFakeServer(liveGameJSON)
		defer server.Close()
		var svc Score
		scores, _ := svc.FetchResults(server.URL + "/any")
		PrintTeamSportsScores(scores, nil)
	})

	t.Run("invalid date errors", func(t *testing.T) {
		server := newFakeServer(invalidDateGameJSON)
		defer server.Close()
		var svc Score
		scores, _ := svc.FetchResults(server.URL + "/any")
		PrintTeamSportsScores(scores, nil)
	})
}

func TestPrintRacingTable(t *testing.T) {
	t.Run("error case", func(t *testing.T) {
		err := fmt.Errorf("Something went wrong")
		PrintRacingTable(nil, err)

	})

	t.Run("print results", func(t *testing.T) {
		server := newRacingFakeServer(multiEventJSON)
		defer server.Close()

		var svc Results
		results, _ := svc.FetchResults(server.URL + "/racing/f1/scoreboard")

		PrintRacingTable(results, nil)
	})
}
