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
	TeamId int `json:"teamId"`
}

type LineStats struct {
	// TEAM|GOAL|SHOT|HIT|TAKEAWAY|BLOCKED_SHOT|MISSED_SHOT|GIVEAWAY|PENALTY|MISSED_SHOT
	Team         int `json:"team"`
	Goals        int `json:"goals"`
	Shots        int `json:"shots"`
	Hits         int `json:"hits"`
	Takeaways    int `json:"takeaways"`
	BlockedShots int `json:"blocked_shots"`
	MissedShots  int `json:"missed_shots"`
	Giveaways    int `json:"giveaways"`
	Penalties    int `json:"penalties"`
}

func Handler(request events.APIGatewayProxyRequest) (Response, error) {
	db := database.GetDatabase()

	var requestBody RequestBody
	json.Unmarshal([]byte(request.Body), &requestBody)

	teamId, err := strconv.Atoi(request.PathParameters["teamId"])
	if err != nil {
		log.Fatal("TeamId must be integer")
	}

	query := query(teamId)
	results, err := db.Query(query)

	if err != nil {
		log.Fatal("Error querying line stats")
	}

	var lineStats []LineStats

	for results.Next() {
		var lineStatRow LineStats
		results.Scan(
			&lineStatRow.Team,
			&lineStatRow.Goals,
			&lineStatRow.Shots,
			&lineStatRow.Hits,
			&lineStatRow.Takeaways,
			&lineStatRow.BlockedShots,
			&lineStatRow.MissedShots,
			&lineStatRow.Giveaways,
			&lineStatRow.Penalties,
		)
		lineStats = append(lineStats, lineStatRow)
	}

	jsonData, err := json.Marshal(lineStats)
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

func query(teamId int) string {
	return fmt.Sprintf(`
	with plays AS (
		SELECT 
		pbp.*
		FROM play_by_play pbp  
		LEFT JOIN games g on pbp.game_pk = g.game_pk
		WHERE 
		g.away_team_id = %s
		OR g.home_team_id = %s
	)

	SELECT 
	CASE WHEN p.team_id = %s THEN 0 ELSE 1 END AS team,
	SUM(CASE WHEN p.event_type_id = 'GOAL' THEN 1 ELSE 0 END) AS goal,
	SUM(CASE WHEN p.event_type_id = 'SHOT' THEN 1 ELSE 0 END) AS shot,
	SUM(CASE WHEN p.event_type_id = 'HIT' THEN 1 ELSE 0 END) AS hit,
	SUM(CASE WHEN p.event_type_id = 'TAKEAWAY' THEN 1 ELSE 0 END) AS takeaway,
	SUM(CASE WHEN p.event_type_id = 'BLOCKED_SHOT' THEN 1 ELSE 0 END) AS blocked_shot,
	SUM(CASE WHEN p.event_type_id = 'MISSED_SHOT' THEN 1 ELSE 0 END) AS missed_shot,
	SUM(CASE WHEN p.event_type_id = 'GIVEAWAY' THEN 1 ELSE 0 END) AS giveaway,
	SUM(CASE WHEN p.event_type_id = 'PENALTY' THEN 1 ELSE 0 END) AS penalty
	FROM plays p
	GROUP BY (
	CASE WHEN p.team_id = %s THEN 0 ELSE 1 END)
	`,
		strconv.Itoa(teamId),
		strconv.Itoa(teamId),
		strconv.Itoa(teamId),
		strconv.Itoa(teamId))
}

func main() {
	lambda.Start(Handler)
}
