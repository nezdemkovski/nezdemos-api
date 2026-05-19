package whoop

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"

	domain "nezdemos-api/internal/domain/whoop"
)

type Service interface {
	Profile(context.Context) (*domain.Profile, error)
	Days(context.Context, string, string) (domain.DateRange, []domain.DailyMetrics, error)
	Latest(context.Context) (*domain.DailyMetrics, error)
	Context(context.Context, int) (domain.ContextPack, error)
}

type ProfileOutput struct {
	Body domain.Profile
}

type DaysInput struct {
	From string `query:"from" doc:"Start date in YYYY-MM-DD format. Defaults to 14 days ago when omitted." example:"2026-05-01"`
	To   string `query:"to" doc:"End date in YYYY-MM-DD format. Defaults to today when omitted." example:"2026-05-19"`
}

type DaysOutput struct {
	Body struct {
		Range domain.DateRange      `json:"range"`
		Days  []domain.DailyMetrics `json:"days"`
	}
}

type LatestOutput struct {
	Body struct {
		Latest *domain.DailyMetrics `json:"latest,omitempty"`
	}
}

type ContextInput struct {
	Days int `query:"days" default:"14" minimum:"1" maximum:"90" doc:"Number of recent days to include"`
}

type ContextOutput struct {
	Body domain.ContextPack
}

func Register(api huma.API, service Service) {
	huma.Register(api, huma.Operation{
		OperationID: "getWhoopProfile",
		Method:      http.MethodGet,
		Path:        "/whoop/profile",
		Summary:     "Get WHOOP profile",
		Description: "Return the latest synced WHOOP user profile. Use this to identify whose health data the other WHOOP endpoints describe.",
		Tags:        []string{"WHOOP"},
		Errors:      []int{http.StatusNotFound},
	}, func(ctx context.Context, input *struct{}) (*ProfileOutput, error) {
		profile, err := service.Profile(ctx)
		if err != nil {
			return nil, err
		}
		if profile == nil {
			return nil, huma.Error404NotFound("whoop profile has not been synced yet")
		}
		return &ProfileOutput{Body: *profile}, nil
	})

	huma.Register(api, huma.Operation{
		OperationID: "listWhoopDays",
		Method:      http.MethodGet,
		Path:        "/whoop/days",
		Summary:     "List WHOOP daily metrics",
		Description: "Return day-level WHOOP metrics by joining cycles, recovery, sleep, and workouts. This is the preferred endpoint when an agent needs a table-like range for analysis.",
		Tags:        []string{"WHOOP"},
		Errors:      []int{http.StatusBadRequest},
	}, func(ctx context.Context, input *DaysInput) (*DaysOutput, error) {
		dateRange, days, err := service.Days(ctx, input.From, input.To)
		if err != nil {
			return nil, huma.Error400BadRequest(err.Error())
		}
		out := &DaysOutput{}
		out.Body.Range = dateRange
		out.Body.Days = days
		return out, nil
	})

	huma.Register(api, huma.Operation{
		OperationID: "getLatestWhoopDay",
		Method:      http.MethodGet,
		Path:        "/whoop/latest",
		Summary:     "Get latest WHOOP day",
		Description: "Return the most recent non-empty daily WHOOP metrics from the last 30 days. Use this for quick current-state checks.",
		Tags:        []string{"WHOOP"},
	}, func(ctx context.Context, input *struct{}) (*LatestOutput, error) {
		latest, err := service.Latest(ctx)
		if err != nil {
			return nil, err
		}
		return &LatestOutput{Body: struct {
			Latest *domain.DailyMetrics `json:"latest,omitempty"`
		}{Latest: latest}}, nil
	})

	huma.Register(api, huma.Operation{
		OperationID: "getWhoopAgentContext",
		Method:      http.MethodGet,
		Path:        "/whoop/context",
		Summary:     "Get WHOOP agent context",
		Description: "Return a compact agent-ready context pack with profile, latest day, recent daily metrics, trend summary, data freshness, and interpretation hints. Prefer this endpoint for LLM tool calls.",
		Tags:        []string{"WHOOP", "Agents"},
	}, func(ctx context.Context, input *ContextInput) (*ContextOutput, error) {
		contextPack, err := service.Context(ctx, input.Days)
		if err != nil {
			return nil, err
		}
		return &ContextOutput{Body: contextPack}, nil
	})
}
