package internal

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

const raceWeekendJSON = `{
	"events": [{
		"shortName": "Australian GP",
		"competitions": [
			{
				"date": "2025-03-14T01:30Z",
				"status": {"type": {"completed": true, "text": "Final", "description": "Final"}, "displayClock": "", "period": 0},
				"type": {"abbreviation": "FP1"},
				"competitors": [
					{"order": 1, "athlete": {"shortName": "VER"}},
					{"order": 2, "athlete": {"shortName": "NOR"}},
					{"order": 3, "athlete": {"shortName": "LEC"}},
					{"order": 4, "athlete": {"shortName": "HAM"}}
				]
			},
			{
				"date": "2025-03-16T05:00Z",
				"status": {"type": {"completed": true, "text": "Final", "description": "Final"}, "displayClock": "", "period": 0},
				"type": {"abbreviation": "RACE"},
				"competitors": [
					{"order": 1, "athlete": {"shortName": "LEC"}},
					{"order": 2, "athlete": {"shortName": "HAM"}},
					{"order": 3, "athlete": {"shortName": "VER"}},
					{"order": 4, "athlete": {"shortName": "NOR"}}
				]
			}
		]
	}]
}`

const upcomingRaceJSON = `{
	"events": [{
		"shortName": "Monaco GP",
		"competitions": [{
			"date": "2025-05-25T13:00Z",
			"status": {"type": {"completed": false, "text": "Scheduled", "description": "Scheduled"}, "displayClock": "", "period": 0},
			"type": {"abbreviation": "RACE"},
			"competitors": []
		}]
	}]
}`

const multiEventJSON = `{
	"events": [
		{
			"shortName": "Australian GP",
			"competitions": [{
				"date": "2025-03-16T05:00Z",
				"status": {"type": {"completed": true, "text": "Final", "description": "Final"}, "displayClock": "", "period": 0},
				"type": {"abbreviation": "RACE"},
				"competitors": [
					{"order": 1, "athlete": {"shortName": "LEC"}},
					{"order": 2, "athlete": {"shortName": "HAM"}},
					{"order": 3, "athlete": {"shortName": "VER"}}
				]
			}]
		},
		{
			"shortName": "Japanese GP",
			"competitions": [{
				"date": "2025-04-06T05:00Z",
				"status": {"type": {"completed": true, "text": "Final", "description": "Final"}, "displayClock": "", "period": 0},
				"type": {"abbreviation": "RACE"},
				"competitors": [
					{"order": 1, "athlete": {"shortName": "VER"}},
					{"order": 2, "athlete": {"shortName": "NOR"}},
					{"order": 3, "athlete": {"shortName": "PIA"}}
				]
			}]
		}
	]
}`

func newRacingFakeServer(body string) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(body))
	}))
}

func TestRacingFetchResults_FullWeekend(t *testing.T) {
	server := newRacingFakeServer(raceWeekendJSON)
	defer server.Close()

	var svc Results
	results, err := svc.FetchResults(server.URL + "/racing/f1/scoreboard")

	if err != nil {
		t.Fatalf("FetchResults returned error: %v", err)
	}
	if len(results) != 2 {
		t.Fatalf("got %d results, want 2 (FP1 + Race)", len(results))
	}

	// FP1 session
	fp1 := results[0]
	if fp1.SessionType != "FP1" {
		t.Errorf("session 1 type = %q, want %q", fp1.SessionType, "FP1")
	}
	if fp1.Location != "Australian GP" {
		t.Errorf("location = %q, want %q", fp1.Location, "Australian GP")
	}
	if !fp1.SessionComplete {
		t.Error("expected FP1 to be complete")
	}
	if fp1.Podium != "VER - NOR - LEC" {
		t.Errorf("FP1 podium = %q, want %q", fp1.Podium, "VER - NOR - LEC")
	}

	// Race session
	race := results[1]
	if race.SessionType != "RACE" {
		t.Errorf("session 2 type = %q, want %q", race.SessionType, "RACE")
	}
	if race.Podium != "LEC - HAM - VER" {
		t.Errorf("race podium = %q, want %q", race.Podium, "LEC - HAM - VER")
	}
}

func TestRacingFetchResults_PodiumExcludesNonTop3(t *testing.T) {
	server := newRacingFakeServer(raceWeekendJSON)
	defer server.Close()

	var svc Results
	results, err := svc.FetchResults(server.URL + "/racing/f1/scoreboard")

	if err != nil {
		t.Fatalf("FetchResults returned error: %v", err)
	}

	// FP1 has 4 drivers but podium should only show top 3
	fp1 := results[0]
	if fp1.Podium != "VER - NOR - LEC" {
		t.Errorf("podium includes P4 driver: %q", fp1.Podium)
	}
}

func TestRacingFetchResults_NoCompetitors(t *testing.T) {
	server := newRacingFakeServer(upcomingRaceJSON)
	defer server.Close()

	var svc Results
	results, err := svc.FetchResults(server.URL + "/racing/f1/scoreboard")

	if err != nil {
		t.Fatalf("FetchResults returned error: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("got %d results, want 1", len(results))
	}
	if results[0].Podium != "TBD" {
		t.Errorf("podium = %q, want %q", results[0].Podium, "TBD")
	}
	if results[0].SessionComplete {
		t.Error("expected session to be incomplete")
	}
}

func TestRacingFetchResults_MultipleEvents(t *testing.T) {
	server := newRacingFakeServer(multiEventJSON)
	defer server.Close()

	var svc Results
	results, err := svc.FetchResults(server.URL + "/racing/f1/scoreboard")

	if err != nil {
		t.Fatalf("FetchResults returned error: %v", err)
	}
	if len(results) != 2 {
		t.Fatalf("got %d results, want 2", len(results))
	}
	if results[0].Location != "Australian GP" {
		t.Errorf("event 1 location = %q, want %q", results[0].Location, "Australian GP")
	}
	if results[1].Location != "Japanese GP" {
		t.Errorf("event 2 location = %q, want %q", results[1].Location, "Japanese GP")
	}
	if results[1].Podium != "VER - NOR - PIA" {
		t.Errorf("event 2 podium = %q, want %q", results[1].Podium, "VER - NOR - PIA")
	}
}

func TestRacingFetchResults_EmptyEvents(t *testing.T) {
	server := newRacingFakeServer(`{"events": []}`)
	defer server.Close()

	var svc Results
	results, err := svc.FetchResults(server.URL + "/racing/f1/scoreboard")

	if err != nil {
		t.Fatalf("FetchResults returned error: %v", err)
	}
	if results != nil {
		t.Errorf("expected nil results for empty events, got %d", len(results))
	}
}

func TestRacingFetchResults_MalformedJSON(t *testing.T) {
	server := newRacingFakeServer(`{broken`)
	defer server.Close()

	var svc Results
	_, err := svc.FetchResults(server.URL + "/racing/f1/scoreboard")

	if err == nil {
		t.Error("expected error for malformed JSON, got nil")
	}
}

func TestRacingFetchResults_UnreachableServer(t *testing.T) {
	var svc Results
	_, err := svc.FetchResults("http://127.0.0.1:0/racing/f1/scoreboard")

	if err == nil {
		t.Error("expected error for unreachable server, got nil")
	}
}
