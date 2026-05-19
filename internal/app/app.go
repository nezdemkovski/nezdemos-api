package app

import (
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	pgwhoop "nezdemos-api/internal/adapters/postgres/whoop"
	"nezdemos-api/internal/buildinfo"
	"nezdemos-api/internal/config"
	domainwhoop "nezdemos-api/internal/domain/whoop"
	"nezdemos-api/internal/httpapi"
	httpwhoop "nezdemos-api/internal/httpapi/whoop"
)

func NewHandler(cfg config.Settings, db *pgxpool.Pool) http.Handler {
	router := chi.NewRouter()
	router.Use(httpapi.APIKeyMiddleware(cfg.APIKey))

	api := humachi.New(router, huma.DefaultConfig("Nezdemos API", buildinfo.Version))
	configureOpenAPI(api)

	httpapi.RegisterHealth(api, buildinfo.Version)
	whoopRepo := pgwhoop.NewRepository(db)
	whoopService := domainwhoop.NewService(whoopRepo)
	httpwhoop.Register(api, whoopService)

	return router
}

func configureOpenAPI(api huma.API) {
	api.OpenAPI().Info.Description = "Read-only API over synced personal data, optimized for agent tool calls."
	api.OpenAPI().Servers = []*huma.Server{{URL: "http://localhost:8080"}}
	api.OpenAPI().Components.SecuritySchemes = map[string]*huma.SecurityScheme{
		"api_key": {Type: "apiKey", In: "header", Name: "X-API-Key"},
	}
	api.OpenAPI().Security = []map[string][]string{{"api_key": {}}}
}
