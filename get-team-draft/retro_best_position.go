package main

import (
	"database/sql"
	_ "embed"
)

//go:embed queries/retro_best_position.sql
var retroSamePositionQuery string

func getRetroBestAvailableSamePosition(db *sql.DB, pick DraftPick, year int) ([]RetroDraftPick, error) {
	rows, err := db.Query(retroSamePositionQuery, year, year, pick.OverallPick, year, pick.OverallPick)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	retroPicks := []RetroDraftPick{}
	for rows.Next() {
		var retro RetroDraftPick
		err := rows.Scan(
			&retro.DraftYear,
			&retro.OverallPick,
			&retro.Round,
			&retro.PickInRound,
			&retro.TeamAbbrev,
			&retro.TeamName,
			&retro.FirstName,
			&retro.LastName,
			&retro.PositionCode,
			&retro.CountryCode,
			&retro.Height,
			&retro.Weight,
			&retro.ApiID,
			&retro.TOI,
			&retro.EvoGAR,
			&retro.EvdGAR,
			&retro.PpoGAR,
			&retro.ShdGAR,
			&retro.OffGAR,
			&retro.DefGAR,
			&retro.GAR,
			&retro.WAR,
			&retro.SPAR,
		)
		if err != nil {
			return nil, err
		}
		retroPicks = append(retroPicks, retro)
	}

	return retroPicks, nil
}
