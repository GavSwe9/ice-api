package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"github.com/gavswe19/ice-api/database"
)

type Response events.APIGatewayProxyResponse

type SeasonTeam struct {
	// season  |team_id|name                 |abbreviation|division_id|division_name|conference_id|conference_name|franchise_id
	Season         int    `json:"season"`
	TeamId         int    `json:"teamId"`
	Name           string `json:"name"`
	Abbreviation   string `json:"abbreviation"`
	DivisionId     int    `json:"divisionId"`
	DivisionName   string `json:"divisionName"`
	ConferenceId   int    `json:"conferenceId"`
	ConferenceName string `json:"conferenceName"`
	FranchiseId    int    `json:"franchiseId"`
}

func Handler(request events.APIGatewayProxyRequest) (Response, error) {
	db := database.GetDatabase()

	season, err := strconv.Atoi(request.PathParameters["season"])

	if err != nil {
		log.Fatal("Season must be an integer")
	}

	query := query(season)
	results, err := db.Query(query)

	if err != nil {
		log.Fatal("Error querying line stats")
	}

	var seasonTeams []SeasonTeam

	for results.Next() {
		var seasonTeam SeasonTeam
		results.Scan(
			&seasonTeam.Season,
			&seasonTeam.TeamId,
			&seasonTeam.Name,
			&seasonTeam.Abbreviation,
			&seasonTeam.DivisionId,
			&seasonTeam.DivisionName,
			&seasonTeam.ConferenceId,
			&seasonTeam.ConferenceName,
			&seasonTeam.FranchiseId,
		)
		seasonTeams = append(seasonTeams, seasonTeam)
	}

	jsonData, err := json.Marshal(seasonTeams)
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

func query(season int) string {
	return fmt.Sprintf(`
	SELECT
	*
	FROM team_seasons ts 
	WHERE ts.season = %s
	`,
		strconv.Itoa(season))
}

func main() {
	lambda.Start(Handler)
}
