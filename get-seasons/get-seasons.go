package main

import (
	"encoding/json"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"github.com/gavswe19/ice-api/database"
)

type Response events.APIGatewayProxyResponse

type SeasonRow struct {
	Season int `json:"season"`
}

func Handler(request events.APIGatewayProxyRequest) (Response, error) {
	db := database.GetDatabase()

	results, err := db.Query(`
	SELECT 
	DISTINCT ts.season  
	FROM team_seasons ts 
	ORDER BY ts.season DESC 
	`)

	if err != nil {
		log.Fatal(err)
	}
	defer results.Close()

	var allSeasons []int

	for results.Next() {
		var allSeason SeasonRow
		results.Scan(
			&allSeason.Season,
		)
		allSeasons = append(allSeasons, allSeason.Season)
	}

	jsonData, err := json.Marshal(allSeasons)
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
