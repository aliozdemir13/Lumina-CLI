// Package internal is managing the API logic, this class handles racing results
package internal

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

// Results parses the races results to display in UI
type Results struct {
	Podium          string // sample: LEC - HAM - VER
	Location        string // sample: Qatar Airways Australian GP
	SessionType     string // FP1 | SQ | SR | Quali | Race
	SessionDate     string
	SessionComplete bool
}

// RacingResponse in use for the API response parsing
type RacingResponse struct {
	Events []RacingEvent `json:"events"`
}

// RacingEvent parses the individual race weeks
type RacingEvent struct {
	ShortName    string              `json:"shortName"` // "Australian GP"
	Competitions []RacingCompetition `json:"competitions"`
}

// RacingCompetition provides details about the races
type RacingCompetition struct {
	Date        string             `json:"date"`
	Status      Status             `json:"status"`
	Type        SessionType        `json:"type"` // "Race", "Qualifying", etc.
	Competitors []RacingCompetitor `json:"competitors"`
}

// SessionType parses type of session
type SessionType struct {
	Abbreviation string `json:"abbreviation"` // "RACE", "PROV", "FP1"
}

// RacingCompetitor parses the racers
type RacingCompetitor struct {
	Order   int `json:"order"`
	Athlete struct {
		Abbreviation string `json:"shortName"` // "LEC", "HAM"
	} `json:"athlete"`
}

// FetchResults pulls the API response for motorsports
func (s *Results) FetchResults(endpoint string) ([]Results, error) {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	resp, err := client.Get(endpoint)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	var wrapper RacingResponse
	if err := json.NewDecoder(resp.Body).Decode(&wrapper); err != nil {
		return nil, err
	}

	var cleanResults []Results

	for _, event := range wrapper.Events {
		// Each competition is a session (Race, Quali, etc.)
		for _, session := range event.Competitions {
			res := Results{
				Location:        event.ShortName,
				SessionType:     session.Type.Abbreviation,
				SessionComplete: session.Status.Type.Completed,
			}

			res.SessionDate, err = FormatToLocal(session.Date)
			if err != nil {
				fmt.Printf(`Error: %s`, err)
			}

			// Build the "LEC - HAM - VER" string from competitors
			var podiumNames []string
			for _, driver := range session.Competitors {
				if driver.Order <= 3 && driver.Order > 0 {
					podiumNames = append(podiumNames, driver.Athlete.Abbreviation)
				}
			}

			// Join the slice into a single string
			if len(podiumNames) > 0 {
				res.Podium = strings.Join(podiumNames, " - ")
			} else {
				res.Podium = "TBD"
			}

			cleanResults = append(cleanResults, res)
		}
	}
	return cleanResults, nil
}
