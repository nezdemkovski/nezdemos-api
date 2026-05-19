package whoop_test

import (
	"context"
	"testing"
	"time"

	"nezdemos-api/internal/domain/whoop"
)

func TestLatestSkipsEmptyCalendarDays(t *testing.T) {
	strain := 11.2
	recovery := 73
	service := whoop.NewService(fakeRepository{days: []whoop.DailyMetrics{
		{Date: "2026-05-19"},
		{Date: "2026-05-18", Strain: &strain, RecoveryScore: &recovery},
	}})

	latest, err := service.Latest(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if latest == nil {
		t.Fatal("expected latest day")
	}
	if latest.Date != "2026-05-18" {
		t.Fatalf("expected 2026-05-18, got %s", latest.Date)
	}
}

func TestSummarizeUsesLatestNonEmptyDay(t *testing.T) {
	strain := 8.5
	recovery := 66
	sleep := 91.0
	service := whoop.NewService(fakeRepository{days: []whoop.DailyMetrics{
		{Date: "2026-05-19"},
		{Date: "2026-05-18", Strain: &strain, RecoveryScore: &recovery, SleepPerformancePercentage: &sleep, WorkoutCount: 1},
	}})

	contextPack, err := service.Context(context.Background(), 14)
	if err != nil {
		t.Fatal(err)
	}
	summary := contextPack.Trends
	if summary.LatestRecoveryScore == nil || *summary.LatestRecoveryScore != recovery {
		t.Fatalf("expected latest recovery %d, got %#v", recovery, summary.LatestRecoveryScore)
	}
	if summary.LatestStrain == nil || *summary.LatestStrain != strain {
		t.Fatalf("expected latest strain %.1f, got %#v", strain, summary.LatestStrain)
	}
	if summary.WorkoutDays != 1 || summary.WorkoutCount != 1 {
		t.Fatalf("expected one workout day and one workout, got %d/%d", summary.WorkoutDays, summary.WorkoutCount)
	}
}

func TestContextRoundsFloatMetrics(t *testing.T) {
	strain := 6.6559
	hrv := 38.944
	sleepPerformance := 82.857142857
	sleepEfficiency := 94.731
	sleepHours := 6.594658333333333
	workoutStrain := 10.5299
	service := whoop.NewService(fakeRepository{days: []whoop.DailyMetrics{
		{
			Date:                       "2026-05-19",
			Strain:                     &strain,
			HRVRMSSDMilli:              &hrv,
			SleepPerformancePercentage: &sleepPerformance,
			SleepEfficiencyPercentage:  &sleepEfficiency,
			SleepDurationHours:         &sleepHours,
			WorkoutCount:               2,
			WorkoutStrain:              &workoutStrain,
		},
	}})

	contextPack, err := service.Context(context.Background(), 14)
	if err != nil {
		t.Fatal(err)
	}
	day := contextPack.RecentDays[0]
	assertFloat(t, "strain", day.Strain, 6.66)
	assertFloat(t, "hrv", day.HRVRMSSDMilli, 38.94)
	assertFloat(t, "sleep performance", day.SleepPerformancePercentage, 82.86)
	assertFloat(t, "sleep efficiency", day.SleepEfficiencyPercentage, 94.73)
	assertFloat(t, "sleep hours", day.SleepDurationHours, 6.59)
	assertFloat(t, "workout strain", day.WorkoutStrain, 10.53)
	assertFloat(t, "strain average", contextPack.Trends.StrainAverage, 6.66)
	assertFloat(t, "sleep performance average", contextPack.Trends.SleepPerformanceAverage, 82.86)
	assertFloat(t, "sleep hours average", contextPack.Trends.SleepDurationAverageHours, 6.59)
}

func assertFloat(t *testing.T, name string, actual *float64, expected float64) {
	t.Helper()
	if actual == nil {
		t.Fatalf("expected %s %.2f, got nil", name, expected)
	}
	if *actual != expected {
		t.Fatalf("expected %s %.2f, got %.12f", name, expected, *actual)
	}
}

type fakeRepository struct {
	days []whoop.DailyMetrics
}

func (r fakeRepository) Profile(context.Context) (*whoop.Profile, error) {
	return &whoop.Profile{UserID: "test-user"}, nil
}

func (r fakeRepository) Daily(context.Context, time.Time, time.Time) ([]whoop.DailyMetrics, error) {
	return r.days, nil
}

func (r fakeRepository) DataFreshness(context.Context) (*time.Time, error) {
	now := time.Date(2026, 5, 19, 12, 0, 0, 0, time.UTC)
	return &now, nil
}
