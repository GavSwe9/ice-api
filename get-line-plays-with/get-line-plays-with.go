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

type LinePlaysWithPlayer struct {
	PlayerId         int     `json:"player_id"`
	PercentageEvents float32 `json:"percentage_events"`
	PlayerName       string  `json:"player_name"`
	Position         string  `json:"position"`
}

func Handler(request events.APIGatewayProxyRequest) (Response, error) {
	db := database.GetDatabase()

	var requestBody RequestBody
	json.Unmarshal([]byte(request.Body), &requestBody)

	query := queryPart1(requestBody.Season) + queryPart2(requestBody.PlayerIds) + queryPart3(requestBody.TeamId, requestBody.PlayerIds)
	results, err := db.Query(query)

	if err != nil {
		log.Fatal("Error querying line stats")
	}

	var linePlaysWithPlayerList []LinePlaysWithPlayer

	for results.Next() {
		var linePlaysWithPlayerRow LinePlaysWithPlayer
		results.Scan(
			&linePlaysWithPlayerRow.PlayerId,
			&linePlaysWithPlayerRow.PercentageEvents,
			&linePlaysWithPlayerRow.PlayerName,
			&linePlaysWithPlayerRow.Position,
		)
		linePlaysWithPlayerList = append(linePlaysWithPlayerList, linePlaysWithPlayerRow)
	}

	jsonData, err := json.Marshal(linePlaysWithPlayerList)
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

func queryPart1(season int) string {
	return fmt.Sprintf(`
	WITH group_skater_lines AS (
	SELECT sl.*
	FROM season_skater_lines sl 
	WHERE sl.season = %s `, strconv.Itoa(season))
}

func queryPart2(playerIds []int) string {
	var whereClause []string
	for _, playerId := range playerIds {
		whereClause = append(whereClause, playerWhereClause(playerId))
	}

	return "\n AND " + strings.Join(whereClause, " AND ")
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

func queryPart3(team_id int, playerIds []int) string {
	return fmt.Sprintf(`
		), 
		
		plays AS (
			SELECT 
			
			gsl.skater_id_1,
			gsl.skater_id_2,
			gsl.skater_id_3,
			gsl.skater_id_4,
			gsl.skater_id_5,
			gsl.skater_id_6
			
			FROM group_skater_lines gsl  
			LEFT JOIN play_by_play_on_ice pbpoi on gsl.line_hash = pbpoi.line_hash and pbpoi.team_id = %s
		),

		all_players AS (
			SELECT  
			p.skater_id_1 AS skater_id 
			FROM plays p
			
			UNION ALL
			
			SELECT 
			p.skater_id_2
			FROM plays p
			
			UNION ALL  
			
			SELECT 
			p.skater_id_3
			FROM plays p 
			
			UNION ALL
			
			SELECT 
			p.skater_id_4
			FROM plays p
			
			UNION ALL 
			
			SELECT 
			p.skater_id_5
			FROM plays p
			
			UNION ALL 
			
			SELECT 
			p.skater_id_6
			FROM plays p 
		),

		skater_counts AS (
			SELECT 
			ap.skater_id,
			COUNT(ap.skater_id) as total
			FROM all_players ap
			WHERE ap.skater_id != 0 
			GROUP BY ap.skater_id
			ORDER BY COUNT(ap.skater_id) DESC
		),

		total_events AS (
			SELECT *
			FROM skater_counts sc
			WHERE sc.skater_id = %s
		)

		SELECT 
		sc.skater_id as player_id,
		sc.total / (SELECT total from total_events) AS percentage_events,
		p.player_name,
		p.position
		FROM skater_counts sc
		LEFT JOIN players p on sc.skater_id = p.player_id 
		WHERE sc.skater_id NOT IN (%s)
	`,
		strconv.Itoa(team_id),
		strconv.Itoa(playerIds[0]),
		strings.Trim(strings.Join(strings.Split(fmt.Sprint(playerIds), " "), ", "), "[]"))
}

func main() {
	lambda.Start(Handler)

	// body := "{\"season\":20222023, \"teamId\":22, \"playerIds\": [8478402, 8477934]}"

	// Handler(events.APIGatewayProxyRequest{Body: body})
}
