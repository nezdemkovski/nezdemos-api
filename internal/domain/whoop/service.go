package whoop

import (
	"context"
	"fmt"
	"math"
	"time"
)

type Repository interface {
	Profile(context.Context) (*Profile, error)
	Daily(context.Context, time.Time, time.Time) ([]DailyMetrics, error)
	DataFreshness(context.Context) (*time.Time, error)
}

type Service struct {
	repo Repository
	now  func() time.Time
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo, now: func() time.Time { return time.Now().UTC() }}
}

func (s *Service) Profile(ctx context.Context) (*Profile, error) {
	return s.repo.Profile(ctx)
}

func (s *Service) Days(ctx context.Context, fromValue, toValue string) (DateRange, []DailyMetrics, error) {
	from, to, err := parseRange(fromValue, toValue, 14, 500, s.now)
	if err != nil {
		return DateRange{}, nil, err
	}
	days, err := s.repo.Daily(ctx, from, to)
	if err != nil {
		return DateRange{}, nil, err
	}
	normalizeDays(days)
	return formatRange(from, to), days, nil
}

func (s *Service) Latest(ctx context.Context) (*DailyMetrics, error) {
	from, to := recentRange(30, s.now)
	days, err := s.repo.Daily(ctx, from, to)
	if err != nil {
		return nil, err
	}
	normalizeDays(days)
	return firstNonEmpty(days), nil
}

func (s *Service) Context(ctx context.Context, days int) (ContextPack, error) {
	if days == 0 {
		days = 14
	}
	from, to := recentRange(days, s.now)
	recentDays, err := s.repo.Daily(ctx, from, to)
	if err != nil {
		return ContextPack{}, err
	}
	normalizeDays(recentDays)
	profile, err := s.repo.Profile(ctx)
	if err != nil {
		return ContextPack{}, err
	}
	freshness, err := s.repo.DataFreshness(ctx)
	if err != nil {
		return ContextPack{}, err
	}
	return ContextPack{
		GeneratedAt:         s.now(),
		DataFreshness:       freshness,
		Range:               formatRange(from, to),
		Profile:             profile,
		Latest:              firstNonEmpty(recentDays),
		Trends:              summarize(recentDays),
		RecentDays:          recentDays,
		InterpretationHints: interpretationHints(),
	}, nil
}

func parseRange(fromValue, toValue string, defaultDays, maxDays int, now func() time.Time) (time.Time, time.Time, error) {
	if fromValue == "" || toValue == "" {
		from, to := recentRange(defaultDays, now)
		return from, to, nil
	}
	from, err := time.Parse(time.DateOnly, fromValue)
	if err != nil {
		return time.Time{}, time.Time{}, err
	}
	to, err := time.Parse(time.DateOnly, toValue)
	if err != nil {
		return time.Time{}, time.Time{}, err
	}
	if to.Before(from) {
		return time.Time{}, time.Time{}, fmt.Errorf("to must be on or after from")
	}
	if int(to.Sub(from).Hours()/24)+1 > maxDays {
		from = to.AddDate(0, 0, -maxDays+1)
	}
	return from, to, nil
}

func recentRange(days int, now func() time.Time) (time.Time, time.Time) {
	if days < 1 {
		days = 1
	}
	to := now().In(time.Local)
	from := to.AddDate(0, 0, -days+1)
	return dateOnly(from), dateOnly(to)
}

func dateOnly(value time.Time) time.Time {
	parsed, _ := time.Parse(time.DateOnly, value.Format(time.DateOnly))
	return parsed
}

func formatRange(from, to time.Time) DateRange {
	return DateRange{From: from.Format(time.DateOnly), To: to.Format(time.DateOnly)}
}

func firstNonEmpty(days []DailyMetrics) *DailyMetrics {
	for i := range days {
		if days[i].Strain != nil || days[i].RecoveryScore != nil || days[i].SleepPerformancePercentage != nil || days[i].WorkoutCount > 0 {
			return &days[i]
		}
	}
	return nil
}

func summarize(days []DailyMetrics) TrendSummary {
	summary := TrendSummary{Days: len(days)}
	if latest := firstNonEmpty(days); latest != nil {
		summary.LatestRecoveryScore = latest.RecoveryScore
		summary.LatestStrain = latest.Strain
		summary.LatestSleepPerformancePercentage = latest.SleepPerformancePercentage
	}
	var recoveryValues, strainValues, sleepPerfValues, sleepHoursValues []float64
	for _, day := range days {
		if day.RecoveryScore != nil {
			recoveryValues = append(recoveryValues, float64(*day.RecoveryScore))
		}
		if day.Strain != nil {
			strainValues = append(strainValues, *day.Strain)
		}
		if day.SleepPerformancePercentage != nil {
			sleepPerfValues = append(sleepPerfValues, *day.SleepPerformancePercentage)
		}
		if day.SleepDurationHours != nil {
			sleepHoursValues = append(sleepHoursValues, *day.SleepDurationHours)
		}
		if day.WorkoutCount > 0 {
			summary.WorkoutDays++
			summary.WorkoutCount += day.WorkoutCount
		}
	}
	summary.RecoveryAverage = average(recoveryValues)
	summary.StrainAverage = average(strainValues)
	summary.SleepPerformanceAverage = average(sleepPerfValues)
	summary.SleepDurationAverageHours = average(sleepHoursValues)
	roundFloatPtr(summary.RecoveryAverage, 2)
	roundFloatPtr(summary.StrainAverage, 2)
	roundFloatPtr(summary.SleepPerformanceAverage, 2)
	roundFloatPtr(summary.SleepDurationAverageHours, 2)
	return summary
}

func average(values []float64) *float64 {
	if len(values) == 0 {
		return nil
	}
	var total float64
	for _, value := range values {
		total += value
	}
	result := total / float64(len(values))
	return &result
}

func normalizeDays(days []DailyMetrics) {
	for i := range days {
		roundFloatPtr(days[i].Strain, 2)
		roundFloatPtr(days[i].HRVRMSSDMilli, 2)
		roundFloatPtr(days[i].SleepPerformancePercentage, 2)
		roundFloatPtr(days[i].SleepEfficiencyPercentage, 2)
		roundFloatPtr(days[i].SleepDurationHours, 2)
		roundFloatPtr(days[i].WorkoutStrain, 2)
	}
}

func roundFloatPtr(value *float64, places int) {
	if value == nil {
		return
	}
	factor := math.Pow10(places)
	*value = math.Round(*value*factor) / factor
}

func interpretationHints() map[string]string {
	return map[string]string{
		"recovery_score":               "0-100, higher usually means more readiness for strain.",
		"strain":                       "0-21 WHOOP cardiovascular load scale; higher means more strain.",
		"workout_strain":               "Sum of workout strain values for the day. Because WHOOP strain is non-linear, treat it as a workout-load hint, not a value directly comparable to day strain.",
		"sleep_performance_percentage": "0-100, actual sleep compared with WHOOP sleep need.",
		"hrv_rmssd_milli":              "Heart-rate variability in milliseconds, best interpreted as personal trend.",
		"data_freshness":               "Most recent updated_at timestamp across synced WHOOP tables.",
	}
}
