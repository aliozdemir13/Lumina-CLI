package internal

import (
	"encoding/json"
	"net/http"
	"strings"
)

// for races
type Results struct {
	URL             string
	Podium          string // sample: LEC - HAM - VER
	Location        string // sample: Qatar Airways Australian GP
	SessionType     string // FP1 | SQ | SR | Quali | Race
	SessionDate     string
	SessionComplete bool
}

type RacingResponse struct {
	Events []RacingEvent `json:"events"`
}

type RacingEvent struct {
	ShortName    string              `json:"shortName"` // "Australian GP"
	Competitions []RacingCompetition `json:"competitions"`
}

type RacingCompetition struct {
	Date        string             `json:"date"`
	Status      Status             `json:"status"`
	Type        SessionType        `json:"type"` // "Race", "Qualifying", etc.
	Competitors []RacingCompetitor `json:"competitors"`
}

type SessionType struct {
	Abbreviation string `json:"abbreviation"` // "RACE", "PROV", "FP1"
}

type RacingCompetitor struct {
	Order   int `json:"order"`
	Athlete struct {
		Abbreviation string `json:"shortName"` // "LEC", "HAM"
	} `json:"athlete"`
}

func (s *Results) FetchResults(path string) ([]Results, error) {
	resp, err := http.Get(s.URL + path)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

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

			res.SessionDate, _ = FormatToLocal(session.Date)

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
