package main

import (
	"encoding/json"
	"log"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	types "github.com/gavswe19/ice-api/Types"
	"github.com/gavswe19/ice-api/database"
)

type Response events.APIGatewayProxyResponse

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

func Handler(request events.APIGatewayProxyRequest) (Response, error) {
	db := database.GetDatabase()

	teamAbv := request.PathParameters["teamAbv"]

	query := `
		SELECT 
		*
		FROM ICE.team_seasons_2 ts 
		WHERE 
		ts.team_abbrev = ?
		ORDER BY season_id DESC 
		LIMIT 1
	`

	results, err := db.Query(query, teamAbv)

	if err != nil {
		log.Fatal(err)
	}
	defer results.Close()

	if !results.Next() {
		resp := Response{
			StatusCode:      404,
			IsBase64Encoded: false,
			Headers: map[string]string{
				"Content-Type":                     "application/json",
				"X-MyCompany-Func-Reply":           "hello-handler",
				"Access-Control-Allow-Origin":      "*",
				"Access-Control-Allow-Credentials": "true",
			},
		}

		return resp, nil
	}

	var teamSeason types.TeamSeason

	err = results.Scan(
		&teamSeason.ID,
		&teamSeason.TeamAbbrev,
		&teamSeason.SeasonID,
		&teamSeason.ConferenceName,
		&teamSeason.DivisionName,
		&teamSeason.TeamName,
		&teamSeason.TeamCommonName,
	)
	if err != nil {
		log.Fatalf("Error scanning row: %v", err)
	}

	jsonData, err := json.Marshal(teamSeason)
	if err != nil {
		log.Fatal(err)
	}

	resp := Response{
		StatusCode:      200,
		IsBase64Encoded: false,
		Body:            string(jsonData),
		Headers: map[string]string{
			"Content-Type":                     "application/json",
			"X-MyCompany-Func-Reply":           "hello-handler",
			"Access-Control-Allow-Origin":      "*",
			"Access-Control-Allow-Credentials": "true",
		},
	}

	return resp, nil
}

func main() {
	lambda.Start(Handler)
}
