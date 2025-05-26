package main

import (
	"database/sql"
	_ "embed"
	"encoding/json"
	"log"
	"strconv"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/gavswe19/ice-api/database"
)

//go:embed queries/team_draft_picks.sql
var teamDraftPicksQuery string

type Response events.APIGatewayProxyResponse

type DraftPick struct {
	DraftYear       int      `json:"draftYear"`
	OverallPick     int      `json:"overallPick"`
	Round           int      `json:"round"`
	PickInRound     int      `json:"pickInRound"`
	TeamID          int      `json:"teamId"`
	TeamAbbrev      string   `json:"teamAbbrev"`
	TeamName        string   `json:"teamName"`
	TeamCommonName  string   `json:"teamCommonName"`
	TeamPlaceName   string   `json:"teamPlaceName"`
	TeamPickHistory string   `json:"teamPickHistory"`
	FirstName       string   `json:"firstName"`
	LastName        string   `json:"lastName"`
	PositionCode    string   `json:"positionCode"`
	CountryCode     string   `json:"countryCode"`
	Height          int      `json:"height"`
	Weight          int      `json:"weight"`
	AmateurLeague   string   `json:"amateurLeague"`
	AmateurClubName string   `json:"amateurClubName"`
	ApiID           *string  `json:"apiId,omitempty"`
	Position        *string  `json:"position,omitempty"`
	TOI             *float64 `json:"toi,omitempty"`
	EvoGAR          *float64 `json:"evoGar,omitempty"`
	EvdGAR          *float64 `json:"evdGar,omitempty"`
	PpoGAR          *float64 `json:"ppoGar,omitempty"`
	ShdGAR          *float64 `json:"shdGar,omitempty"`
	OffGAR          *float64 `json:"offGar,omitempty"`
	DefGAR          *float64 `json:"defGar,omitempty"`
	GAR             *float64 `json:"gar,omitempty"`
	WAR             *float64 `json:"war,omitempty"`
	SPAR            *float64 `json:"spar,omitempty"`
}

type RetroDraftPick struct {
	DraftYear    int     `json:"draftYear"`
	OverallPick  int     `json:"overallPick"`
	Round        int     `json:"round"`
	PickInRound  int     `json:"pickInRound"`
	TeamAbbrev   string  `json:"teamAbbrev"`
	TeamName     string  `json:"teamName"`
	FirstName    string  `json:"firstName"`
	LastName     string  `json:"lastName"`
	PositionCode string  `json:"positionCode"`
	CountryCode  string  `json:"countryCode"`
	Height       int     `json:"height"`
	Weight       int     `json:"weight"`
	ApiID        string  `json:"apiId"`
	TOI          float64 `json:"toi"`
	EvoGAR       float64 `json:"evoGar"`
	EvdGAR       float64 `json:"evdGar"`
	PpoGAR       float64 `json:"ppoGar"`
	ShdGAR       float64 `json:"shdGar"`
	OffGAR       float64 `json:"offGar"`
	DefGAR       float64 `json:"defGar"`
	GAR          float64 `json:"gar"`
	WAR          float64 `json:"war"`
	SPAR         float64 `json:"spar"`
}

type DraftPickWithRetro struct {
	DraftPick                      DraftPick        `json:"draftPick"`
	RetroBestAvailableOverall      []RetroDraftPick `json:"retroBestAvailableOverall"`
	RetroBestAvailableSamePosition []RetroDraftPick `json:"retroBestAvailableSamePosition"`
}

type TeamDraftResponse struct {
	TeamAbbrev string               `json:"teamAbbrev"`
	DraftYear  int                  `json:"draftYear"`
	DraftPicks []DraftPickWithRetro `json:"draftPicks"`
}

func Handler(request events.APIGatewayProxyRequest) (Response, error) {
	log.Printf("NEW VERSION - Using embedded queries") // Add this line
	db := database.GetDatabase()
	defer db.Close()

	// Parse path parameters
	teamAbv := request.PathParameters["teamAbv"]
	yearStr := request.PathParameters["year"]

	if teamAbv == "" {
		return errorResponse(400, "Invalid team abbreviation"), nil
	}

	year, err := strconv.Atoi(yearStr)
	if err != nil {
		return errorResponse(400, "Invalid year"), nil
	}

	// Get team's draft picks for the year
	draftPicks, err := getTeamDraftPicks(db, teamAbv, year)
	if err != nil {
		log.Printf("Error getting draft picks: %v", err)
		return errorResponse(500, "Internal server error"), nil
	}

	if len(draftPicks) == 0 {
		return errorResponse(404, "No draft picks found for team/year"), nil
	}

	// For each draft pick, get retro best available
	var draftPicksWithRetro []DraftPickWithRetro
	for _, pick := range draftPicks {
		retroOverall, err := getRetroBestAvailableOverall(db, pick, year)
		if err != nil {
			log.Printf("Error getting retro overall for pick %d: %v", pick.OverallPick, err)
			retroOverall = []RetroDraftPick{}
		}

		retroSamePosition, err := getRetroBestAvailableSamePosition(db, pick, year)
		if err != nil {
			log.Printf("Error getting retro same position for pick %d: %v", pick.OverallPick, err)
			retroSamePosition = []RetroDraftPick{}
		}

		draftPicksWithRetro = append(draftPicksWithRetro, DraftPickWithRetro{
			DraftPick:                      pick,
			RetroBestAvailableOverall:      retroOverall,
			RetroBestAvailableSamePosition: retroSamePosition,
		})
	}

	response := TeamDraftResponse{
		TeamAbbrev: teamAbv,
		DraftYear:  year,
		DraftPicks: draftPicksWithRetro,
	}

	jsonData, err := json.Marshal(response)
	if err != nil {
		log.Printf("Error marshaling response: %v", err)
		return errorResponse(500, "Internal server error"), nil
	}

	return Response{
		StatusCode:      200,
		IsBase64Encoded: false,
		Body:            string(jsonData),
		Headers: map[string]string{
			"Content-Type":                     "application/json",
			"X-MyCompany-Func-Reply":           "draft-handler",
			"Access-Control-Allow-Origin":      "*",
			"Access-Control-Allow-Credentials": "true",
		},
	}, nil
}

func getTeamDraftPicks(db *sql.DB, teamAbv string, year int) ([]DraftPick, error) {
	rows, err := db.Query(teamDraftPicksQuery, year, teamAbv, year)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var picks []DraftPick
	for rows.Next() {
		var pick DraftPick
		var apiID sql.NullString
		var position sql.NullString
		var toi sql.NullFloat64
		var evoGAR sql.NullFloat64
		var evdGAR sql.NullFloat64
		var ppoGAR sql.NullFloat64
		var shdGAR sql.NullFloat64
		var offGAR sql.NullFloat64
		var defGAR sql.NullFloat64
		var gar sql.NullFloat64
		var war sql.NullFloat64
		var spar sql.NullFloat64

		err := rows.Scan(
			&pick.DraftYear,
			&pick.OverallPick,
			&pick.Round,
			&pick.PickInRound,
			&pick.TeamID,
			&pick.TeamAbbrev,
			&pick.TeamName,
			&pick.TeamCommonName,
			&pick.TeamPlaceName,
			&pick.TeamPickHistory,
			&pick.FirstName,
			&pick.LastName,
			&pick.PositionCode,
			&pick.CountryCode,
			&pick.Height,
			&pick.Weight,
			&pick.AmateurLeague,
			&pick.AmateurClubName,
			&apiID,
			&position,
			&toi,
			&evoGAR,
			&evdGAR,
			&ppoGAR,
			&shdGAR,
			&offGAR,
			&defGAR,
			&gar,
			&war,
			&spar,
		)
		if err != nil {
			return nil, err
		}

		// Convert SQL null types to pointers
		if apiID.Valid {
			pick.ApiID = &apiID.String
		}
		if position.Valid {
			pick.Position = &position.String
		}
		if toi.Valid {
			pick.TOI = &toi.Float64
		}
		if evoGAR.Valid {
			pick.EvoGAR = &evoGAR.Float64
		}
		if evdGAR.Valid {
			pick.EvdGAR = &evdGAR.Float64
		}
		if ppoGAR.Valid {
			pick.PpoGAR = &ppoGAR.Float64
		}
		if shdGAR.Valid {
			pick.ShdGAR = &shdGAR.Float64
		}
		if offGAR.Valid {
			pick.OffGAR = &offGAR.Float64
		}
		if defGAR.Valid {
			pick.DefGAR = &defGAR.Float64
		}
		if gar.Valid {
			pick.GAR = &gar.Float64
		}
		if war.Valid {
			pick.WAR = &war.Float64
		}
		if spar.Valid {
			pick.SPAR = &spar.Float64
		}

		picks = append(picks, pick)
	}

	return picks, nil
}

func errorResponse(statusCode int, message string) Response {
	body := map[string]string{"error": message}
	jsonBody, _ := json.Marshal(body)

	return Response{
		StatusCode:      statusCode,
		IsBase64Encoded: false,
		Body:            string(jsonBody),
		Headers: map[string]string{
			"Content-Type":                     "application/json",
			"X-MyCompany-Func-Reply":           "draft-handler",
			"Access-Control-Allow-Origin":      "*",
			"Access-Control-Allow-Credentials": "true",
		},
	}
}

func main() {
	lambda.Start(Handler)
}
