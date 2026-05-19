package whoop

import "time"

type DateRange struct {
	From string `json:"from" example:"2026-05-01" doc:"Inclusive start date in YYYY-MM-DD format."`
	To   string `json:"to" example:"2026-05-19" doc:"Inclusive end date in YYYY-MM-DD format."`
}

type Profile struct {
	UserID    string `json:"user_id" doc:"WHOOP user identifier."`
	Email     string `json:"email,omitempty" doc:"Email from the WHOOP profile, when available."`
	FirstName string `json:"first_name,omitempty" doc:"First name from the WHOOP profile."`
	LastName  string `json:"last_name,omitempty" doc:"Last name from the WHOOP profile."`
}

type DailyMetrics struct {
	Date                       string   `json:"date" example:"2026-05-19" doc:"Calendar date represented by these metrics."`
	Strain                     *float64 `json:"strain,omitempty" doc:"WHOOP day strain, 0-21 scale." example:"10.4"`
	AverageHeartRate           *int     `json:"average_heart_rate,omitempty" doc:"Average heart rate for the cycle, in bpm." example:"71"`
	MaxHeartRate               *int     `json:"max_heart_rate,omitempty" doc:"Maximum heart rate for the cycle, in bpm." example:"168"`
	RecoveryScore              *int     `json:"recovery_score,omitempty" doc:"WHOOP recovery score, 0-100 scale." example:"67"`
	RestingHeartRate           *int     `json:"resting_heart_rate,omitempty" doc:"Resting heart rate from recovery, in bpm." example:"54"`
	HRVRMSSDMilli              *float64 `json:"hrv_rmssd_milli,omitempty" doc:"HRV RMSSD in milliseconds." example:"42.7"`
	SleepPerformancePercentage *float64 `json:"sleep_performance_percentage,omitempty" doc:"Sleep performance percentage, 0-100." example:"84.5"`
	SleepEfficiencyPercentage  *float64 `json:"sleep_efficiency_percentage,omitempty" doc:"Sleep efficiency percentage, 0-100." example:"91.2"`
	SleepDurationHours         *float64 `json:"sleep_duration_hours,omitempty" doc:"Approximate sleep duration in hours, based on in-bed time." example:"7.6"`
	WorkoutCount               int      `json:"workout_count" doc:"Number of workouts started on this date." example:"1"`
	WorkoutStrain              *float64 `json:"workout_strain,omitempty" doc:"Sum of workout strain values for this date." example:"8.2"`
}

type TrendSummary struct {
	Days                             int      `json:"days" doc:"Number of calendar days included in the summary." example:"14"`
	RecoveryAverage                  *float64 `json:"recovery_average,omitempty" doc:"Average recovery score across days with recovery data." example:"64.3"`
	StrainAverage                    *float64 `json:"strain_average,omitempty" doc:"Average day strain across days with cycle data." example:"9.8"`
	SleepPerformanceAverage          *float64 `json:"sleep_performance_average,omitempty" doc:"Average sleep performance across days with sleep data." example:"82.1"`
	SleepDurationAverageHours        *float64 `json:"sleep_duration_average_hours,omitempty" doc:"Average sleep duration in hours across days with sleep data." example:"7.2"`
	WorkoutDays                      int      `json:"workout_days" doc:"Number of days with at least one workout." example:"5"`
	WorkoutCount                     int      `json:"workout_count" doc:"Total workouts in the range." example:"7"`
	LatestRecoveryScore              *int     `json:"latest_recovery_score,omitempty" doc:"Recovery score from the most recent non-empty day." example:"71"`
	LatestStrain                     *float64 `json:"latest_strain,omitempty" doc:"Strain from the most recent non-empty day." example:"11.1"`
	LatestSleepPerformancePercentage *float64 `json:"latest_sleep_performance_percentage,omitempty" doc:"Sleep performance from the most recent non-empty day." example:"86.4"`
}

type ContextPack struct {
	GeneratedAt         time.Time         `json:"generated_at" doc:"Timestamp when this API response was generated."`
	DataFreshness       *time.Time        `json:"data_freshness,omitempty" doc:"Most recent updated_at timestamp across synced WHOOP tables."`
	Range               DateRange         `json:"range" doc:"Calendar range included in the context pack."`
	Profile             *Profile          `json:"profile,omitempty" doc:"Latest synced WHOOP user profile."`
	Latest              *DailyMetrics     `json:"latest,omitempty" doc:"Most recent non-empty daily metrics in the requested range."`
	Trends              TrendSummary      `json:"trends" doc:"Compact trend summary for the requested range."`
	RecentDays          []DailyMetrics    `json:"recent_days" doc:"Daily metrics ordered newest first."`
	InterpretationHints map[string]string `json:"interpretation_hints" doc:"Short field interpretation hints for agents."`
}
