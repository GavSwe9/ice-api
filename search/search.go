package main

import (
	"encoding/json"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"github.com/gavswe19/ice-api/types"
)

func Handler(request events.APIGatewayProxyRequest) (types.Response, error) {
	searchResult := types.SearchResult{
		ItemsFound: false,
		Teams:      make([]types.TeamSeason, 0),
		Players:    make([]types.PlayerBio, 0),
	}

	searchText := request.PathParameters["searchText"]
	log.Print(searchText)

	if len(searchText) == 0 {
		jsonData, err := json.Marshal(searchResult)
		if err != nil {
			log.Fatal(err)
		}

		resp := types.Response{
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
	searchResult.Teams = searchTeams(searchText)
	searchResult.Players = searchPlayers(searchText)
	searchResult.ItemsFound = len(searchResult.Teams) > 0 || len(searchResult.Players) > 0

	jsonData, err := json.Marshal(searchResult)
	if err != nil {
		log.Fatal(err)
	}

	resp := types.Response{
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
