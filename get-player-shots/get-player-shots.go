package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"github.com/gavswe19/ice-api/database"
)

type Response events.APIGatewayProxyResponse

type ShotRow struct {
	PlayerId       int       `json:"player_id"`
	PlayerType     string    `json:"player_type"`
	GamePk         int       `json:"game_pk"`
	EventIdx       int       `json:"event_idx"`
	EventId        int       `json:"event_id"`
	Period         int       `json:"period"`
	PeriodType     string    `json:"periodType"`
	PeriodTime     string    `json:"periodTime"`
	DateTime       time.Time `json:"date_time"`
	GoalsAway      int       `json:"goals_away"`
	GoalsHome      int       `json:"goals_home"`
	Event          string    `json:"event"`
	EventCode      string    `json:"event_code"`
	EventTypeId    string    `json:"event_type_id"`
	Description    string    `json:"description"`
	SecondaryType  string    `json:"secondary_type"`
	TeamId         int       `json:"team_id"`
	XCoordinate    float32   `json:"x_coordinate"`
	YCoordinate    float32   `json:"y_coordinate"`
	AdjXCoordinate float32   `json:"adj_x_coordinate"`
	AdjYCoordinate float32   `json:"adj_y_coordinate"`
}

func Handler(request events.APIGatewayProxyRequest) (Response, error) {
	// func Handler() () {
	db := database.GetDatabase()

	playerId, err := strconv.Atoi(request.PathParameters["playerId"])
	fmt.Println(playerId)
	if err != nil {
		log.Fatal("PlayerId must be integer")
	}

	fmt.Println(playerId)
	fmt.Println(" -------------------- ")
	results, err := db.Query(fmt.Sprintf(
		`SELECT 
		c.player_id,
		c.player_type,
		p.*
		FROM play_by_play_contributor c
		LEFT JOIN play_by_play p ON c.game_pk = p.game_pk AND c.event_idx = p.event_idx 
		WHERE c.player_id = %s
		AND p.event_type_id <> 'BLOCKED_SHOT'
		AND c.player_type = 'Shooter'`,
		strconv.Itoa(playerId)))

	if err != nil {
		log.Fatal(err)
	}
	defer results.Close()

	var allPlays []ShotRow

	for results.Next() {
		var shotRow ShotRow
		results.Scan(
			&shotRow.PlayerId,
			&shotRow.PlayerType,
			&shotRow.GamePk,
			&shotRow.EventIdx,
			&shotRow.EventId,
			&shotRow.Period,
			&shotRow.PeriodType,
			&shotRow.PeriodTime,
			&shotRow.DateTime,
			&shotRow.GoalsAway,
			&shotRow.GoalsHome,
			&shotRow.Event,
			&shotRow.EventCode,
			&shotRow.EventTypeId,
			&shotRow.Description,
			&shotRow.SecondaryType,
			&shotRow.TeamId,
			&shotRow.XCoordinate,
			&shotRow.YCoordinate,
			&shotRow.AdjXCoordinate,
			&shotRow.AdjYCoordinate,
		)

		if err != nil {
			// Handle the error
			fmt.Println("Error scanning the row:", err)
		}
		allPlays = append(allPlays, shotRow)
	}

	jsonData, err := json.Marshal(allPlays)
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

	// playerMap := make(map[string]string)
	// playerMap["playerId"] = "8478402"

	// Handler(events.APIGatewayProxyRequest{PathParameters: playerMap})
}
