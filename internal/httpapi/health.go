package httpapi

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

type HealthOutput struct {
	Body struct {
		Status string `json:"status" example:"ok"`
	}
}

func RegisterHealth(api huma.API) {
	huma.Register(api, huma.Operation{
		OperationID: "getHealth",
		Method:      http.MethodGet,
		Path:        "/healthz",
		Summary:     "Check API health",
		Description: "Public liveness endpoint for load balancers and uptime checks.",
		Tags:        []string{"System"},
		Security:    []map[string][]string{},
	}, func(ctx context.Context, input *struct{}) (*HealthOutput, error) {
		out := &HealthOutput{}
		out.Body.Status = "ok"
		return out, nil
	})
}
