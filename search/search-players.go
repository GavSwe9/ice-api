package main

import (
	"fmt"
	"log"

	"github.com/gavswe19/ice-api/database"
	"github.com/gavswe19/ice-api/types"
)

func searchPlayers(searchText string) []types.PlayerBio {
	db := database.GetDatabase()
	playerList := make([]types.PlayerBio, 0)

	log.Print(searchText)

	if len(searchText) == 0 {
		return playerList
	}

	query := `
		SELECT * FROM ICE.player_bio 
		WHERE full_name like ?
		LIMIT 10
	`

	log.Printf("%%%s%%", searchText)
	results, err := db.Query(query, fmt.Sprintf("%%%s%%", searchText))

	if err != nil {
		log.Fatal(err)
	}
	defer results.Close()

	for results.Next() {
		var player types.PlayerBio

		// Scan the row into the Player struct
		err := results.Scan(
			&player.PlayerId,
			&player.IsActive,
			&player.CurrentTeamId,
			&player.CurrentTeamAbbrev,
			&player.FullTeamName,
			&player.FirstName,
			&player.LastName,
			&player.FullName,
			&player.SweaterNumber,
			&player.Position,
			&player.Headshot,
			&player.HeroImage,
			&player.HeightInInches,
			&player.HeightInCentimeters,
			&player.WeightInPounds,
			&player.WeightInKilograms,
			&player.BirthDate,
			&player.BirthCity,
			&player.BirthStateProvince,
			&player.BirthCountry,
			&player.ShootsCatches,
			&player.DraftYear,
			&player.DraftTeamAbbrev,
			&player.DraftRound,
			&player.DraftPickInRound,
			&player.DraftOverallPick,
		)
		if err != nil {
			log.Fatalf("Error scanning row: %v", err)
		}

		// Append the player to the list
		playerList = append(playerList, player)
	}

	return playerList
}
