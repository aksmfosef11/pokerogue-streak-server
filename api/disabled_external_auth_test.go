package api

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestExternalAccountRoutesAreUnavailable(t *testing.T) {
	mux := http.NewServeMux()
	registerRoutes(mux)
	for _, route := range []struct{ method, path string }{
		{"GET", "/auth/google/callback"},
		{"POST", "/auth/google/logout"},
		{"GET", "/auth/discord/callback"},
		{"POST", "/auth/discord/logout"},
		{"POST", "/admin/account/googleLink"},
		{"POST", "/admin/account/googleUnlink"},
		{"POST", "/admin/account/discordLink"},
		{"POST", "/admin/account/discordUnlink"},
	} {
		response := httptest.NewRecorder()
		mux.ServeHTTP(response, httptest.NewRequest(route.method, route.path, nil))
		if response.Code != http.StatusNotFound {
			t.Errorf("%s %s: got %d, want 404", route.method, route.path, response.Code)
		}
	}
	for _, route := range []struct{ method, path string }{
		{"POST", "/account/login"},
		{"POST", "/account/register"},
		{"GET", "/admin/account/adminSearch"},
		{"POST", "/savedata/updateall"},
	} {
		_, pattern := mux.Handler(httptest.NewRequest(route.method, route.path, nil))
		if pattern == "" {
			t.Errorf("%s %s: expected route to remain available", route.method, route.path)
		}
	}
}
