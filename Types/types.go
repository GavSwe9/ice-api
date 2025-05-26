package types

import (
	"time"

	"github.com/aws/aws-lambda-go/events"
)

type Response events.APIGatewayProxyResponse

type SearchResult struct {
	ItemsFound bool         `json:"itemsFound"`
	Teams      []TeamSeason `json:"teams"`
	Players    []PlayerBio  `json:"players"`
}

type PlayerBio struct {
	PlayerId            int       `json:"playerId"`
	IsActive            bool      `json:"isActive"`
	CurrentTeamId       int       `json:"currentTeamId"`
	CurrentTeamAbbrev   string    `json:"currentTeamAbbrev"`
	FullTeamName        string    `json:"fullTeamName"`
	FirstName           string    `json:"firstName"`
	LastName            string    `json:"lastName"`
	FullName            string    `json:"fullName"`
	SweaterNumber       int       `json:"sweaterNumber"`
	Position            string    `json:"position"`
	Headshot            string    `json:"headshot"`
	HeroImage           string    `json:"heroImage"`
	HeightInInches      int       `json:"heightInInches"`
	HeightInCentimeters int       `json:"heightInCentimeters"`
	WeightInPounds      int       `json:"weightInPounds"`
	WeightInKilograms   int       `json:"weightInKilograms"`
	BirthDate           time.Time `json:"birthDate"`
	BirthCity           string    `json:"birthCity"`
	BirthStateProvince  string    `json:"birthStateProvince"`
	BirthCountry        string    `json:"birthCountry"`
	ShootsCatches       string    `json:"shootsCatches"`
	DraftYear           int       `json:"draftYear"`
	DraftTeamAbbrev     string    `json:"draftTeamAbbrev"`
	DraftRound          int       `json:"draftRound"`
	DraftPickInRound    int       `json:"draftPickInRound"`
	DraftOverallPick    int       `json:"draftOverallPick"`
}

type TeamSeason struct {
	ID             int    `json:"id"`
	SeasonID       int    `json:"seasonId"`
	ConferenceName string `json:"conferenceName"`
	DivisionName   string `json:"divisionName"`
	TeamName       string `json:"teamName"`
	TeamCommonName string `json:"teamCommonName"`
	TeamAbbrev     string `json:"teamAbbrev"`
}
