package main

import (
	"fmt"
	"log"

	"github.com/gavswe19/ice-api/database"
	"github.com/gavswe19/ice-api/types"
)

func searchTeams(searchText string) []types.TeamSeason {
	db := database.GetDatabase()
	teamList := make([]types.TeamSeason, 0)

	log.Print(searchText)

	if len(searchText) == 0 {
		return teamList
	}

	query := `
		SELECT 
		*
		FROM ICE.team_seasons_2 ts 
		WHERE
		ts.season_id = 20232024
		and (
			team_name like ?
			or team_abbrev like ?
		)
	`

	searchWild := fmt.Sprintf("%%%s%%", searchText)
	results, err := db.Query(query, searchWild, searchWild)

	if err != nil {
		log.Fatal(err)
	}
	defer results.Close()

	for results.Next() {
		var teamSeason types.TeamSeason

		// Scan the row into the Player struct
		err := results.Scan(
			&teamSeason.ID,
			&teamSeason.TeamAbbrev,
			&teamSeason.SeasonID,
			&teamSeason.ConferenceName,
			&teamSeason.DivisionName,
			&teamSeason.TeamName,
			&teamSeason.TeamCommonName,
		)
		if err != nil {
			log.Fatalf("Error scanning row: %v", err)
		}

		// Append the player to the list
		teamList = append(teamList, teamSeason)
	}

	return teamList
}
