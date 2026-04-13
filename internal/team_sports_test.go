package internal

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// --- generateHighlights tests (pure function, no network) ---

func TestGenerateHighlights(t *testing.T) {
	home := Team{Id: "1", DisplayName: "Arsenal"}
	away := Team{Id: "2", DisplayName: "Chelsea"}

	t.Run("single goal event", func(t *testing.T) {
		details := []Details{
			{
				Type:        Type{Text: "Goal - Normal"},
				AthleteName: []AthletesInvolved{{DisplayName: "Saka", Team: home}},
				Clock:       Clock{Minute: "23'"},
			},
		}

		highlights := generateHighlights(details, home, away)
		if len(highlights) != 1 {
			t.Fatalf("got %d highlights, want 1", len(highlights))
		}
		if highlights[0].Player != "Saka" {
			t.Errorf("player = %q, want %q", highlights[0].Player, "Saka")
		}
		if highlights[0].Team != "Arsenal" {
			t.Errorf("team = %q, want %q", highlights[0].Team, "Arsenal")
		}
		if highlights[0].Text != "Goal - Normal" {
			t.Errorf("text = %q, want %q", highlights[0].Text, "Goal - Normal")
		}
		if highlights[0].Minute != "23'" {
			t.Errorf("minute = %q, want %q", highlights[0].Minute, "23'")
		}
	})

	t.Run("away team goal resolves correctly", func(t *testing.T) {
		details := []Details{
			{
				Type:        Type{Text: "Goal - Normal"},
				AthleteName: []AthletesInvolved{{DisplayName: "Palmer", Team: away}},
				Clock:       Clock{Minute: "55'"},
			},
		}

		highlights := generateHighlights(details, home, away)
		if highlights[0].Team != "Chelsea" {
			t.Errorf("team = %q, want %q", highlights[0].Team, "Chelsea")
		}
	})

	t.Run("multiple events preserve order", func(t *testing.T) {
		details := []Details{
			{
				Type:        Type{Text: "Yellow Card"},
				AthleteName: []AthletesInvolved{{DisplayName: "Rice", Team: home}},
				Clock:       Clock{Minute: "30'"},
			},
			{
				Type:        Type{Text: "Red Card"},
				AthleteName: []AthletesInvolved{{DisplayName: "Caicedo", Team: away}},
				Clock:       Clock{Minute: "44'"},
			},
		}

		highlights := generateHighlights(details, home, away)
		if len(highlights) != 2 {
			t.Fatalf("got %d highlights, want 2", len(highlights))
		}
		if highlights[0].Text != "Yellow Card" {
			t.Errorf("first event = %q, want %q", highlights[0].Text, "Yellow Card")
		}
		if highlights[1].Text != "Red Card" {
			t.Errorf("second event = %q, want %q", highlights[1].Text, "Red Card")
		}
	})

	t.Run("empty details returns nil", func(t *testing.T) {
		highlights := generateHighlights([]Details{}, home, away)
		if highlights != nil {
			t.Errorf("got %v, want nil", highlights)
		}
	})

	t.Run("detail with no athlete has empty player and team", func(t *testing.T) {
		details := []Details{
			{
				Type:        Type{Text: "End Period"},
				AthleteName: []AthletesInvolved{},
				Clock:       Clock{Minute: "45'"},
			},
		}

		highlights := generateHighlights(details, home, away)
		if highlights[0].Player != "" {
			t.Errorf("player = %q, want empty string", highlights[0].Player)
		}
		if highlights[0].Team != "" {
			t.Errorf("team = %q, want empty string", highlights[0].Team)
		}
	})
}

// --- FetchResults tests (httptest fake server) ---

const completedGameJSON = `{
	"events": [{
		"id": "401",
		"date": "2025-06-15T18:30Z",
		"name": "LAL vs BOS",
		"competitions": [{
			"competitors": [
				{
					"homeAway": "home",
					"score": "110",
					"team": {"id": "13", "displayName": "Lakers", "abbreviation": "LAL"},
					"winner": true
				},
				{
					"homeAway": "away",
					"score": "105",
					"team": {"id": "2", "displayName": "Celtics", "abbreviation": "BOS"},
					"winner": false
				}
			],
			"details": []
		}],
		"status": {
			"type": {"completed": true, "text": "Final", "description": "Final"},
			"displayClock": "0:00",
			"period": 4
		}
	}]
}`

const multiGameJSON = `{
	"events": [
		{
			"id": "401",
			"date": "2025-06-15T18:30Z",
			"name": "LAL vs BOS",
			"competitions": [{
				"competitors": [
					{"homeAway": "home", "score": "110", "team": {"id": "1", "displayName": "Lakers", "abbreviation": "LAL"}, "winner": true},
					{"homeAway": "away", "score": "105", "team": {"id": "2", "displayName": "Celtics", "abbreviation": "BOS"}, "winner": false}
				],
				"details": []
			}],
			"status": {"type": {"completed": true, "text": "Final", "description": "Final"}, "displayClock": "0:00", "period": 4}
		},
		{
			"id": "402",
			"date": "2025-06-15T20:00Z",
			"name": "GSW vs MIA",
			"competitions": [{
				"competitors": [
					{"homeAway": "home", "score": "98", "team": {"id": "3", "displayName": "Warriors", "abbreviation": "GSW"}, "winner": false},
					{"homeAway": "away", "score": "102", "team": {"id": "4", "displayName": "Heat", "abbreviation": "MIA"}, "winner": true}
				],
				"details": []
			}],
			"status": {"type": {"completed": true, "text": "Final", "description": "Final"}, "displayClock": "0:00", "period": 4}
		}
	]
}`

const gameWithHighlightsJSON = `{
	"events": [{
		"id": "501",
		"date": "2025-06-15T20:00Z",
		"name": "Arsenal vs Chelsea",
		"competitions": [{
			"competitors": [
				{"homeAway": "home", "score": "2", "team": {"id": "1", "displayName": "Arsenal", "abbreviation": "ARS"}, "winner": true},
				{"homeAway": "away", "score": "1", "team": {"id": "2", "displayName": "Chelsea", "abbreviation": "CHE"}, "winner": false}
			],
			"details": [
				{
					"type": {"text": "Goal - Normal", "completed": false, "description": ""},
					"athletesInvolved": [{"displayName": "Saka", "team": {"id": "1", "displayName": "Arsenal", "abbreviation": "ARS"}}],
					"clock": {"displayValue": "23'"}
				},
				{
					"type": {"text": "Yellow Card", "completed": false, "description": ""},
					"athletesInvolved": [{"displayName": "Caicedo", "team": {"id": "2", "displayName": "Chelsea", "abbreviation": "CHE"}}],
					"clock": {"displayValue": "40'"}
				}
			]
		}],
		"status": {"type": {"completed": true, "text": "Final", "description": "Final"}, "displayClock": "0:00", "period": 2}
	}]
}`

func newFakeServer(responseBody string) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(responseBody))
	}))
}

func TestFetchResults_CompletedGame(t *testing.T) {
	server := newFakeServer(completedGameJSON)
	defer server.Close()

	var svc Score
	scores, err := svc.FetchResults(server.URL + "/basketball/nba/scoreboard")

	if err != nil {
		t.Fatalf("FetchResults returned error: %v", err)
	}
	if len(scores) != 1 {
		t.Fatalf("got %d scores, want 1", len(scores))
	}

	s := scores[0]
	if s.HomeTeam != "Lakers" {
		t.Errorf("HomeTeam = %q, want %q", s.HomeTeam, "Lakers")
	}
	if s.AwayTeam != "Celtics" {
		t.Errorf("AwayTeam = %q, want %q", s.AwayTeam, "Celtics")
	}
	if s.HomeScore != "110" {
		t.Errorf("HomeScore = %q, want %q", s.HomeScore, "110")
	}
	if s.AwayScore != "105" {
		t.Errorf("AwayScore = %q, want %q", s.AwayScore, "105")
	}
	if !s.Completed {
		t.Error("expected Completed to be true")
	}
	if s.Period != 4 {
		t.Errorf("Period = %d, want 4", s.Period)
	}
	if s.Status != "Final" {
		t.Errorf("Status = %q, want %q", s.Status, "Final")
	}
}

func TestFetchResults_MultipleGames(t *testing.T) {
	server := newFakeServer(multiGameJSON)
	defer server.Close()

	var svc Score
	scores, err := svc.FetchResults(server.URL + "/basketball/nba/scoreboard")

	if err != nil {
		t.Fatalf("FetchResults returned error: %v", err)
	}
	if len(scores) != 2 {
		t.Fatalf("got %d scores, want 2", len(scores))
	}
	if scores[0].HomeTeam != "Lakers" {
		t.Errorf("game 1 HomeTeam = %q, want %q", scores[0].HomeTeam, "Lakers")
	}
	if scores[1].HomeTeam != "Warriors" {
		t.Errorf("game 2 HomeTeam = %q, want %q", scores[1].HomeTeam, "Warriors")
	}
}

func TestFetchResults_WithHighlights(t *testing.T) {
	server := newFakeServer(gameWithHighlightsJSON)
	defer server.Close()

	var svc Score
	scores, err := svc.FetchResults(server.URL + "/soccer/eng.1/scoreboard")

	if err != nil {
		t.Fatalf("FetchResults returned error: %v", err)
	}
	if len(scores) != 1 {
		t.Fatalf("got %d scores, want 1", len(scores))
	}

	details := scores[0].Details
	if len(details) != 2 {
		t.Fatalf("got %d details, want 2", len(details))
	}
	if details[0].Player != "Saka" {
		t.Errorf("first detail player = %q, want %q", details[0].Player, "Saka")
	}
	if details[1].Text != "Yellow Card" {
		t.Errorf("second detail text = %q, want %q", details[1].Text, "Yellow Card")
	}
}

func TestFetchResults_EmptyEvents(t *testing.T) {
	server := newFakeServer(`{"events": []}`)
	defer server.Close()

	var svc Score
	scores, err := svc.FetchResults(server.URL + "/basketball/nba/scoreboard")

	if err != nil {
		t.Fatalf("FetchResults returned error: %v", err)
	}
	if scores != nil {
		t.Errorf("expected nil scores for empty events, got %d", len(scores))
	}
}

func TestFetchResults_MalformedJSON(t *testing.T) {
	server := newFakeServer(`{not valid json}`)
	defer server.Close()

	var svc Score
	_, err := svc.FetchResults(server.URL + "/basketball/nba/scoreboard")

	if err == nil {
		t.Error("expected error for malformed JSON, got nil")
	}
}

func TestFetchResults_ServerError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"events": []}`))
	}))
	defer server.Close()

	var svc Score
	scores, err := svc.FetchResults(server.URL + "/basketball/nba/scoreboard")

	// Current implementation doesn't check status codes, so it will still
	// parse the body. This test documents that behavior — you may want to
	// add status code checking in the future.
	if err != nil {
		t.Fatalf("FetchResults returned error: %v", err)
	}
	if scores != nil {
		t.Errorf("expected nil scores, got %d", len(scores))
	}
}

func TestFetchResults_UnreachableServer(t *testing.T) {
	var svc Score // port 0 = nothing listening
	_, err := svc.FetchResults("http://127.0.0.1:0/basketball/nba/scoreboard")

	if err == nil {
		t.Error("expected error for unreachable server, got nil")
	}
}
