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

	var retroPicks = []RetroDraftPick{}
	for rows.Next() {
		var retro RetroDraftPick

		// Scan all columns: dp.* (20 columns) + eh.* (5 columns)
		err := rows.Scan(
			// draft_picks columns (dp.*)
			&retro.DraftYear,
			&retro.OverallPick,
			&retro.Round,
			&retro.PickInRound,
			&retro.TeamID,
			&retro.TeamAbbrev,
			&retro.TeamName,
			&retro.TeamCommonName,
			&retro.TeamPlaceName,
			&retro.TeamPickHistory,
			&retro.FirstName,
			&retro.LastName,
			&retro.PositionCode,
			&retro.CountryCode,
			&retro.Height,
			&retro.Weight,
			&retro.AmateurLeague,
			&retro.AmateurClubName,
			&retro.CreatedAt,
			&retro.UpdatedAt,
			// evolving_hockey columns (eh.*) - skip duplicate draft_year/overall_pick
			&retro.NHLID,
			&retro.TOI,
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
