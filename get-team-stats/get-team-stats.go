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

type LineStats struct {
	TeamName     string `json:"teamName"`
	TeamId       int    `json:"teamId"`
	Goals        int    `json:"goals"`
	Shots        int    `json:"shots"`
	Hits         int    `json:"hits"`
	Takeaways    int    `json:"takeaways"`
	BlockedShots int    `json:"blockedShots"`
	MissedShots  int    `json:"missedShots"`
	Giveaways    int    `json:"giveaways"`
	Penalties    int    `json:"penalties"`
}

func Handler(request events.APIGatewayProxyRequest) (Response, error) {
	db := database.GetDatabase()

	teamId, err := strconv.Atoi(request.PathParameters["teamId"])
	season, err := strconv.Atoi(request.PathParameters["season"])

	if err != nil {
		log.Fatal("TeamId must be integer")
	}

	query := query(teamId, season)
	results, err := db.Query(query)

	if err != nil {
		log.Fatal("Error querying line stats")
	}

	var lineStats []LineStats

	for results.Next() {
		var lineStatRow LineStats
		results.Scan(
			&lineStatRow.TeamName,
			&lineStatRow.TeamId,
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

func query(teamId int, season int) string {
	return fmt.Sprintf(`
		with plays AS (
			SELECT 
			pbp.*
			FROM play_by_play pbp  
			LEFT JOIN games g on pbp.game_pk = g.game_pk
			WHERE 
			g.season = %s AND (
				g.away_team_id = %s 
				OR g.home_team_id = %s
			)
		),

		team_stats AS (
			SELECT 
			CASE WHEN p.team_id = %s THEN %s ELSE 0 END AS team_id,
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
			CASE WHEN p.team_id = %s THEN %s ELSE 0 END)
			ORDER BY (
				CASE WHEN p.team_id = %s THEN %s ELSE 0 END) DESC
		)

		SELECT 
		CASE 
			WHEN ts.team_id = %s THEN t.name 
			ELSE 'Opponents'
		END AS team_name,
		ts.*
		FROM team_stats ts
		LEFT JOIN team_seasons t on ts.team_id = t.team_id AND t.season = %s
	`,
		strconv.Itoa(season),
		strconv.Itoa(teamId),
		strconv.Itoa(teamId),
		strconv.Itoa(teamId),
		strconv.Itoa(teamId),
		strconv.Itoa(teamId),
		strconv.Itoa(teamId),
		strconv.Itoa(teamId),
		strconv.Itoa(teamId),
		strconv.Itoa(teamId),
		strconv.Itoa(season))
}

func main() {
	lambda.Start(Handler)
}
