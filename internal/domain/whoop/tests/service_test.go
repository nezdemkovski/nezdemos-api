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
