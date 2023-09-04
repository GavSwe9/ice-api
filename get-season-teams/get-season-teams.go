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

type RequestBody struct {
	Season int `json:"season"`
}

type SeasonTeam struct {
	// season  |team_id|name                 |abbreviation|division_id|division_name|conference_id|conference_name|franchise_id
	Season         int    `json:"season"`
	TeamId         int    `json:"team_id"`
	Name           string `json:"name"`
	Abbreviation   string `json:"abbreviation"`
	DivisionId     int    `json:"division_id"`
	DivisionName   string `json:"division_name"`
	ConferenceId   int    `json:"conference_id"`
	ConferenceName string `json:"conference_name"`
	FranchiseId    int    `json:"franchise_id"`
}

func Handler(request events.APIGatewayProxyRequest) (Response, error) {
	db := database.GetDatabase()

	var requestBody RequestBody
	json.Unmarshal([]byte(request.Body), &requestBody)

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
