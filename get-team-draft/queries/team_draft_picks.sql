WITH
	eh_totals AS (
		SELECT
			eh.draft_yr,
			eh.draft_ov,
			eh.api_id,
			eh.position,
			SUM(toi_all) as toi,
			SUM(evo_gar) AS evo_gar,
			SUM(evd_gar) AS evd_gar,
			SUM(ppo_gar) AS ppo_gar,
			SUM(shd_gar) AS shd_gar,
			SUM(off_gar) AS off_gar,
			SUM(def_gar) AS def_gar,
			SUM(gar) AS gar,
			SUM(war) AS war,
			SUM(spar) AS spar
		FROM
			ICE.evolving_hockey_player_seasons eh
		WHERE 
			eh.draft_yr = ?
		GROUP BY
			eh.draft_yr,
			eh.draft_ov,
			eh.api_id,
			eh.position
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
	eh.api_id,
	eh.position,
	eh.toi,
	eh.evo_gar,
	eh.evd_gar,
	eh.ppo_gar,
	eh.shd_gar,
	eh.off_gar,
	eh.def_gar,
	eh.gar,
	eh.war,
	eh.spar
FROM
	ICE.draft_picks dp
	LEFT JOIN eh_totals eh ON dp.draft_year = eh.draft_yr
	AND dp.overall_pick = eh.draft_ov
WHERE
	dp.team_abbrev = ?
	AND dp.draft_year = ?
ORDER BY dp.overall_pick
