package whoop

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	domain "nezdemos-api/internal/domain/whoop"
)

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Profile(ctx context.Context) (*domain.Profile, error) {
	row := r.db.QueryRow(ctx, profileSQL)
	var profile domain.Profile
	if err := row.Scan(&profile.UserID, &profile.Email, &profile.FirstName, &profile.LastName); err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &profile, nil
}

func (r *Repository) Daily(ctx context.Context, from, to time.Time) ([]domain.DailyMetrics, error) {
	rows, err := r.db.Query(ctx, dailySQL, from, to)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var days []domain.DailyMetrics
	for rows.Next() {
		var day domain.DailyMetrics
		if err := rows.Scan(
			&day.Date,
			&day.Strain,
			&day.AverageHeartRate,
			&day.MaxHeartRate,
			&day.RecoveryScore,
			&day.RestingHeartRate,
			&day.HRVRMSSDMilli,
			&day.SleepPerformancePercentage,
			&day.SleepEfficiencyPercentage,
			&day.SleepDurationHours,
			&day.WorkoutCount,
			&day.WorkoutStrain,
		); err != nil {
			return nil, err
		}
		days = append(days, day)
	}
	return days, rows.Err()
}

func (r *Repository) DataFreshness(ctx context.Context) (*time.Time, error) {
	row := r.db.QueryRow(ctx, dataFreshnessSQL)
	var freshness *time.Time
	if err := row.Scan(&freshness); err != nil {
		return nil, err
	}
	return freshness, nil
}
