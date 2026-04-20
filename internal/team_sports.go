// Package internal is managing the API logic, this class handles team sport logic
package internal

import (
	"encoding/json"
	"net/http"
	"time"
)

// Score for NFL, NBA, Champions League and Europa League response parse
type Score struct {
	HomeTeam     string
	AwayTeam     string
	HomeScore    string
	AwayScore    string
	DisplayClock string
	Period       int
	Completed    bool
	Date         string
	Details      []Highlights
	Status       string
}

// Highlights struct is helper for Score to display match higlights
type Highlights struct {
	Text   string
	Player string
	Team   string
	Minute string
}

// ESPNResponse in use for parsing API response for team sports
type ESPNResponse struct {
	Events []Event `json:"events"`
}

// Event parses the match details
type Event struct {
	ID           string        `json:"id"`
	Date         string        `json:"date"`
	ShortName    string        `json:"name"` // shortName || name
	Competitions []Competition `json:"competitions"`
	Status       Status        `json:"status"`
}

// Competition provides details of the match and competitors
type Competition struct {
	Competitors []Competitor `json:"competitors"`
	Details     []Details    `json:"details"`
}

// Competitor represent teams playing, score and status
type Competitor struct {
	HomeAway string `json:"homeAway"`
	Score    string `json:"score"`
	Team     Team   `json:"team"`
	Winner   bool   `json:"winner"`
}

// Team struct provides the details of the each team
type Team struct {
	Id           string `json:"id"`
	DisplayName  string `json:"displayName"`
	Abbreviation string `json:"abbreviation"`
}

// Status struct is a for game status and time
type Status struct {
	Type         Type   `json:"type"`
	DisplayClock string `json:"displayClock"`
	Period       int    `json:"period"`
}

// Type struct parses the game status
type Type struct {
	Completed   bool   `json:"completed"`
	Text        string `json:"text"`
	Description string `json:"description"`
}

// Details struct parses the higlight of the game
type Details struct {
	Type        Type               `json:"type"`
	AthleteName []AthletesInvolved `json:"athletesInvolved"`
	Clock       Clock              `json:"clock"`
}

// AthletesInvolved diplays the actor of the highlight
type AthletesInvolved struct {
	DisplayName string `json:"displayName"`
	Team        Team   `json:"team"`
}

// Clock struct parses the time of the event
type Clock struct {
	Minute string `json:"displayValue"`
}

// FetchResults pulls the API result for the team sport events
func (s *Score) FetchResults(endpoint string) ([]*Score, error) {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	resp, err := client.Get(endpoint)
	if err != nil { // always check errors before closing the door to avoid panic results
		return nil, err
	}
	defer func() {
		_ = resp.Body.Close()
	}() // ALWAYS close the "door" when you're done

	var wrapper ESPNResponse
	if err := json.NewDecoder(resp.Body).Decode(&wrapper); err != nil {
		return nil, err
	}

	var cleanScores []*Score

	for _, event := range wrapper.Events {
		// ESPN usually has 1 competition per event
		comp := event.Competitions[0]

		currentScore := &Score{
			Status:       event.Status.Type.Description,
			Date:         event.Date,
			Completed:    event.Status.Type.Completed,
			Period:       event.Status.Period,
			DisplayClock: event.Status.DisplayClock,
		}

		var homeTeam Team
		var awayTeam Team
		// Loop through competitors to find Home vs Away
		for _, teamData := range comp.Competitors {
			if teamData.HomeAway == "home" {
				currentScore.HomeTeam = teamData.Team.DisplayName // can switch abbreviation as well if desired
				currentScore.HomeScore = teamData.Score           // Note: ESPN scores are often strings!
				homeTeam = teamData.Team
			} else {
				currentScore.AwayTeam = teamData.Team.DisplayName
				currentScore.AwayScore = teamData.Score
				awayTeam = teamData.Team
			}
		}

		highlights := generateHighlights(comp.Details, homeTeam, awayTeam)
		currentScore.Details = highlights
		cleanScores = append(cleanScores, currentScore)
	}

	return cleanScores, nil
}

func generateHighlights(details []Details, homeTeam Team, awayTeam Team) []Highlights {
	teamNames := map[string]string{
		homeTeam.Id: homeTeam.DisplayName,
		awayTeam.Id: awayTeam.DisplayName,
	}
	var highlights []Highlights
	if len(details) == 0 {
		return nil
	}
	for _, detail := range details {
		playerName := ""
		teamName := ""
		if len(detail.AthleteName) > 0 {
			playerName = detail.AthleteName[0].DisplayName
			teamName = teamNames[detail.AthleteName[0].Team.Id]
		}
		highlight := Highlights{
			Text:   detail.Type.Text,
			Player: playerName,
			Team:   teamName,
			Minute: detail.Clock.Minute,
		}

		highlights = append(highlights, highlight)
	}

	return highlights
}
