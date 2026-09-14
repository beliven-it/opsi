package helpers

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRequestSendsMethodBodyQueryAndHeaders(t *testing.T) {
	var gotMethod, gotQuery, gotHeader, gotContentType string
	var gotBody []byte

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotQuery = r.URL.Query().Get("page")
		gotHeader = r.Header.Get("PRIVATE-TOKEN")
		gotContentType = r.Header.Get("Content-Type")
		gotBody, _ = io.ReadAll(r.Body)

		w.Write([]byte(`{"id":1}`))
	}))
	defer server.Close()

	response, err := Request(
		"PUT",
		server.URL,
		map[string]any{"mentions_disabled": true},
		map[string]string{"page": "2"},
		map[string]string{"Content-Type": "application/json", "PRIVATE-TOKEN": "secret"},
	)
	if err != nil {
		t.Fatal(err)
	}

	if gotMethod != "PUT" {
		t.Errorf("method = %q", gotMethod)
	}

	if gotQuery != "2" {
		t.Errorf("page query param = %q", gotQuery)
	}

	if gotHeader != "secret" || gotContentType != "application/json" {
		t.Errorf("headers = %q, %q", gotHeader, gotContentType)
	}

	var sent map[string]any
	if err := json.Unmarshal(gotBody, &sent); err != nil {
		t.Fatalf("body is not json: %q", gotBody)
	}

	if sent["mentions_disabled"] != true {
		t.Errorf("body = %v", sent)
	}

	if string(response) != `{"id":1}` {
		t.Errorf("response = %q", response)
	}
}

func TestRequestSendsNoBodyWhenNil(t *testing.T) {
	var gotBody []byte

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotBody, _ = io.ReadAll(r.Body)
	}))
	defer server.Close()

	if _, err := Request("GET", server.URL, nil, nil, nil); err != nil {
		t.Fatal(err)
	}

	if len(gotBody) != 0 {
		t.Errorf("a nil body should send nothing, got %q", gotBody)
	}
}

func TestRequestStatusHandling(t *testing.T) {
	cases := []struct {
		name      string
		status    int
		wantError bool
	}{
		{"ok", http.StatusOK, false},
		{"created", http.StatusCreated, false},
		// The accepted range ends at 226, so redirects and up are errors
		{"im used", http.StatusIMUsed, false},
		{"moved", http.StatusMovedPermanently, true},
		{"bad request", http.StatusBadRequest, true},
		{"server error", http.StatusInternalServerError, true},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(c.status)
				w.Write([]byte("the reason"))
			}))
			defer server.Close()

			response, err := Request("GET", server.URL, nil, nil, nil)

			if c.wantError {
				if err == nil {
					t.Fatalf("status %d should be an error", c.status)
				}

				// The caller needs the API message, not just the status
				if err.Error() != "the reason" {
					t.Errorf("error should carry the response body, got %q", err.Error())
				}

				if response != nil {
					t.Errorf("a failed request should return no payload, got %q", response)
				}

				return
			}

			if err != nil {
				t.Fatalf("status %d should not be an error: %v", c.status, err)
			}
		})
	}
}

func TestRequestFailsOnUnreachableHost(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	url := server.URL
	server.Close()

	if _, err := Request("GET", url, nil, nil, nil); err == nil {
		t.Error("a closed server should return an error")
	}
}
