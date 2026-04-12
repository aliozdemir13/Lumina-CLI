package internal

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// for NFL, NBA, Champions League and Europa League response parse
type Score struct {
	URL          string
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

type Highlights struct {
	Text   string
	Player string
	Team   string
	Minute string
}

type EspnResponse struct {
	Events []Event `json:"events"`
}

type Event struct {
	ID           string        `json:"id"`
	Date         string        `json:"date"`
	ShortName    string        `json:"name"` // shortName || name
	Competitions []Competition `json:"competitions"`
	Status       Status        `json:"status"`
}

type Competition struct {
	Competitors []Competitor `json:"competitors"`
	Details     []Details    `json:"details"`
}

type Competitor struct {
	HomeAway string `json:"homeAway"`
	Score    string `json:"score"`
	Team     Team   `json:"team"`
	Winner   bool   `json:"winner"`
}

type Team struct {
	Id           string `json:"id"`
	DisplayName  string `json:"displayName"`
	Abbreviation string `json:"abbreviation"`
}

type Status struct {
	Type         Type   `json:"type"`
	DisplayClock string `json:"displayClock"`
	Period       int    `json:"period"`
}

type Type struct {
	Completed   bool   `json:"completed"`
	Text        string `json:"text"`
	Description string `json:"description"`
}

type Details struct {
	Type        Type               `json:"type"`
	AthleteName []AthletesInvolved `json:"athletesInvolved"`
	Clock       Clock              `json:"clock"`
}

type AthletesInvolved struct {
	DisplayName string `json:"displayName"`
	Team        Team   `json:"team"`
}

type Clock struct {
	Minute string `json:"displayValue"`
}

func (s *Score) FetchResults(sportsType string) ([]*Score, error) {
	resp, err := http.Get(s.URL + sportsType)
	if err != nil { // always check errors before closing the door to avoid panic results
		return nil, err
	}
	defer resp.Body.Close() // ALWAYS close the "door" when you're done

	var wrapper EspnResponse
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

		highlights, err := generateHighlights(comp.Details, homeTeam, awayTeam)
		if err != nil {
			fmt.Println("Error: ", err)
		}
		currentScore.Details = highlights
		cleanScores = append(cleanScores, currentScore)
	}

	return cleanScores, nil
}

func generateHighlights(details []Details, homeTeam Team, awayTeam Team) ([]Highlights, error) {
	teamNames := map[string]string{
		homeTeam.Id: homeTeam.DisplayName,
		awayTeam.Id: awayTeam.DisplayName,
	}
	var highlights []Highlights
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

	return highlights, nil
}
