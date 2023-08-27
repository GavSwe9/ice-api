package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"github.com/gavswe19/ice-api/database"
)

type Response events.APIGatewayProxyResponse

type RequestBody struct {
	TeamId    int   `json:"teamId"`
	PlayerIds []int `json:"playerIds"`
}

type LineStats struct {
	// TEAM|GOAL|SHOT|HIT|TAKEAWAY|BLOCKED_SHOT|MISSED_SHOT|GIVEAWAY|PENALTY|MISSED_SHOT
	Team  int `json:"team"`
	Goals int `json:"goals"`
	Shots int `json:"shots"`
	Hits  int `json:"hits"`

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

	query := queryPart1() + queryPart2(requestBody.PlayerIds) + queryPart3(requestBody.TeamId)
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

func queryPart1() string {
	return `
	WITH group_skater_lines AS (
	SELECT sl.*
	FROM skater_lines sl `
}

func queryPart2(playerIds []int) string {
	var whereClause []string
	for _, playerId := range playerIds {
		whereClause = append(whereClause, playerWhereClause(playerId))
	}

	return "\nWHERE " + strings.Join(whereClause, " AND ")
}

func queryPart3(teamId int) string {
	return fmt.Sprintf(`		
	),

	plays AS (
		SELECT 
		pbp.*
		FROM group_skater_lines gsl  
		LEFT JOIN play_by_play_on_ice pbpoi on gsl.line_hash = pbpoi.line_hash 
		LEFT JOIN play_by_play pbp on pbpoi.game_pk = pbp.game_pk AND pbpoi.event_idx = pbp.event_idx 
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
		strconv.Itoa(teamId))
}

func playerWhereClause(playerId int) string {
	return fmt.Sprintf(`
	(
		sl.skater_id_1 = %s 
		OR sl.skater_id_2 = %s
		OR sl.skater_id_3 = %s
		OR sl.skater_id_4 = %s
		OR sl.skater_id_5 = %s
		OR sl.skater_id_6 = %s 
	)
	`,
		strconv.Itoa(playerId),
		strconv.Itoa(playerId),
		strconv.Itoa(playerId),
		strconv.Itoa(playerId),
		strconv.Itoa(playerId),
		strconv.Itoa(playerId))
}

func main() {
	lambda.Start(Handler)
}
