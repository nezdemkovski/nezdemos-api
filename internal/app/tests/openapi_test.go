package app_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"nezdemos-api/internal/app"
	"nezdemos-api/internal/config"
)

func TestOpenAPISpecIsAgentFriendly(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/openapi.json", nil)
	rec := httptest.NewRecorder()
	app.NewHandler(config.Settings{APIKey: "secret"}, nil).ServeHTTP(rec, req)
	resp := rec.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	spec := string(body)
	for _, expected := range []string{
		`"operationId":"getWhoopAgentContext"`,
		`"operationId":"listWhoopDays"`,
		`"operationId":"getLatestWhoopDay"`,
		`"operationId":"getWhoopProfile"`,
		`"operationId":"getHealth"`,
		`"tags":["WHOOP","Agents"]`,
		`"security":[]`,
		`"X-API-Key"`,
		`"WHOOP recovery score, 0-100 scale."`,
		`"Application version embedded at build time."`,
	} {
		if !strings.Contains(spec, expected) {
			t.Fatalf("expected OpenAPI spec to contain %s", expected)
		}
	}
}

func TestHealthIncludesVersion(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	rec := httptest.NewRecorder()
	app.NewHandler(config.Settings{APIKey: "secret"}, nil).ServeHTTP(rec, req)

	resp := rec.Result()
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	for _, expected := range []string{`"status":"ok"`, `"version":"dev"`} {
		if !strings.Contains(string(body), expected) {
			t.Fatalf("expected health response to contain %s, got %s", expected, string(body))
		}
	}
}
