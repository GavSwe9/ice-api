WITH
	eh_totals AS (
		SELECT
			eh.draft_year,
			eh.overall_pick,
			eh.nhl_id,
			SUM(toi_all) AS toi,
			SUM(gar) AS gar,
			SUM(war) AS war,
			SUM(spar) AS spar
		FROM
			ICE.evolving_hockey_player_seasons_gar eh
		WHERE 
			eh.draft_year = ?
		GROUP BY
			eh.draft_year,
			eh.overall_pick,
			eh.nhl_id
	)
	
SELECT
	dp.draft_year,
	dp.overall_pick,
	dp.round,
	dp.pick_in_round,
	dp.team_id,
	dp.team_abbrev,
	dp.team_name,
	dp.team_common_name,
	dp.team_place_name,
	dp.team_pick_history,
	dp.first_name,
	dp.last_name,
	dp.position_code,
	dp.country_code,
	dp.height,
	dp.weight,
	dp.amateur_league,
	dp.amateur_club_name,
	dp.created_at,
	dp.updated_at,
	eh.nhl_id,
	eh.toi,
	eh.gar,
	eh.war,
	eh.spar
FROM
	ICE.draft_picks dp
	LEFT JOIN eh_totals eh ON dp.draft_year = eh.draft_year
	AND dp.overall_pick = eh.overall_pick
WHERE
	dp.team_abbrev = ?
	AND dp.draft_year = ?
ORDER BY dp.overall_pick