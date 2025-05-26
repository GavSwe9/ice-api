WITH
	eh_totals AS (
		SELECT
			eh.draft_yr,
			eh.draft_ov,
			eh.api_id,
			SUM(toi_all) AS toi,
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
			eh.api_id
	),
	target_stat AS (
		SELECT
			COALESCE(eh.spar, -1000) AS benchmark_stat
		FROM
			ICE.draft_picks dp
			LEFT JOIN eh_totals eh ON dp.draft_year = eh.draft_yr
			AND dp.overall_pick = eh.draft_ov
		WHERE
			dp.draft_year = ?
			AND dp.overall_pick = ?
	)
SELECT
	dp.draft_year,
	dp.overall_pick,
	dp.round,
	dp.pick_in_round,
	dp.team_abbrev,
	dp.team_name,
	dp.first_name,
	dp.last_name,
	dp.position_code,
	dp.country_code,
	dp.height,
	dp.weight,
	eh.api_id,
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
	CROSS JOIN target_stat ts
WHERE
	dp.draft_year = ?
	AND dp.overall_pick > ?
	AND COALESCE(eh.spar, -1000) > ts.benchmark_stat
ORDER BY
	eh.spar DESC
LIMIT 25
