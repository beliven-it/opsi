package gitlab

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestTokenExpiration(t *testing.T) {
	now := time.Date(2026, 10, 6, 15, 0, 0, 0, time.UTC)

	cases := []struct {
		name     string
		expires  string
		expiring int
		want     string
		wantShow bool
	}{
		{"no expiration, no filter", "", -1, "never expires", true},
		{"no expiration, filtered", "", 30, "never expires", false},
		{"expired is always shown", "2026-10-01", 30, "EXPIRED on 2026-10-01 (5 days ago)", true},
		{"expires today", "2026-10-06", 0, "expires today (2026-10-06)", true},
		{"within the window", "2026-10-26", 30, "expires 2026-10-26 (in 20 days)", true},
		{"outside the window", "2027-01-01", 30, "expires 2027-01-01 (in 87 days)", false},
		{"outside the window but no filter", "2027-01-01", -1, "expires 2027-01-01 (in 87 days)", true},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got, show := tokenExpiration(gitlabAccessToken{ExpiresAt: c.expires}, now, c.expiring)
			if got != c.want || show != c.wantShow {
				t.Errorf("got (%q, %v), want (%q, %v)", got, show, c.want, c.wantShow)
			}
		})
	}
}

func TestListGroupAccessTokens(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/groups":
			fmt.Fprint(w, `[
			  {"id":1,"full_path":"acme","web_url":"https://gl/groups/acme"},
			  {"id":2,"full_path":"acme/empty","web_url":"https://gl/groups/acme/empty"},
			  {"id":3,"full_path":"acme/locked","web_url":"https://gl/groups/acme/locked"}]`)
		case "/groups/1/access_tokens":
			fmt.Fprint(w, `[
			  {"name":"ci","scopes":["api"],"expires_at":"2020-01-01","revoked":false},
			  {"name":"old","scopes":["api"],"expires_at":"2020-01-01","revoked":true}]`)
		case "/groups/2/access_tokens":
			fmt.Fprint(w, `[]`)
		case "/groups/3/access_tokens":
			w.WriteHeader(http.StatusForbidden)
			fmt.Fprint(w, `{"message":"403 Forbidden"}`)
		default:
			t.Errorf("unexpected request %s", r.URL.Path)
		}
	}))
	defer server.Close()

	g := &gitlab{apiURL: server.URL, token: "t"}

	tokens, err := g.listGroupAccessTokens(1)
	if err != nil || len(tokens) != 1 || tokens[0].Name != "ci" {
		t.Errorf("revoked tokens must be excluded, got %v, %v", tokens, err)
	}

	if _, err := g.listGroupAccessTokens(3); err == nil || !strings.Contains(err.Error(), "403") {
		t.Errorf("expected the 403 error, got %v", err)
	}

	// A group without permission must not stop the whole listing
	if err := g.ListGroupAccessTokens(-1); err != nil {
		t.Errorf("listing failed: %v", err)
	}
}
