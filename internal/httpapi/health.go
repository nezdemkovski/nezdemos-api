package httpapi

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

type HealthOutput struct {
	Body struct {
		Status  string `json:"status" example:"ok" doc:"Health status."`
		Version string `json:"version" example:"0.1.0" doc:"Application version embedded at build time."`
	}
}

func RegisterHealth(api huma.API, version string) {
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
		out.Body.Version = version
		return out, nil
	})
}
