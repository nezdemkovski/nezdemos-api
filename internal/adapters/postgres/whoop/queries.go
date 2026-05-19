package whoop

const profileSQL = `
SELECT user_id, COALESCE(email, ''), COALESCE(first_name, ''), COALESCE(last_name, '')
FROM whoop_user_profile
ORDER BY updated_at DESC
LIMIT 1
`

const dataFreshnessSQL = `
SELECT MAX(updated_at) FROM (
	SELECT updated_at FROM whoop_user_profile
	UNION ALL SELECT updated_at FROM whoop_body_measurement
	UNION ALL SELECT updated_at FROM whoop_cycle
	UNION ALL SELECT updated_at FROM whoop_recovery
	UNION ALL SELECT updated_at FROM whoop_sleep
	UNION ALL SELECT updated_at FROM whoop_workout
) AS updates
`

const dailySQL = `
WITH bounds AS (
	SELECT $1::date AS from_date, $2::date AS to_date
),
calendar AS (
	SELECT generate_series(from_date, to_date, interval '1 day')::date AS day
	FROM bounds
),
cycles AS (
	SELECT
		start_time::date AS day,
		AVG(score_strain)::float8 AS strain,
		ROUND(AVG(score_average_heart_rate))::int AS average_heart_rate,
		MAX(score_max_heart_rate)::int AS max_heart_rate
	FROM whoop_cycle, bounds
	WHERE start_time::date BETWEEN from_date AND to_date
	GROUP BY 1
),
recoveries AS (
	SELECT
		c.start_time::date AS day,
		ROUND(AVG(r.score_recovery_score))::int AS recovery_score,
		ROUND(AVG(r.score_resting_heart_rate))::int AS resting_heart_rate,
		AVG(r.score_hrv_rmssd_milli)::float8 AS hrv_rmssd_milli
	FROM whoop_recovery r
	JOIN whoop_cycle c ON c.id = r.cycle_id
	JOIN bounds ON true
	WHERE c.start_time::date BETWEEN from_date AND to_date
	GROUP BY 1
),
sleeps AS (
	SELECT
		COALESCE(end_time, start_time)::date AS day,
		AVG(score_sleep_performance_percentage)::float8 AS sleep_performance_percentage,
		AVG(score_sleep_efficiency_percentage)::float8 AS sleep_efficiency_percentage,
		AVG(score_stage_summary_total_in_bed_time_milli)::float8 / 1000 / 60 / 60 AS sleep_duration_hours
	FROM whoop_sleep, bounds
	WHERE COALESCE(end_time, start_time)::date BETWEEN from_date AND to_date
	GROUP BY 1
),
workouts AS (
	SELECT
		start_time::date AS day,
		COUNT(*)::int AS workout_count,
		SUM(score_strain)::float8 AS workout_strain
	FROM whoop_workout, bounds
	WHERE start_time::date BETWEEN from_date AND to_date
	GROUP BY 1
)
SELECT
	calendar.day::text,
	cycles.strain,
	cycles.average_heart_rate,
	cycles.max_heart_rate,
	recoveries.recovery_score,
	recoveries.resting_heart_rate,
	recoveries.hrv_rmssd_milli,
	sleeps.sleep_performance_percentage,
	sleeps.sleep_efficiency_percentage,
	sleeps.sleep_duration_hours,
	COALESCE(workouts.workout_count, 0),
	workouts.workout_strain
FROM calendar
LEFT JOIN cycles ON cycles.day = calendar.day
LEFT JOIN recoveries ON recoveries.day = calendar.day
LEFT JOIN sleeps ON sleeps.day = calendar.day
LEFT JOIN workouts ON workouts.day = calendar.day
ORDER BY calendar.day DESC
`
