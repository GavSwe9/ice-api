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
	Season    int   `json:"season"`
	TeamId    int   `json:"teamId"`
	PlayerIds []int `json:"playerIds"`
}

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

	var requestBody RequestBody
	json.Unmarshal([]byte(request.Body), &requestBody)

	query := queryPart1(requestBody.Season, requestBody.TeamId) + queryPart2(requestBody.PlayerIds) + queryPart3(requestBody.Season, requestBody.TeamId)
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

func queryPart1(season int, teamId int) string {
	return fmt.Sprintf(`
	WITH group_skater_lines AS (
	SELECT sl.*
	FROM team_season_skater_lines sl 
	WHERE sl.season = %s AND sl.team_id = %s`,
		strconv.Itoa(season), strconv.Itoa(teamId))
}

func queryPart2(playerIds []int) string {
	var whereClause []string
	for _, playerId := range playerIds {
		whereClause = append(whereClause, playerWhereClause(playerId))
	}

	return "\n AND NOT (" + strings.Join(whereClause, " AND ") + ")"
}

func queryPart3(season int, teamId int) string {
	return fmt.Sprintf(`		
	),

	plays AS (
		SELECT 
		pbp.*
		FROM group_skater_lines gsl  
		LEFT JOIN play_by_play_on_ice pbpoi on gsl.line_hash = pbpoi.line_hash 
		LEFT JOIN play_by_play pbp on pbpoi.game_pk = pbp.game_pk AND pbpoi.event_idx = pbp.event_idx   
		LEFT JOIN games g on pbp.game_pk = g.game_pk
		WHERE g.game_type = 'R'
		AND g.season = 20222023 
		AND pbp.period_type != 'SHOOTOUT'
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
		GROUP BY (CASE WHEN p.team_id = %s THEN p.team_id ELSE 0 END)
		ORDER BY (CASE WHEN p.team_id = %s THEN p.team_id ELSE 0 END) DESC
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
		strconv.Itoa(teamId),
		strconv.Itoa(teamId),
		strconv.Itoa(teamId),
		strconv.Itoa(teamId),
		strconv.Itoa(teamId),
		strconv.Itoa(season))
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

	// body := "{\"season\":20222023, \"teamId\":22, \"playerIds\": [8478402, 8477934]}"
	// Handler(events.APIGatewayProxyRequest{Body: body})
}
