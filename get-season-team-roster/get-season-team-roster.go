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

type SeasonTeamRosterPlayer struct {
	PlayerId int    `json:"playerId"`
	FullName string `json:"fullName"`
	Position string `json:"position"`
}

func Handler(request events.APIGatewayProxyRequest) (Response, error) {
	db := database.GetDatabase()

	season, err := strconv.Atoi(request.PathParameters["season"])

	if err != nil {
		log.Fatal("Season must be an integer")
	}

	team_id, err := strconv.Atoi(request.PathParameters["teamId"])

	if err != nil {
		log.Fatal("TeamId must be an integer")
	}

	query := query(season, team_id)
	results, err := db.Query(query)

	if err != nil {
		log.Fatal("Error querying season team roster")
	}

	var seasonTeamRoster []SeasonTeamRosterPlayer

	for results.Next() {
		var seasonTeamRosterPlayer SeasonTeamRosterPlayer
		results.Scan(
			&seasonTeamRosterPlayer.PlayerId,
			&seasonTeamRosterPlayer.FullName,
			&seasonTeamRosterPlayer.Position,
		)
		seasonTeamRoster = append(seasonTeamRoster, seasonTeamRosterPlayer)
	}

	jsonData, err := json.Marshal(seasonTeamRoster)
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

func query(season int, teamId int) string {
	return fmt.Sprintf(`
	SELECT 
	p.*
	FROM team_season_players tsp 
	LEFT JOIN players p ON tsp.player_id = p.player_id 
	WHERE 
	tsp.season = %s AND 
	tsp.team_id = %s
	`,
		strconv.Itoa(season),
		strconv.Itoa(teamId))
}

func main() {
	lambda.Start(Handler)
}

// 7685
