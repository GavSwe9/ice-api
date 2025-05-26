package main

import (
	"encoding/json"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	types "github.com/gavswe19/ice-api/Types"
	"github.com/gavswe19/ice-api/database"
)

type Response events.APIGatewayProxyResponse

func Handler(request events.APIGatewayProxyRequest) (Response, error) {
	db := database.GetDatabase()

	playerId := request.PathParameters["playerId"]

	query := `
		SELECT * FROM ICE.player_bio 
		WHERE player_id = ?
	`

	results, err := db.Query(query, playerId)

	if err != nil {
		log.Fatal(err)
	}
	defer results.Close()

	var player types.PlayerBio

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
	// Scan the row into the Player struct
	err = results.Scan(
		&player.PlayerId,
		&player.IsActive,
		&player.CurrentTeamId,
		&player.CurrentTeamAbbrev,
		&player.FullTeamName,
		&player.FirstName,
		&player.LastName,
		&player.FullName,
		&player.SweaterNumber,
		&player.Position,
		&player.Headshot,
		&player.HeroImage,
		&player.HeightInInches,
		&player.HeightInCentimeters,
		&player.WeightInPounds,
		&player.WeightInKilograms,
		&player.BirthDate,
		&player.BirthCity,
		&player.BirthStateProvince,
		&player.BirthCountry,
		&player.ShootsCatches,
		&player.DraftYear,
		&player.DraftTeamAbbrev,
		&player.DraftRound,
		&player.DraftPickInRound,
		&player.DraftOverallPick,
	)
	if err != nil {
		log.Fatalf("Error scanning row: %v", err)
	}

	jsonData, err := json.Marshal(player)
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
